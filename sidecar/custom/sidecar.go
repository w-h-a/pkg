package custom

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/w-h-a/pkg/broker"
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/sidecar"
	"github.com/w-h-a/pkg/store"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/telemetry/trace"
	"github.com/w-h-a/pkg/telemetry/trace/memory"
	"github.com/w-h-a/pkg/utils/datautils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	bk, ok := s.options.Brokers[event.EventName]
	if !ok {
		log.Warnf("broker %s was not found", event.EventName)
		return sidecar.ErrComponentNotFound
	}

	if err := bk.Publish(event.Data, *bk.Options().PublishOptions); err != nil {
		return err
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
			EventName: brokerId,
			Data:      body,
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

func (s *customSidecar) ReadFromSecretStore(ctx context.Context, secretStore string, name string) (*sidecar.Secret, error) {
	tracer := trace.GetTracer()

	if tracer == nil {
		log.Error("setting default memory tracer")
		tracer = memory.NewTrace()
		trace.SetTracer(tracer)
	}

	_, span, err := tracer.Start(
		ctx,
		"customSidecar.ReadFromSecretStore",
		map[string]string{
			"secretStore": secretStore,
			"name":        name,
			"error":       "",
		},
	)
	if err != nil {
		log.Errorf("failed to start span: %v", err)
		return nil, trace.ErrStart
	}

	sc, ok := s.options.Secrets[secretStore]
	if !ok {
		log.Warnf("secret store %s was not found", secretStore)
		span.Metadata["error"] = fmt.Sprintf("secret store %s was not found", secretStore)
		tracer.Finish(span)
		return nil, sidecar.ErrComponentNotFound
	}

	mp, err := sc.GetSecret(name)
	if err != nil {
		span.Metadata["error"] = err.Error()
		tracer.Finish(span)
		return nil, err
	}

	tracer.Finish(span)

	return &sidecar.Secret{
		Data: mp,
	}, nil
}

func (s *customSidecar) String() string {
	return "custom"
}

func (s *customSidecar) sendEventToService(event *sidecar.Event) error {
	url := fmt.Sprintf("%s:%s", s.options.ServiceName, s.options.ServicePort.Port)

	p, _ := strconv.Atoi(s.options.ServicePort.Port)

	opts := []client.RequestOption{
		client.RequestWithNamespace(s.options.ServiceName),
		client.RequestWithName(s.options.ServiceName),
		client.RequestWithPort(p),
	}

	parts := strings.Split(event.EventName, "-")

	if len(parts) != 2 {
		return sidecar.ErrInvalidGroupName
	}

	if s.options.ServicePort.Protocol == "grpc" {
		caser := cases.Title(language.English)
		receiver := caser.String(parts[0])
		method := caser.String(parts[1])
		pbEvent, err := sidecar.SerializeEvent(event)
		if err != nil {
			return err
		}
		opts = append(
			opts,
			client.RequestWithMethod(fmt.Sprintf("%s.%s", receiver, method)),
			client.RequestWithUnmarshaledRequest(pbEvent),
		)
	} else {
		receiver := strings.ToLower(parts[0])
		method := strings.ToLower(parts[1])
		opts = append(
			opts,
			client.RequestWithMethod(fmt.Sprintf("%s/%s", receiver, method)),
			client.RequestWithUnmarshaledRequest(event),
		)
	}

	req := s.options.Client.NewRequest(opts...)

	var rsp interface{}

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
