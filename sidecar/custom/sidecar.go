package custom

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/w-h-a/pkg/broker"
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/sidecar"
	"github.com/w-h-a/pkg/store"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/utils/datautils"
)

type customSidecar struct {
	options     sidecar.SidecarOptions
	subscribers map[string]broker.Subscriber
	mtx         sync.RWMutex
}

func (s *customSidecar) Options() sidecar.SidecarOptions {
	return s.options
}

func (s *customSidecar) SaveStateToStore(state *sidecar.State) error {
	if len(state.Records) == 0 {
		return nil
	}

	st, ok := s.options.Stores[state.StoreId]
	if !ok {
		log.Warnf("store %s was not found", state.StoreId)
		return sidecar.ErrComponentNotFound
	}

	for _, record := range state.Records {
		storeRecord := &store.Record{
			Key: record.Key,
		}

		data := record.Value

		bs, err := datautils.Stringify(data)
		if err != nil {
			return err
		}

		storeRecord.Value = bs

		if err := st.Write(storeRecord); err != nil {
			return err
		}
	}

	return nil
}

func (c *customSidecar) ListStateFromStore(storeId string) ([]*store.Record, error) {
	st, ok := c.options.Stores[storeId]
	if !ok {
		log.Warnf("store %s was not found", storeId)
		return nil, sidecar.ErrComponentNotFound
	}

	// TODO: limit + offset
	recs, err := st.Read("", store.ReadWithPrefix())
	if err != nil {
		return nil, err
	}

	return recs, nil
}

func (s *customSidecar) SingleStateFromStore(storeId, key string) ([]*store.Record, error) {
	st, ok := s.options.Stores[storeId]
	if !ok {
		log.Warnf("store %s was not found", storeId)
		return nil, sidecar.ErrComponentNotFound
	}

	recs, err := st.Read(key)
	if err != nil {
		return nil, err
	}

	return recs, nil
}

func (s *customSidecar) RemoveStateFromStore(storeId, key string) error {
	st, ok := s.options.Stores[storeId]
	if !ok {
		log.Warnf("store %s was not found", storeId)
		return sidecar.ErrComponentNotFound
	}

	if err := st.Delete(key); err != nil {
		return err
	}

	return nil
}

func (s *customSidecar) WriteEventToBroker(event *sidecar.Event) error {
	if len(event.To) == 0 {
		log.Warnf("event %#+event has no address", event)
		return nil
	}

	if len(event.Concurrent) > 0 {
		s.sendEventToTargetsConcurrently(event)
	} else {
		if err := s.sendEventToTargetsSequentially(event); err != nil {
			return err
		}
	}

	return nil
}

func (s *customSidecar) ReadEventsFromBroker(brokerId string) {
	bk, ok := s.options.Brokers[brokerId]
	if !ok {
		log.Warnf("broker %s was not found", brokerId)
		return
	}

	s.mtx.RLock()

	_, ok = s.subscribers[brokerId]
	if ok {
		log.Warnf("a subscriber for broker %s was already found", brokerId)
		s.mtx.RUnlock()
		return
	}

	s.mtx.RUnlock()

	sub := bk.Subscribe(func(b []byte) error {
		var body interface{}
		if err := json.Unmarshal(b, &body); err != nil {
			return err
		}

		event := &sidecar.Event{
			Data:      body,
			EventName: brokerId,
			CreatedAt: time.Now(),
		}

		return s.sendEventToService(event)
	}, *bk.Options().SubscribeOptions)

	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.subscribers[brokerId] = sub
}

func (s *customSidecar) UnsubscribeFromBroker(brokerId string) error {
	s.mtx.RLock()

	sub, ok := s.subscribers[brokerId]
	if !ok {
		s.mtx.RUnlock()
		return nil
	}

	s.mtx.RUnlock()

	if err := sub.Unsubscribe(); err != nil {
		return err
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.subscribers, brokerId)

	return nil
}

func (s *customSidecar) String() string {
	return "custom"
}

func (s *customSidecar) sendEventToTargetsConcurrently(event *sidecar.Event) {
	for _, target := range event.To {
		go func() {
			err := s.sendEventToTarget(target, event)
			if err != nil {
				log.Errorf("failed to send event %s to target %s: %v", event.EventName, target, err)
			}
		}()
	}
}

func (s *customSidecar) sendEventToTargetsSequentially(event *sidecar.Event) error {
	for _, target := range event.To {
		err := s.sendEventToTarget(target, event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *customSidecar) sendEventToTarget(target string, event *sidecar.Event) error {
	bk, ok := s.options.Brokers[target]
	if !ok {
		log.Warnf("broker %s was not found", target)
		return sidecar.ErrComponentNotFound
	}

	if err := bk.Publish(event.Data, *bk.Options().PublishOptions); err != nil {
		return err
	}

	return nil
}

func (s *customSidecar) sendEventToService(event *sidecar.Event) error {
	url := fmt.Sprintf("%s:%s", s.options.ServiceName, s.options.ServicePort.Port)

	p, _ := strconv.Atoi(s.options.ServicePort.Port)

	req := s.options.Client.NewRequest(
		client.RequestWithNamespace(s.options.ServiceName),
		client.RequestWithName(s.options.ServiceName),
		client.RequestWithPort(p),
		// TODO: make this better
		client.RequestWithMethod(event.EventName),
		// TODO: does the service accept proto?
		client.RequestWithUnmarshaledRequest(event),
	)

	rsp := &sidecar.Event{}

	if err := s.options.Client.Call(context.Background(), req, rsp, client.CallWithAddress(url)); err != nil {
		return err
	}

	return nil
}

func NewSidecar(opts ...sidecar.SidecarOption) sidecar.Sidecar {
	options := sidecar.NewSidecarOptions(opts...)

	s := &customSidecar{
		options:     options,
		subscribers: map[string]broker.Subscriber{},
		mtx:         sync.RWMutex{},
	}

	return s
}
