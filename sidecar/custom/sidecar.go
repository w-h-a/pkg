package custom

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/sidecar"
	"github.com/w-h-a/pkg/store"
	"github.com/w-h-a/pkg/telemetry/log"
)

type customSidecar struct {
	options sidecar.SidecarOptions
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

func (s *customSidecar) SaveStateToStore(storeId string, state []*store.Record) error {
	if len(state) == 0 {
		return nil
	}

	st, ok := s.options.Stores[storeId]
	if !ok {
		return nil
	}

	for _, record := range state {
		key := fmt.Sprintf("%s:%s", s.options.Id, record.Key)
		if err := st.Write(&store.Record{
			Key:   key,
			Value: record.Value,
		}); err != nil {
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

	prefix := fmt.Sprintf("%s:%s", s.options.Id, key)

	recs, err := st.Read(prefix, store.ReadWithPrefix())
	if err != nil {
		return nil, err
	}

	return recs, nil
}

func (s *customSidecar) ReadEventsFromBroker(brokerId, eventName string) {
	bk, ok := s.options.Brokers[brokerId]
	if !ok {
		log.Warnf("broker %s was not found", brokerId)
		return
	}

	go func() {
		bk.Subscribe(func(b []byte) error {
			var body interface{}
			if err := json.Unmarshal(b, &body); err != nil {
				return err
			}

			event := &sidecar.Event{
				Data:      body,
				EventName: eventName,
				CreatedAt: time.Now(),
			}

			return s.postEventToApp(event)
		}, bk.Options().SubscribeOptions)
	}()
}

func (s *customSidecar) String() string {
	return "custom"
}

func (s *customSidecar) actOnEventFromApp(event *sidecar.Event) error {
	if event.State != nil && len(event.State.Records) > 0 {
		if err := s.SaveStateToStore(event.State.StoreId, event.State.Records); err != nil {
			return err
		}
	}

	for _, target := range event.To {
		if err := s.sendEventToTarget(target, event); err != nil {
			return err
		}
		log.Infof("successfully send event %s to target %s", event.EventName, target)
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

		url := fmt.Sprintf("http://%s:%s", name, s.options.HttpPort.Port)

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
	url := fmt.Sprintf("http://%s:%s", s.options.ServiceName, s.options.ServicePort)

	rsp, err := s.sendEventViaHttp(s.options.ServiceName, s.options.ServiceName, s.options.ServicePort.Port, event.EventName, url, event)
	if err != nil {
		return err
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
		return nil, err
	}

	return rsp, nil
}

func NewSidecar(opts ...sidecar.SidecarOption) sidecar.Sidecar {
	options := sidecar.NewSidecarOptions(opts...)

	s := &customSidecar{
		options: options,
	}

	return s
}
