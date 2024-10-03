package streamer

import "errors"

var (
	ErrChannelClosed      = errors.New("channel closed")
	ErrMaxAttempts        = errors.New("max reconnect attempts reached")
	ErrConnectionToServer = errors.New("failed to connect to server")
	ErrContextTimout      = errors.New("context timeout")
	ErrStreamIsClosed     = errors.New("stream is closed")
)
