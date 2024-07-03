package grpcstreams

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/store/memory"
	"github.com/w-h-a/pkg/streams"
)

func TestGrpcStreams(t *testing.T) {
	s := NewStreams(GrpcStreamsWithStore(memory.NewStore()))
	require.NotNil(t, s)

	runTestStream(t, s)
}

func runTestStream(t *testing.T, s streams.Streams) {
	t.Run("Test Consume Topic AutoAck", func(t *testing.T) {
		payload := map[string]string{"message": "hello world"}
		metadata := map[string]string{"foo": "bar"}
		id := uuid.New().String()
		topic := "test1"

		err := s.Subscribe(
			id,
			streams.SubscribeWithTopic(topic),
			streams.SubscribeWithAutoAck(true, 0),
		)
		require.NoError(t, err)

		defer func() {
			err = s.Unsubscribe(id)
			require.NoError(t, err)
		}()

		evChan, err := s.Consume(id)
		require.NoError(t, err)

		timeCh := make(chan time.Time)
		doneCh := make(chan bool)

		go func() {
			timeout := time.NewTicker(time.Second * 4)

			select {
			case event := <-evChan:
				require.NotNil(t, event)
				require.Equal(t, metadata, event.Metadata)

				result := map[string]string{}

				err := event.Unmarshal(&result)
				require.NoError(t, err)
				require.Equal(t, payload, result)

				doneCh <- true
			case t := <-timeout.C:
				timeCh <- t
			}
		}()

		err = s.Produce(
			topic,
			payload,
			streams.ProduceWithMetadata(metadata),
		)
		require.NoError(t, err)

		select {
		case <-timeCh:
			t.Fatal("event was not received")
		case <-doneCh:
		}
	})

	t.Run("Test Consume Group AutoAck", func(t *testing.T) {
		payload := map[string]string{"message": "hello world"}
		metadata := map[string]string{"foo": "bar"}
		id1 := uuid.New().String()
		id2 := uuid.New().String()
		group1 := uuid.New().String()
		group2 := uuid.New().String()
		topic := "test2"

		err := s.Subscribe(
			id1,
			streams.SubscribeWithGroup(group1),
			streams.SubscribeWithTopic(topic),
			streams.SubscribeWithAutoAck(true, 0),
		)
		require.NoError(t, err)

		defer func() {
			err = s.Unsubscribe(id1)
			require.NoError(t, err)
		}()

		evChan1, err := s.Consume(id1)
		require.NoError(t, err)

		timeCh1 := make(chan time.Time)
		doneCh1 := make(chan bool)

		go func() {
			timeout := time.NewTicker(time.Second * 4)

			select {
			case event := <-evChan1:
				require.NotNil(t, event)
				require.Equal(t, metadata, event.Metadata)

				result := map[string]string{}

				err := event.Unmarshal(&result)
				require.NoError(t, err)
				require.Equal(t, payload, result)

				doneCh1 <- true
			case t := <-timeout.C:
				timeCh1 <- t
			}
		}()

		err = s.Produce(
			topic,
			payload,
			streams.ProduceWithMetadata(metadata),
		)
		require.NoError(t, err)

		err = s.Subscribe(
			id2,
			streams.SubscribeWithGroup(group2),
			streams.SubscribeWithTopic(topic),
			streams.SubscribeWithAutoAck(true, 0),
			streams.SubscribeWithOffset(time.Now().Add(time.Minute*-1)),
		)
		require.NoError(t, err)

		defer func() {
			err = s.Unsubscribe(id2)
			require.NoError(t, err)
		}()

		evChan2, err := s.Consume(id2)
		require.NoError(t, err)

		timeCh2 := make(chan time.Time)
		doneCh2 := make(chan bool)

		go func() {
			timeout := time.NewTicker(time.Second * 4)

			select {
			case event := <-evChan2:
				require.NotNil(t, event)
				require.Equal(t, metadata, event.Metadata)

				result := map[string]string{}

				err := event.Unmarshal(&result)
				require.NoError(t, err)
				require.Equal(t, payload, result)

				doneCh2 <- true
			case t := <-timeout.C:
				timeCh2 <- t
			}
		}()

		select {
		case <-timeCh1:
			t.Fatal("first subscriber did not receive event")
		case <-doneCh1:
		}

		select {
		case <-timeCh2:
			t.Fatal("second subscriber did not receive event")
		case <-doneCh2:
		}
	})

	t.Run("Test Consume With Manual Acking", func(t *testing.T) {
		payload1 := map[string]string{"message": "hello world 1"}
		payload2 := map[string]string{"message": "hello world 2"}
		id := uuid.New().String()
		topic := "test3"

		err := s.Subscribe(
			id,
			streams.SubscribeWithTopic(topic),
		)
		require.NoError(t, err)

		defer func() {
			err = s.Unsubscribe(id)
			require.NoError(t, err)
		}()

		evChan, err := s.Consume(id)
		require.NoError(t, err)

		require.NoError(t, s.Produce(topic, payload1))
		require.NoError(t, s.Produce(topic, payload2))

		ev := <-evChan
		ev.Ack()

		ev = <-evChan
		ev.Nack()

		nackedId := ev.Id

		select {
		case ev = <-evChan:
			require.Equal(t, nackedId, ev.Id)
			require.NoError(t, ev.Ack())
		case <-time.After(8 * time.Second):
			t.Fatal("timed out waiting for event to be put back on the queue")
		}
	})

	t.Run("Test Consume With Multiple Topics", func(t *testing.T) {
		payload := map[string]string{"message": "hello world"}
		id1 := uuid.New().String()
		id2 := uuid.New().String()
		topic1 := "test4"
		topic2 := "test5"

		err := s.Subscribe(
			id1,
			streams.SubscribeWithTopic(topic1),
		)
		require.NoError(t, err)

		defer func() {
			err = s.Unsubscribe(id1)
			require.NoError(t, err)
		}()

		evChan1, err := s.Consume(id1)
		require.NoError(t, err)

		err = s.Subscribe(
			id2,
			streams.SubscribeWithTopic(topic2),
		)
		require.NoError(t, err)

		defer func() {
			err = s.Unsubscribe(id2)
			require.NoError(t, err)
		}()

		evChan2, err := s.Consume(id2)
		require.NoError(t, err)

		require.NoError(t, s.Produce(topic1, payload))
		require.NoError(t, s.Produce(topic2, payload))

		wg := &sync.WaitGroup{}

		wg.Add(2)

		go func() {
			ev := <-evChan1
			require.Equal(t, topic1, ev.Topic)
			ev.Ack()
			wg.Done()
		}()

		go func() {
			ev := <-evChan2
			ev.Ack()
			require.Equal(t, topic2, ev.Topic)
			ev.Ack()
			wg.Done()
		}()

		wg.Wait()
	})
}
