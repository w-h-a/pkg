package snssqs

import (
	"context"
	"encoding/json"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/w-h-a/pkg/broker"
	"github.com/w-h-a/pkg/telemetry/log"
)

const (
	defaultVisibilityTimeout int32 = 4
	defaultWaitSeconds       int32 = 8
)

type snssqs struct {
	options   broker.BrokerOptions
	snsClient SnsClient
	sqsClient SqsClient
}

func (b *snssqs) Options() broker.BrokerOptions {
	return b.options
}

func (b *snssqs) Publish(data interface{}, options broker.PublishOptions) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := b.snsClient.ProduceToTopic(bs, b.options.Topic); err != nil {
		return err
	}

	return nil
}

func (b *snssqs) Subscribe(callback func([]byte) error, options broker.SubscribeOptions) {
	for {
		b.sqsClient.ConsumeFromGroup(callback, b.options.Group, options)
		time.Sleep(time.Second)
	}
}

func (b *snssqs) String() string {
	return "snssqs"
}

func (b *snssqs) configure() error {
	if sns, ok := GetSnsClientFromContext(b.options.Context); ok {
		b.snsClient = sns
	}

	if sqs, ok := GetSqsClientFromContext(b.options.Context); ok {
		b.sqsClient = sqs
	}

	if b.snsClient != nil && b.sqsClient != nil {
		return nil
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion("us-west-2"))
	if err != nil {
		return err
	}

	b.snsClient = &snsClient{sns.NewFromConfig(cfg)}

	visibilityTimeout := defaultVisibilityTimeout

	waitTimeSeconds := defaultWaitSeconds

	if timeout, ok := GetVisibilityTimeoutFromContext(b.options.SubscribeOptions.Context); ok {
		visibilityTimeout = timeout
	}

	if waitTime, ok := GetWaitTimeSecondsFromContext(b.options.SubscribeOptions.Context); ok {
		waitTimeSeconds = waitTime
	}

	b.sqsClient = &sqsClient{sqs.NewFromConfig(cfg), visibilityTimeout, waitTimeSeconds}

	return nil
}

func NewBroker(opts ...broker.BrokerOption) broker.Broker {
	options := broker.NewBrokerOptions(opts...)

	b := &snssqs{
		options: options,
	}

	if err := b.configure(); err != nil {
		log.Fatal(err)
	}

	return b
}
