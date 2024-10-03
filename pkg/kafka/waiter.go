package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"studentgit.kata.academy/quant/torque/pkg/changroup"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type WaitingHandler struct {
	logger        logster.Logger
	next          sarama.ConsumerGroupHandler
	client        sarama.Client
	groupID       string
	invalid       *changroup.Group[struct{}]
	mu            sync.Mutex
	session       *WaitingSession
	claims        []sarama.ConsumerGroupClaim
	ready         chan struct{} // ready is closed if all claims are received
	shouldConsume map[string]map[int32]struct{}
}

func NewWaitingHandler(logger logster.Logger, next sarama.ConsumerGroupHandler, groupID string) *WaitingHandler {
	return &WaitingHandler{
		logger:        logger,
		next:          next,
		client:        nil, // will be filled in SetClient
		groupID:       groupID,
		invalid:       changroup.NewGroup[struct{}](),
		mu:            sync.Mutex{},
		session:       nil,
		claims:        nil,
		ready:         make(chan struct{}),
		shouldConsume: nil,
	}
}

func (w *WaitingHandler) SetClient(client sarama.Client) {
	if w.client != nil {
		panic("client already set")
	}
	w.client = client
}

// WaitNewestMarked gets current high water mark offsets for all consumed topics/partitions
// and blocks until those offsets (or greater) are marked.
// If consumer hasn't started yet (or rebalance happening), it waits as well.
// If context is done, it returns ctx.Err().
func (w *WaitingHandler) WaitNewestMarked(ctx context.Context) error {
	for {
		retry, err := w.waitNewestMarked(ctx)
		if err != nil {
			return err
		}
		if !retry {
			return nil
		}
	}
}

func (w *WaitingHandler) waitNewestMarked(ctx context.Context) (retry bool, _ error) { //nolint:nonamedreturns // docs
	invalid, release := w.invalid.Acquire()
	defer release()

	w.mu.Lock()
	ready := w.ready
	w.mu.Unlock()

	select {
	case <-ready:
	case <-invalid: // edge case: consumption stopped before all claims are registered
		return true, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}

	w.mu.Lock()
	session := w.session
	claims := w.claims
	w.mu.Unlock()

	// this may happen because we wait for <-ready without lock
	if session == nil || len(claims) == 0 {
		return true, nil
	}

	waitCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, len(claims)+1)

	for _, claim := range claims {
		topic := claim.Topic()
		partition := claim.Partition()
		highWaterMarkOffset := claim.HighWaterMarkOffset()
		go func() {
			errs <- session.WaitMarkOffset(waitCtx, topic, partition, highWaterMarkOffset)
		}()
	}

	var errInvalid = errors.New("invalid: rebalance or consumption stopped")
	go func() {
		<-invalid
		errs <- errInvalid
	}()

	for range claims {
		err := <-errs
		if err != nil && errors.Is(err, errInvalid) {
			return true, nil
		}
		if err != nil {
			return false, err
		}
	}

	return false, nil
}

func (w *WaitingHandler) Setup(session sarama.ConsumerGroupSession) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.session != nil {
		return fmt.Errorf("something went wrong: setup is called on ready wrapper")
	}

	count := 0
	claims := session.Claims()
	w.shouldConsume = make(map[string]map[int32]struct{}, len(claims))
	for topic, partitions := range claims {
		count += len(partitions)
		m := make(map[int32]struct{}, len(partitions))
		for _, partition := range partitions {
			m[partition] = struct{}{}
		}
		w.shouldConsume[topic] = m
	}
	w.claims = make([]sarama.ConsumerGroupClaim, 0, count)

	coordinator, err := w.client.Coordinator(w.groupID)
	if err != nil {
		return fmt.Errorf("client.Coordinator: %w", err)
	}

	offsets, err := getCommittedOffsets(coordinator, w.groupID, w.shouldConsume)
	if err != nil {
		return fmt.Errorf("getCommittedOffsets: %w", err)
	}

	w.session = NewWaitingSession(w.logger, session, offsets)

	return w.next.Setup(w.session)
}

func (w *WaitingHandler) Cleanup(sarama.ConsumerGroupSession) error {
	defer func() {
		w.mu.Lock()
		defer w.mu.Unlock()
		w.session = nil
		w.claims = nil
		w.ready = make(chan struct{})
		w.shouldConsume = nil
		w.invalid.SendAsync(struct{}{}) // SendAsync is used to avoid deadlock
	}()
	return w.next.Cleanup(w.session)
}

func (w *WaitingHandler) ConsumeClaim(_ sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	err := w.registerClaim(claim)
	if err != nil {
		return err
	}

	return w.next.ConsumeClaim(w.session, claim)
}

func (w *WaitingHandler) registerClaim(claim sarama.ConsumerGroupClaim) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	topic := claim.Topic()

	partitions, ok := w.shouldConsume[topic]
	if !ok {
		return fmt.Errorf(
			"something went wrong: unexpected topic received %q, should consume: %v",
			topic,
			w.shouldConsume,
		)
	}

	partition := claim.Partition()

	_, ok = partitions[partition]
	if !ok {
		return fmt.Errorf(
			"something went wrong: unexpected partition received %q/%d, should consume: %v",
			topic,
			partition,
			w.shouldConsume,
		)
	}

	w.claims = append(w.claims, claim)

	delete(partitions, partition)

	if len(partitions) == 0 {
		delete(w.shouldConsume, topic)
	}

	if len(w.shouldConsume) == 0 {
		close(w.ready)
	}

	return nil
}

type WaitingSession struct {
	logger  logster.Logger
	next    sarama.ConsumerGroupSession
	notify  *changroup.Group[struct{}]
	mu      sync.RWMutex
	offsets map[string]map[int32]int64
}

func NewWaitingSession(
	logger logster.Logger,
	next sarama.ConsumerGroupSession,
	offsets map[string]map[int32]int64,
) *WaitingSession {
	return &WaitingSession{
		logger:  logger,
		next:    next,
		notify:  changroup.NewGroup[struct{}](),
		mu:      sync.RWMutex{},
		offsets: offsets,
	}
}

func (w *WaitingSession) WaitMarkOffset(ctx context.Context, topic string, partition int32, offset int64) error {
	w.mu.RLock()
	unlocked := false
	defer func() {
		if !unlocked {
			w.mu.RUnlock()
		}
	}()

	partitions, ok := w.offsets[topic]
	if !ok {
		return fmt.Errorf("unknown topic %q", topic)
	}

	current, ok := partitions[partition]
	if !ok {
		return fmt.Errorf("unknown partition %q/%d", topic, partition)
	}

	fields := logster.Fields{
		"topic":           topic,
		"partition":       partition,
		"expected_offset": offset,
		"current_offset":  current,
	}

	if current >= offset {
		w.logger.WithFields(fields).Debugf("WaitingSession: offset already marked")
		return nil
	}

	notify, release := w.notify.Acquire()
	defer release()

	w.mu.RUnlock()
	unlocked = true

	start := time.Now()
	w.logger.WithFields(fields).Debugf("WaitingSession: start waiting offset")

	for {
		select {
		case <-ctx.Done():
			delete(fields, "current_offset")
			fields["duration"] = time.Since(start)
			w.logger.WithFields(fields).Debugf("WaitingSession: context done")
			return ctx.Err()
		case <-notify:
			w.mu.RLock()
			current = w.offsets[topic][partition]
			w.mu.RUnlock()

			if current >= offset {
				fields["current_offset"] = current
				fields["duration"] = time.Since(start)
				w.logger.WithFields(fields).Debugf("WaitingSession: offset marked")
				return nil
			}
		}
	}
}

func (w *WaitingSession) mark(topic string, partition int32, offset int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	partitions, ok := w.offsets[topic]
	if !ok {
		return
	}

	_, ok = partitions[partition]
	if !ok {
		return
	}

	partitions[partition] = offset

	w.notify.SendAsync(struct{}{}) // SendAsync is used to not block without need
}

func (w *WaitingSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	w.next.MarkOffset(topic, partition, offset, metadata)
	w.mark(topic, partition, offset)
}

func (w *WaitingSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	w.next.MarkMessage(msg, metadata)
	// +1 is taken from [sarama.consumerGroupSession.MarkMessage]
	w.mark(msg.Topic, msg.Partition, msg.Offset+1)
}

func (w *WaitingSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
	w.next.ResetOffset(topic, partition, offset, metadata)
	w.mark(topic, partition, offset)
}

func (w *WaitingSession) Claims() map[string][]int32 {
	return w.next.Claims()
}

func (w *WaitingSession) MemberID() string {
	return w.next.MemberID()
}

func (w *WaitingSession) GenerationID() int32 {
	return w.next.GenerationID()
}

func (w *WaitingSession) Commit() {
	w.next.Commit()
}

func (w *WaitingSession) Context() context.Context {
	return w.next.Context()
}
