package broker

import "context"

type BrokerOption func(o *BrokerOptions)

type BrokerOptions struct {
	Nodes            []string
	Topic            string
	Group            string
	PublishOptions   PublishOptions
	SubscribeOptions SubscribeOptions
	Context          context.Context
}

func BrokerWithNodes(addrs ...string) BrokerOption {
	return func(o *BrokerOptions) {
		o.Nodes = addrs
	}
}

func BrokerWithTopic(topic string) BrokerOption {
	return func(o *BrokerOptions) {
		o.Topic = topic
	}
}

func BrokerWithGroup(group string) BrokerOption {
	return func(o *BrokerOptions) {
		o.Group = group
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

func NewBrokerOptions(opts ...BrokerOption) BrokerOptions {
	options := BrokerOptions{
		PublishOptions: PublishOptions{
			Context: context.Background(),
		},
		SubscribeOptions: SubscribeOptions{
			Context: context.Background(),
		},
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type PublishOption func(o *PublishOptions)

type PublishOptions struct {
	Context context.Context
}

type SubscribeOption func(o *SubscribeOptions)

type SubscribeOptions struct {
	Context context.Context
}
