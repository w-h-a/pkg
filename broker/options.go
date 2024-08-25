package broker

import "context"

type BrokerOption func(o *BrokerOptions)

type BrokerOptions struct {
	Nodes            []string
	PublishOptions   PublishOptions
	SubscribeOptions SubscribeOptions
	Context          context.Context
}

func BrokerWithNodes(addrs ...string) BrokerOption {
	return func(o *BrokerOptions) {
		o.Nodes = addrs
	}
}

func BrokerWithPublishOptions(options PublishOptions) BrokerOption {
	return func(o *BrokerOptions) {
		o.PublishOptions = options
	}
}

func BrokerWithSubscribeOptions(options SubscribeOptions) BrokerOption {
	return func(o *BrokerOptions) {
		o.SubscribeOptions = options
	}
}

func NewBrokerOptions(publishOpts []PublishOption, subscribeOpts []SubscribeOption, opts []BrokerOption) BrokerOptions {
	options := BrokerOptions{
		PublishOptions:   NewPublishOptions(publishOpts...),
		SubscribeOptions: NewSubscribeOptions(subscribeOpts...),
		Context:          context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type PublishOption func(o *PublishOptions)

type PublishOptions struct {
	Topic   string
	Context context.Context
}

func PublishWithTopic(topic string) PublishOption {
	return func(o *PublishOptions) {
		o.Topic = topic
	}
}

func NewPublishOptions(opts ...PublishOption) PublishOptions {
	options := PublishOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type SubscribeOption func(o *SubscribeOptions)

type SubscribeOptions struct {
	Group   string
	Context context.Context
}

func SubscribeWithGroup(group string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Group = group
	}
}

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	options := SubscribeOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
