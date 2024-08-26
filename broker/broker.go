package broker

type Broker interface {
	Options() BrokerOptions
	Publish(data interface{}, options PublishOptions) error
	Subscribe(callback func([]byte) error, options SubscribeOptions) Subscriber
	String() string
}
