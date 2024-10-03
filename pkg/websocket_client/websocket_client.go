package websocketclient

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"studentgit.kata.academy/quant/torque/pkg/logster"

	"github.com/gorilla/websocket"
)

type ConnectionState int

const (
	Disconnected ConnectionState = iota + 1
	Connecting
	Connected
)

const (
	messageRetryInterval = time.Millisecond * 100
)

type Client struct {
	addr     string
	logger   logster.Logger
	messages chan []byte

	conn      *websocket.Conn
	ctx       context.Context
	connState int32
	errCh     chan error
	mu        sync.Mutex
	cancel    context.CancelFunc

	// Options
	OnError                    func(error)
	OnConnected                func()
	OnDisconnected             func()
	sendQueue                  chan []byte
	sendBufferSize             int
	receiveBufferSize          int
	retryTimes                 int
	handshakeTimeout           time.Duration
	reconnectInterval          time.Duration
	readDeadline               time.Duration
	sendMessagesAfterReconnect [][]byte
	reconnectHistory           []time.Time   // Stores the timestamps of reconnect attempts
	maxReconnects              int           // Maximum allowed reconnects within the specified time frame
	reconnectWindow            time.Duration // Time frame for checking reconnect attempts
}

//nolint:gomnd // we specify default values here, it's ok
func NewClient(addr string, logger logster.Logger, options ...Option) *Client {
	client := &Client{

		addr:     addr,
		logger:   logger,
		messages: make(chan []byte),
		errCh:    make(chan error),

		ctx:    nil,
		cancel: nil,

		sendBufferSize:    1024,                 // default
		receiveBufferSize: 1024,                 // default
		sendQueue:         make(chan []byte, 5), // default
		retryTimes:        2,                    // default
		handshakeTimeout:  5 * time.Second,      // default
		reconnectInterval: 2 * time.Second,      // default
		readDeadline:      2 * time.Minute,
		maxReconnects:     5,                // default value
		reconnectWindow:   30 * time.Second, // default value

		OnError:                    nil,
		OnConnected:                nil,
		OnDisconnected:             nil,
		sendMessagesAfterReconnect: nil,

		conn:             nil,
		connState:        int32(Disconnected),
		reconnectHistory: []time.Time{},
		mu:               sync.Mutex{},
	}
	for _, option := range options {
		option(client)
	}
	return client
}

func (c *Client) Send(message []byte) error {
	// don't allow to send messages to closed channel
	select {
	case c.sendQueue <- message:
		return nil
	default:
		return fmt.Errorf("send queue is full")
	}
}

func (c *Client) GetMessages() <-chan []byte {
	return c.messages
}

func (c *Client) Connect(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.ctx = ctx
	c.cancel = cancel
	defer func() {
		go c.manageConnection()
		go c.readMessages()
		go c.sendMessages()
	}()
	return c.connect()
}

func (c *Client) Shutdown() {
	c.cancel()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.WithError(err).Errorf("Error closing connection")
		}
	}
	close(c.sendQueue)
	close(c.errCh)
	close(c.messages)
	c.logger.Infof("Client shutdown")
}

func (c *Client) UpdateHost(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.addr = host
}

func (c *Client) GetHost() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.addr
}

func (c *Client) SendMessageAfterReconnect(msgs [][]byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sendMessagesAfterReconnect = msgs
}

func (c *Client) setConnectionState(state ConnectionState) {
	atomic.StoreInt32(&c.connState, int32(state))
}

func (c *Client) IsConnected() bool {
	return c.getConnectionState() == Connected
}

func (c *Client) getConnectionState() ConnectionState {
	return ConnectionState(atomic.LoadInt32(&c.connState))
}

func (c *Client) connect() error {
	c.setConnectionState(Connecting)
	//nolint:exhaustruct // OK
	dialer := websocket.Dialer{
		HandshakeTimeout: c.handshakeTimeout,
		ReadBufferSize:   c.receiveBufferSize,
		WriteBufferSize:  c.sendBufferSize,
	}
	now := time.Now()
	c.logger.WithField("address", c.addr).Infof("Sending connect")
	conn, resp, err := dialer.Dial(c.addr, nil)
	//nolint:nestif // Need for logging if some unusual event happens
	if err != nil {
		duration := time.Since(now)
		if resp != nil {
			bodyBytes, errRead := io.ReadAll(resp.Body)
			if errRead != nil {
				c.logger.WithError(errRead).Errorf("Error reading response body")
			}
			errX := resp.Body.Close()
			if errX != nil {
				c.logger.Infof("failed to close body")
			}
			c.logger.WithFields(logster.Fields{
				"status":  resp.Status,
				"headers": resp.Header,
				"body":    string(bodyBytes),
			}).WithError(err).WithField("duration", duration).Errorf("Connection failed with response")
		} else {
			c.logger.WithError(err).WithField("duration", duration).Errorf("Connection failed without response")
		}

		c.setConnectionState(Disconnected)
		if c.OnError != nil {
			c.OnError(err)
		}
		return err
	}
	duration := time.Since(now)
	c.logger.WithField("address", c.addr).
		WithField("connectionDuration", duration).Infof("Successfully connected")
	c.conn = conn
	c.setConnectionState(Connected)
	if c.OnConnected != nil {
		c.OnConnected()
	}
	return nil
}

func (c *Client) reconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.connect()
	if err == nil {
		go c.readMessages()
		if c.sendMessagesAfterReconnect != nil {
			for _, msg := range c.sendMessagesAfterReconnect {
				c.sendQueue <- msg
				time.Sleep(messageRetryInterval) // If we need to send multiple messages,
				// the first message can be login, so we have to wait for some time before sending the next message
			}
		}
	}
	return err
}

func (c *Client) readMessages() {
	type ResponseResult struct {
		data []byte
		err  error
	}
	resultChan := make(chan ResponseResult)
	// We need this goroutine to be able to exit when the context is done because ReadMessage is blocking.
	go func() {
		for {
			_ = c.conn.SetReadDeadline(time.Now().Add(c.readDeadline))
			_, data, err := c.conn.ReadMessage()
			if err != nil {
				c.logger.WithError(err).Debugf("error occurred on reading from channel")
				resultChan <- ResponseResult{
					data: nil,
					err:  err,
				}
				close(resultChan)
				return
			}
			resultChan <- ResponseResult{
				data: data,
				err:  nil,
			}
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Debugf("context is done, exiting readMessages")
			return
		case result, ok := <-resultChan:
			if !ok {
				c.logger.Debugf("result channel is closed, exiting readMessages")
				return
			}
			if result.err != nil {
				c.logger.WithError(result.err).Debugf("received error reading message")
				c.errCh <- result.err
				return
			}
			c.messages <- result.data
		}
	}
}

func (c *Client) sendMessages() {
	retryMap := make(map[string]int)
	for {
		select {
		case <-c.ctx.Done():
			c.logger.Debugf("context is done, exiting sendMessages")
			return
		case message, ok := <-c.sendQueue:
			if !ok {
				c.logger.Debugf("send channel is closed")
				return
			}
			c.logger.WithField("message", string(message)).Debugf("getting ready to send message")
			if c.getConnectionState() != Connected {
				c.logger.Debugf("connection is not ready, sending message back to the queue")
				c.sendQueue <- message
				time.Sleep(messageRetryInterval) // Wait if the connection is not ready.
				continue
			}
			messageID := string(message)
			if retryMap[messageID] >= c.retryTimes {
				c.logger.WithField("message", messageID).WithField("retryTimes", c.retryTimes).Errorf(
					"failed to send message")
				delete(retryMap, messageID)
				continue
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				c.logger.WithError(err).Errorf("failed to send message, sending message back to the queue")
				c.sendQueue <- message
				retryMap[messageID]++
				if c.OnError != nil {
					c.OnError(err)
				}
				// Send error to connection manager.
				// The manageConnection method will handle reconnection.
				c.errCh <- err
				continue
			}
			c.logger.WithField("message", messageID).Debugf("message sent successfully")
			delete(retryMap, messageID)
		}
	}
}

func (c *Client) manageConnection() {
	for {
		select {
		case <-c.ctx.Done():
			c.logger.Debugf("context is done, exiting manageConnection")
			return
		case <-c.errCh:
			now := time.Now()
			c.reconnectHistory = append(c.reconnectHistory, now)

			// Remove timestamps outside the time window
			cutoff := now.Add(-c.reconnectWindow)
			for len(c.reconnectHistory) > 0 && c.reconnectHistory[0].Before(cutoff) {
				c.reconnectHistory = c.reconnectHistory[1:]
			}

			if len(c.reconnectHistory) > c.maxReconnects {
				panic(fmt.Errorf("exceeded maximum reconnect attempts (%d) within %v", c.maxReconnects, c.reconnectWindow))
			}

			c.setConnectionState(Disconnected)
			c.logger.Warnf("Got disconnected")
			err := c.reconnect()
			if err != nil {
				c.logger.WithError(err).WithField("reconnect_interval", c.reconnectInterval).Errorf("failed to connect with error")
				time.Sleep(c.reconnectInterval)
				go func(err error) {
					c.errCh <- err // because the channel is unbuffered we need to send message in goroutine not to block
				}(err)
			}
		}
	}
}
