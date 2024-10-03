package kafka_test

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/kernle32dll/testcontainers-go-canned/kafkatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"studentgit.kata.academy/quant/torque/pkg/kafka"
	"studentgit.kata.academy/quant/torque/pkg/kafka/mocks"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

func TestWaitingHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("test starts kafka in container")
	}
	if os.Getenv("CI") != "" {
		t.Skip("testcontainers can be used in gitlab-ci (for now?)")
	}
	t.Parallel()
	const messagesBeforeStart = 10
	ctx := context.Background()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var groupID = "group" + strconv.Itoa(r.Int())
	var topic = "topic" + strconv.Itoa(r.Int())

	// PREPARATION

	cluster, err := kafkatest.NewKafkaCluster(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { cluster.StopCluster(ctx) })

	addr, err := cluster.GetKafkaHostAndPort(ctx)
	require.NoError(t, err)

	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true // for SyncProducer

	client, err := sarama.NewClient([]string{addr}, conf)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, client.Close()) })

	producer, err := sarama.NewSyncProducerFromClient(client)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, producer.Close()) })

	for i := 0; i < messagesBeforeStart; i++ {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder("before start " + strconv.Itoa(i)),
		}
		_, _, err = producer.SendMessage(msg)
		require.NoError(t, err)
	}

	group, err := sarama.NewConsumerGroupFromClient(groupID, client)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, group.Close()) })

	handler := mocks.NewConsumerGroupHandler(t)
	waiter := kafka.NewWaitingHandler(logster.Discard, handler, groupID)
	waiter.SetClient(client)

	ready := make(chan struct{})
	offsets := make(chan int64)
	stop := make(chan struct{})

	mocks.InOrder(
		handler.EXPECT().Setup(mock.Anything).Return(nil).Once(),
		handler.EXPECT().
			ConsumeClaim(mock.Anything, mock.Anything).
			Run(func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) {
				close(ready)
				for {
					select {
					case <-stop:
						return
					case msg, ok := <-claim.Messages():
						if !ok {
							require.Fail(t, "unexpected stop")
						}
						offsets <- msg.Offset
						session.MarkMessage(msg, "")
					}
				}
			}).
			Return(nil).
			Once(),
		handler.EXPECT().Cleanup(mock.Anything).Return(nil).Once(),
	)

	consumeErr := make(chan error)
	go func() {
		consumeErr <- group.Consume(ctx, []string{topic}, waiter)
	}()

	waitChan(t, ready)

	t.Run("consuming new messages", func(t *testing.T) {
		// must be non-parallel

		for i := 0; i < 5; i++ {
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder("after start " + strconv.Itoa(i)),
			}
			_, _, err = producer.SendMessage(msg)
			require.NoError(t, err)
		}

		require.EqualValues(t, messagesBeforeStart+0, waitChan(t, offsets))
		require.EqualValues(t, messagesBeforeStart+1, waitChan(t, offsets))

		waitErr := make(chan error)
		go func() {
			waitErr <- waiter.WaitNewestMarked(ctx)
		}()

		require.EqualValues(t, messagesBeforeStart+2, waitChan(t, offsets))
		assertChanBlocked(t, waitErr)
		require.EqualValues(t, messagesBeforeStart+3, waitChan(t, offsets))
		time.Sleep(time.Second) // to be sure waiter does the job
		assertChanBlocked(t, waitErr)
		require.EqualValues(t, messagesBeforeStart+4, waitChan(t, offsets))
		require.NoError(t, waitChan(t, waitErr))
		assertChanBlocked(t, offsets)
	})

	t.Run("all messages consumed", func(t *testing.T) {
		// must be non-parallel

		waitErr := make(chan error)
		go func() {
			waitErr <- waiter.WaitNewestMarked(ctx)
		}()

		require.NoError(t, waitChan(t, waitErr))
		assertChanBlocked(t, offsets)
	})

	close(stop)
	require.NoError(t, waitChan(t, consumeErr))
}

func TestWaitingSession(t *testing.T) {
	t.Parallel()
	t.Run("no need to wait", func(t *testing.T) {
		t.Parallel()
		testInitial := func(t *testing.T, expected, initial int64) {
			t.Parallel()
			const topic = "topic"
			const partition = 4
			session := mocks.NewConsumerGroupSession(t)
			offsets := map[string]map[int32]int64{"A": {1: 5, 2: 6}, topic: {3: 7, partition: initial, 5: 9}}
			waiter := kafka.NewWaitingSession(logster.Discard, session, offsets)
			err := make(chan error)
			go func() {
				err <- waiter.WaitMarkOffset(context.Background(), topic, partition, expected)
			}()
			require.NoError(t, waitChan(t, err))
		}
		testCurrent := func(t *testing.T, initial, expected, current int64) {
			t.Parallel()
			const topic = "topic"
			const partition = 4
			session := mocks.NewConsumerGroupSession(t)
			offsets := map[string]map[int32]int64{"A": {1: 5, 2: 6}, topic: {3: 7, partition: initial, 5: 9}}
			waiter := kafka.NewWaitingSession(logster.Discard, session, offsets)
			session.EXPECT().MarkOffset(topic, int32(partition), current, "m").Once()
			waiter.MarkOffset(topic, partition, current, "m")
			err := make(chan error)
			go func() {
				err <- waiter.WaitMarkOffset(context.Background(), topic, partition, expected)
			}()
			require.NoError(t, waitChan(t, err))
		}
		t.Run("expected = initial", func(t *testing.T) { testInitial(t, 8, 8) })
		t.Run("expected < initial", func(t *testing.T) { testInitial(t, 7, 8) })
		t.Run("initial < expected = current", func(t *testing.T) { testCurrent(t, 8, 9, 9) })
		t.Run("initial < expected < current", func(t *testing.T) { testCurrent(t, 8, 9, 10) })
	})
	t.Run("wait", func(t *testing.T) {
		t.Parallel()
		t.Run("MarkOffset", func(t *testing.T) {
			t.Parallel()
			t.Run("offset = expected", func(t *testing.T) {
				t.Parallel()
				const topic = "topic"
				const partition = 4
				const initial = 8
				const expected = 10
				const offset = expected
				session := mocks.NewConsumerGroupSession(t)
				offsets := map[string]map[int32]int64{"A": {1: 5, 2: 6}, topic: {3: 7, partition: initial, 5: 9}}
				waiter := kafka.NewWaitingSession(logster.Discard, session, offsets)
				err := make(chan error)
				go func() {
					err <- waiter.WaitMarkOffset(context.Background(), topic, partition, expected)
				}()

				assertChanBlocked(t, err)

				session.EXPECT().MarkOffset("A", int32(2), int64(20), "x").Once()
				waiter.MarkOffset("A", 2, 20, "x") // different topic

				assertChanBlocked(t, err)

				session.EXPECT().MarkOffset("X", int32(2), int64(20), "y").Once()
				waiter.MarkOffset("X", 2, 20, "y") // unknown topic

				assertChanBlocked(t, err)

				session.EXPECT().MarkOffset("A", int32(4), int64(20), "z").Once()
				waiter.MarkOffset("A", 4, 20, "z") // unknown partition

				assertChanBlocked(t, err)

				session.EXPECT().MarkOffset(topic, int32(partition), int64(2), "").Once()
				waiter.MarkOffset(topic, partition, 2, "") // offset < expected

				assertChanBlocked(t, err)

				session.EXPECT().MarkOffset(topic, int32(partition), int64(9), "").Once()
				waiter.MarkOffset(topic, partition, 9, "") // offset < expected

				assertChanBlocked(t, err)

				session.EXPECT().MarkOffset(topic, int32(partition), int64(offset), "").Once()
				waiter.MarkOffset(topic, partition, offset, "") // offset = expected

				require.NoError(t, waitChan(t, err))
			})
			t.Run("offset > expected", func(t *testing.T) {
				t.Parallel()
				const topic = "topic"
				const partition = 4
				const initial = 8
				const expected = 10
				const offset = 20
				session := mocks.NewConsumerGroupSession(t)
				offsets := map[string]map[int32]int64{"A": {1: 5, 2: 6}, topic: {3: 7, partition: initial, 5: 9}}
				waiter := kafka.NewWaitingSession(logster.Discard, session, offsets)
				err := make(chan error)
				go func() {
					err <- waiter.WaitMarkOffset(context.Background(), topic, partition, expected)
				}()

				session.EXPECT().MarkOffset(topic, int32(partition), int64(offset), "").Once()
				waiter.MarkOffset(topic, partition, offset, "")

				require.NoError(t, waitChan(t, err))
			})
		})
		t.Run("ResetOffset", func(t *testing.T) {
			t.Parallel()
			t.Run("offset = expected", func(t *testing.T) {
				t.Parallel()
				const topic = "topic"
				const partition = 4
				const initial = 8
				const expected = 10
				const offset = expected
				session := mocks.NewConsumerGroupSession(t)
				offsets := map[string]map[int32]int64{"A": {1: 5, 2: 6}, topic: {3: 7, partition: initial, 5: 9}}
				waiter := kafka.NewWaitingSession(logster.Discard, session, offsets)
				err := make(chan error)
				go func() {
					err <- waiter.WaitMarkOffset(context.Background(), topic, partition, expected)
				}()

				assertChanBlocked(t, err)

				session.EXPECT().ResetOffset("A", int32(2), int64(20), "x").Once()
				waiter.ResetOffset("A", 2, 20, "x") // different topic

				assertChanBlocked(t, err)

				session.EXPECT().ResetOffset("X", int32(2), int64(20), "y").Once()
				waiter.ResetOffset("X", 2, 20, "y") // unknown topic

				assertChanBlocked(t, err)

				session.EXPECT().ResetOffset("A", int32(4), int64(20), "z").Once()
				waiter.ResetOffset("A", 4, 20, "z") // unknown partition

				assertChanBlocked(t, err)

				session.EXPECT().ResetOffset(topic, int32(partition), int64(2), "").Once()
				waiter.ResetOffset(topic, partition, 2, "") // offset < expected

				assertChanBlocked(t, err)

				session.EXPECT().ResetOffset(topic, int32(partition), int64(9), "").Once()
				waiter.ResetOffset(topic, partition, 9, "") // offset < expected

				assertChanBlocked(t, err)

				session.EXPECT().ResetOffset(topic, int32(partition), int64(offset), "").Once()
				waiter.ResetOffset(topic, partition, offset, "") // offset = expected

				require.NoError(t, waitChan(t, err))
			})
			t.Run("offset > expected", func(t *testing.T) {
				t.Parallel()
				const topic = "topic"
				const partition = 4
				const initial = 8
				const expected = 10
				const offset = 20
				session := mocks.NewConsumerGroupSession(t)
				offsets := map[string]map[int32]int64{"A": {1: 5, 2: 6}, topic: {3: 7, partition: initial, 5: 9}}
				waiter := kafka.NewWaitingSession(logster.Discard, session, offsets)
				err := make(chan error)
				go func() {
					err <- waiter.WaitMarkOffset(context.Background(), topic, partition, expected)
				}()

				session.EXPECT().ResetOffset(topic, int32(partition), int64(offset), "").Once()
				waiter.ResetOffset(topic, partition, offset, "")

				require.NoError(t, waitChan(t, err))
			})
		})
		t.Run("MarkMessage", func(t *testing.T) {
			t.Parallel()
			t.Run("offset+1 = expected", func(t *testing.T) {
				t.Parallel()
				const topic = "topic"
				const partition = 4
				const initial = 8
				const expected = 10
				const offset = 9
				session := mocks.NewConsumerGroupSession(t)
				offsets := map[string]map[int32]int64{"A": {1: 5, 2: 6}, topic: {3: 7, partition: initial, 5: 9}}
				waiter := kafka.NewWaitingSession(logster.Discard, session, offsets)
				err := make(chan error)
				go func() {
					err <- waiter.WaitMarkOffset(context.Background(), topic, partition, expected)
				}()

				msg := &sarama.ConsumerMessage{
					Headers:        nil,
					Timestamp:      time.Now(),
					BlockTimestamp: time.Now(),
					Key:            sarama.ByteEncoder("key"),
					Value:          sarama.ByteEncoder("value"),
					Topic:          topic,
					Partition:      partition,
					Offset:         offset,
				}
				session.EXPECT().MarkMessage(msg, "x").Once()
				waiter.MarkMessage(msg, "x")

				require.NoError(t, waitChan(t, err))
			})
			t.Run("offset+1 > expected", func(t *testing.T) {
				t.Parallel()
				const topic = "topic"
				const partition = 4
				const initial = 8
				const expected = 10
				const offset = 20
				session := mocks.NewConsumerGroupSession(t)
				offsets := map[string]map[int32]int64{"A": {1: 5, 2: 6}, topic: {3: 7, partition: initial, 5: 9}}
				waiter := kafka.NewWaitingSession(logster.Discard, session, offsets)
				err := make(chan error)
				go func() {
					err <- waiter.WaitMarkOffset(context.Background(), topic, partition, expected)
				}()

				msg := &sarama.ConsumerMessage{
					Headers:        nil,
					Timestamp:      time.Now(),
					BlockTimestamp: time.Now(),
					Key:            sarama.ByteEncoder("key"),
					Value:          sarama.ByteEncoder("value"),
					Topic:          topic,
					Partition:      partition,
					Offset:         offset,
				}
				session.EXPECT().MarkMessage(msg, "x").Once()
				waiter.MarkMessage(msg, "x")

				require.NoError(t, waitChan(t, err))
			})
		})
	})
	t.Run("unknown topic", func(t *testing.T) {
		t.Parallel()
		session := mocks.NewConsumerGroupSession(t)
		offsets := map[string]map[int32]int64{"A": {1: 5, 2: 6}, "B": {3: 7, 4: 8, 5: 9}}
		waiter := kafka.NewWaitingSession(logster.Discard, session, offsets)
		err := make(chan error)
		go func() {
			err <- waiter.WaitMarkOffset(context.Background(), "X", 2, 10)
		}()
		require.EqualError(t, waitChan(t, err), `unknown topic "X"`)
	})
	t.Run("unknown partition", func(t *testing.T) {
		t.Parallel()
		session := mocks.NewConsumerGroupSession(t)
		offsets := map[string]map[int32]int64{"A": {1: 5, 2: 6}, "B": {3: 7, 4: 8, 5: 9}}
		waiter := kafka.NewWaitingSession(logster.Discard, session, offsets)
		err := make(chan error)
		go func() {
			err <- waiter.WaitMarkOffset(context.Background(), "A", 3, 10)
		}()
		require.EqualError(t, waitChan(t, err), `unknown partition "A"/3`)
	})
}

func waitChan[T any](t *testing.T, ch <-chan T) T {
	select {
	case v := <-ch:
		return v
	case <-time.After(5 * time.Second):
		require.Fail(t, "timeout")
		panic("")
	}
}

func assertChanBlocked[T any](t *testing.T, ch <-chan T) {
	select {
	case v := <-ch:
		require.Failf(t, "chan not blocked", "%+v", v)
	default:
	}
}
