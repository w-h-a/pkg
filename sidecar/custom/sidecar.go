package custom

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/w-h-a/pkg/broker"
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/sidecar"
	"github.com/w-h-a/pkg/store"
	"github.com/w-h-a/pkg/telemetry/log"
)

type customSidecar struct {
	options     sidecar.SidecarOptions
	subscribers map[string]broker.Subscriber
	mtx         sync.RWMutex
}

func (s *customSidecar) Options() sidecar.SidecarOptions {
	return s.options
}

func (s *customSidecar) OnEventPublished(event *sidecar.Event) error {
	var err error

	if len(event.To) > 0 {
		err = s.actOnEventFromApp(event)
	} else {
		err = s.postEventToApp(event)
	}

	return err
}

func (s *customSidecar) SaveStateToStore(state *sidecar.State) error {
	if len(state.Records) == 0 {
		return nil
	}

	log.Infof("RECEIVED STATE %+v", state)

	log.Infof("ID", state.StoreId)

	log.Infof("RECORDS %+v", state.Records)

	st, ok := s.options.Stores[state.StoreId]
	if !ok {
		return nil
	}

	for _, record := range state.Records {
		if err := st.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func (s *customSidecar) RetrieveStateFromStore(storeId, key string) ([]*store.Record, error) {
	st, ok := s.options.Stores[storeId]
	if !ok {
		return nil, nil
	}

	recs, err := st.Read(key, store.ReadWithPrefix())
	if err != nil {
		return nil, err
	}

	return recs, nil
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

		return s.postEventToApp(event)
	}, bk.Options().SubscribeOptions)

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

func (s *customSidecar) actOnEventFromApp(event *sidecar.Event) error {
	if event.State != nil && len(event.State.Records) > 0 {
		if err := s.SaveStateToStore(event.State); err != nil {
			return err
		}
	}

	for _, target := range event.To {
		if err := s.sendEventToTarget(target, event); err != nil {
			return err
		}
		log.Infof("successfully sent event %s to target %s", event.EventName, target)
	}

	return nil
}

func (s *customSidecar) sendEventToTarget(target string, event *sidecar.Event) error {
	bk, ok := s.options.Brokers[target]
	if ok {
		if err := bk.Publish(event.Data, bk.Options().PublishOptions); err != nil {
			return fmt.Errorf("failed to send event %s to target %s: %v", event.EventName, target, err)
		}
	} else {
		name := fmt.Sprintf("%s-action", target)

		url := fmt.Sprintf("%s:%s", name, s.options.HttpPort.Port)

		newEvent := &sidecar.Event{
			EventName: event.EventName,
			Data:      event.Data,
			CreatedAt: time.Now(),
		}

		if _, err := s.sendEventViaHttp(name, name, s.options.HttpPort.Port, "publish", url, newEvent); err != nil {
			return err
		}
	}

	return nil
}

func (s *customSidecar) postEventToApp(event *sidecar.Event) error {
	var rsp *sidecar.Event
	var err error

	url := fmt.Sprintf("%s:%s", s.options.ServiceName, s.options.ServicePort.Port)

	if s.options.ServicePort.Protocol == "rpc" {
		// TODO: refactor
		serviceTitle := strings.Title(s.options.ServiceName)

		eventTitle := strings.Title(event.EventName)

		// TODO: not a great assumption
		method := fmt.Sprintf("%s.%s", serviceTitle, eventTitle)

		rsp, err = s.sendEventViaRpc(s.options.ServiceName, s.options.ServiceName, s.options.ServicePort.Port, method, url, event)
		if err != nil {
			return err
		}
	} else {
		rsp, err = s.sendEventViaHttp(s.options.ServiceName, s.options.ServiceName, s.options.ServicePort.Port, event.EventName, url, event)
		if err != nil {
			return err
		}
	}

	if rsp != nil {
		if err := s.actOnEventFromApp(rsp); err != nil {
			return err
		}
	}

	return nil
}

func (s *customSidecar) sendEventViaHttp(namespace, name, port, endpoint, baseUrl string, event *sidecar.Event) (*sidecar.Event, error) {
	p, _ := strconv.Atoi(port)

	req := s.options.HttpClient.NewRequest(
		client.RequestWithNamespace(namespace),
		client.RequestWithName(name),
		client.RequestWithPort(p),
		client.RequestWithMethod(endpoint),
		client.RequestWithUnmarshaledRequest(event),
	)

	rsp := &sidecar.Event{}

	if err := s.options.HttpClient.Call(context.Background(), req, rsp, client.CallWithAddress(baseUrl)); err != nil {
		log.Errorf("RECEIVED ERR %+v", err)
		return nil, err
	}

	log.Infof("RESPONSE FROM APP %+v", rsp.State)

	return rsp, nil
}

func (s *customSidecar) sendEventViaRpc(namespace, name, port, method, baseUrl string, event *sidecar.Event) (*sidecar.Event, error) {
	p, _ := strconv.Atoi(port)

	req := s.options.RpcClient.NewRequest(
		client.RequestWithNamespace(namespace),
		client.RequestWithName(name),
		client.RequestWithPort(p),
		client.RequestWithMethod(method),
		client.RequestWithUnmarshaledRequest(event),
	)

	rsp := &sidecar.Event{}

	if err := s.options.RpcClient.Call(context.Background(), req, rsp, client.CallWithAddress(baseUrl)); err != nil {
		return nil, err
	}

	return rsp, nil
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
