package custom

import (
	"context"
	"encoding/hex"
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
	"github.com/w-h-a/pkg/telemetry/tracev2"
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

func (s *customSidecar) SaveStateToStore(ctx context.Context, state *sidecar.State) error {
	_, spanId := s.options.Tracer.Start(ctx, "customSidecar.SaveStateToStore")
	defer s.options.Tracer.Finish(spanId)

	records, _ := json.Marshal(state.Records)

	s.options.Tracer.AddMetadata(spanId, map[string]string{
		"storeId": state.StoreId,
		"records": string(records),
	})

	if len(state.Records) == 0 {
		s.options.Tracer.UpdateStatus(spanId, 2, "success")
		return nil
	}

	st, ok := s.options.Stores[state.StoreId]
	if !ok {
		log.Warnf("store %s was not found", state.StoreId)
		s.options.Tracer.UpdateStatus(spanId, 1, fmt.Sprintf("store %s was not found", state.StoreId))
		return sidecar.ErrComponentNotFound
	}

	for _, record := range state.Records {
		storeRecord := &store.Record{
			Key: record.Key,
		}

		data := record.Value

		bs, err := datautils.Stringify(data)
		if err != nil {
			s.options.Tracer.UpdateStatus(spanId, 1, err.Error())
			return err
		}

		storeRecord.Value = bs

		if err := st.Write(storeRecord); err != nil {
			s.options.Tracer.UpdateStatus(spanId, 1, err.Error())
			return err
		}
	}

	s.options.Tracer.UpdateStatus(spanId, 2, "success")

	return nil
}

func (s *customSidecar) ListStateFromStore(ctx context.Context, storeId string) ([]*store.Record, error) {
	_, spanId := s.options.Tracer.Start(ctx, "customSidecar.ListStateFromStore")
	defer s.options.Tracer.Finish(spanId)

	s.options.Tracer.AddMetadata(spanId, map[string]string{
		"storeId": storeId,
	})

	st, ok := s.options.Stores[storeId]
	if !ok {
		log.Warnf("store %s was not found", storeId)
		s.options.Tracer.UpdateStatus(spanId, 1, fmt.Sprintf("store %s was not found", storeId))
		return nil, sidecar.ErrComponentNotFound
	}

	// TODO: limit + offset
	recs, err := st.Read("", store.ReadWithPrefix())
	if err != nil {
		s.options.Tracer.UpdateStatus(spanId, 1, err.Error())
		return nil, err
	}

	s.options.Tracer.UpdateStatus(spanId, 2, "success")

	return recs, nil
}

func (s *customSidecar) SingleStateFromStore(ctx context.Context, storeId, key string) ([]*store.Record, error) {
	_, spanId := s.options.Tracer.Start(ctx, "customSidecar.SingleStateFromStore")
	defer s.options.Tracer.Finish(spanId)

	s.options.Tracer.AddMetadata(spanId, map[string]string{
		"storeId": storeId,
		"key":     key,
	})

	st, ok := s.options.Stores[storeId]
	if !ok {
		log.Warnf("store %s was not found", storeId)
		s.options.Tracer.UpdateStatus(spanId, 1, fmt.Sprintf("store %s was not found", storeId))
		return nil, sidecar.ErrComponentNotFound
	}

	recs, err := st.Read(key)
	if err != nil {
		s.options.Tracer.UpdateStatus(spanId, 1, err.Error())
		return nil, err
	}

	s.options.Tracer.UpdateStatus(spanId, 2, "success")

	return recs, nil
}

func (s *customSidecar) RemoveStateFromStore(ctx context.Context, storeId, key string) error {
	_, spanId := s.options.Tracer.Start(ctx, "customSidecar.RemoveStateFromStore")
	defer s.options.Tracer.Finish(spanId)

	s.options.Tracer.AddMetadata(spanId, map[string]string{
		"storeId": storeId,
		"key":     key,
	})

	st, ok := s.options.Stores[storeId]
	if !ok {
		log.Warnf("store %s was not found", storeId)
		s.options.Tracer.UpdateStatus(spanId, 1, fmt.Sprintf("store %s was not found", storeId))
		return sidecar.ErrComponentNotFound
	}

	if err := st.Delete(key); err != nil {
		s.options.Tracer.UpdateStatus(spanId, 1, err.Error())
		return err
	}

	s.options.Tracer.UpdateStatus(spanId, 2, "success")

	return nil
}

func (s *customSidecar) WriteEventToBroker(ctx context.Context, event *sidecar.Event) error {
	newCtx, spanId := s.options.Tracer.Start(ctx, "customSidecar.WriteEventToBroker")
	defer s.options.Tracer.Finish(spanId)

	if traceparent, found := tracev2.TraceParentFromContext(newCtx); found {
		if _, ok := event.Payload[tracev2.TraceParentKey].(string); !ok {
			event.Payload[tracev2.TraceParentKey] = hex.EncodeToString(traceparent[:])
		}
	}

	payload, _ := json.Marshal(event.Payload)

	s.options.Tracer.AddMetadata(spanId, map[string]string{
		"eventName": event.EventName,
		"payload":   string(payload),
	})

	bk, ok := s.options.Brokers[event.EventName]
	if !ok {
		log.Warnf("broker %s was not found", event.EventName)
		s.options.Tracer.UpdateStatus(spanId, 1, fmt.Sprintf("broker %s was not found", event.EventName))
		return sidecar.ErrComponentNotFound
	}

	if err := bk.Publish(event.Payload, *bk.Options().PublishOptions); err != nil {
		s.options.Tracer.UpdateStatus(spanId, 1, err.Error())
		return err
	}

	s.options.Tracer.UpdateStatus(spanId, 2, "success")

	return nil
}

func (s *customSidecar) ReadEventsFromBroker(ctx context.Context, brokerId string) {
	_, spanId := s.options.Tracer.Start(ctx, "customSidecar.ReadEventsFromBroker")
	defer s.options.Tracer.Finish(spanId)

	s.options.Tracer.AddMetadata(spanId, map[string]string{
		"brokerId": brokerId,
	})

	bk, ok := s.options.Brokers[brokerId]
	if !ok {
		log.Warnf("broker %s was not found", brokerId)
		s.options.Tracer.UpdateStatus(spanId, 1, fmt.Sprintf("broker %s was not found", brokerId))
		return
	}

	s.mtx.RLock()

	_, ok = s.subscribers[brokerId]
	if ok {
		log.Warnf("a subscriber for broker %s was already found", brokerId)
		s.mtx.RUnlock()
		s.options.Tracer.UpdateStatus(spanId, 1, fmt.Sprintf("a subscriber for broker %s was already found", brokerId))
		return
	}

	s.mtx.RUnlock()

	sub := bk.Subscribe(func(b []byte) error {
		var payload map[string]interface{}

		if err := json.Unmarshal(b, &payload); err != nil {
			s.options.Tracer.UpdateStatus(spanId, 1, err.Error())
			return err
		}

		var traceparent [16]byte

		var spanId string

		if encoded, ok := payload[tracev2.TraceParentKey].(string); ok {
			decoded, _ := hex.DecodeString(encoded)
			copy(traceparent[:], decoded)
			ctx, _ := tracev2.ContextWithTraceParent(context.Background(), traceparent)
			_, spanId = s.options.Tracer.Start(ctx, fmt.Sprintf("%s.Handler", brokerId))
		} else {
			_, spanId = s.options.Tracer.Start(context.Background(), fmt.Sprintf("%s.Handler", brokerId))
		}

		defer s.options.Tracer.Finish(spanId)

		s.options.Tracer.AddMetadata(spanId, map[string]string{
			"brokerId": brokerId,
			"payload":  string(b),
		})

		event := &sidecar.Event{
			EventName: brokerId,
			Payload:   payload,
		}

		s.options.Tracer.UpdateStatus(spanId, 2, "success")

		return s.sendEventToService(event)
	}, *bk.Options().SubscribeOptions)

	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.subscribers[brokerId] = sub

	s.options.Tracer.UpdateStatus(spanId, 2, "success")
}

func (s *customSidecar) UnsubscribeFromBroker(ctx context.Context, brokerId string) error {
	_, spanId := s.options.Tracer.Start(ctx, "customSidecar.UnsubscribeFromBroker")
	defer s.options.Tracer.Finish(spanId)

	s.options.Tracer.AddMetadata(spanId, map[string]string{
		"brokerId": brokerId,
	})

	s.mtx.RLock()

	sub, ok := s.subscribers[brokerId]
	if !ok {
		s.mtx.RUnlock()
		s.options.Tracer.UpdateStatus(spanId, 1, fmt.Sprintf("broker %s was not found", brokerId))
		return nil
	}

	s.mtx.RUnlock()

	if err := sub.Unsubscribe(); err != nil {
		s.options.Tracer.UpdateStatus(spanId, 1, err.Error())
		return err
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.subscribers, brokerId)

	s.options.Tracer.UpdateStatus(spanId, 2, "success")

	return nil
}

func (s *customSidecar) ReadFromSecretStore(ctx context.Context, secretStore string, name string) (*sidecar.Secret, error) {
	_, spanId := s.options.Tracer.Start(ctx, "customSidecar.ReadFromSecretStore")
	defer s.options.Tracer.Finish(spanId)

	s.options.Tracer.AddMetadata(spanId, map[string]string{
		"secretStore": secretStore,
		"name":        name,
	})

	sc, ok := s.options.Secrets[secretStore]
	if !ok {
		log.Warnf("secret store %s was not found", secretStore)
		s.options.Tracer.UpdateStatus(spanId, 1, fmt.Sprintf("secret store %s was not found", secretStore))
		return nil, sidecar.ErrComponentNotFound
	}

	mp, err := sc.GetSecret(name)
	if err != nil {
		s.options.Tracer.UpdateStatus(spanId, 1, fmt.Sprintf("failed to get secret: %v", err))
		return nil, err
	}

	s.options.Tracer.UpdateStatus(spanId, 2, "success")

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
