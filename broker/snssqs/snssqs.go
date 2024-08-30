package snssqs

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"github.com/w-h-a/pkg/broker"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/utils/datautils"
)

const (
	defaultVisibilityTimeout int32 = 8
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
	bs, err := datautils.Stringify(data)
	if err != nil {
		return err
	}

	if err := b.snsClient.ProduceToTopic(bs, options.Topic); err != nil {
		return err
	}

	return nil
}

func (b *snssqs) Subscribe(callback func([]byte) error, options broker.SubscribeOptions) broker.Subscriber {
	sub := &subscriber{
		options: options,
		id:      uuid.New().String(),
		handler: callback,
		exit:    make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-sub.exit:
				return
			default:
				b.sqsClient.ConsumeFromGroup(sub)
				time.Sleep(time.Second)
			}
		}
	}()

	return sub
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

	if b.snsClient != nil || b.sqsClient != nil {
		return nil
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion("us-west-2"))
	if err != nil {
		return err
	}

	if b.options.PublishOptions != nil {
		b.snsClient = &snsClient{sns.NewFromConfig(cfg)}
	}

	if b.options.SubscribeOptions != nil {
		visibilityTimeout := defaultVisibilityTimeout

		waitTimeSeconds := defaultWaitSeconds

		if timeout, ok := GetVisibilityTimeoutFromContext(b.options.SubscribeOptions.Context); ok {
			visibilityTimeout = timeout
		}

		if waitTime, ok := GetWaitTimeSecondsFromContext(b.options.SubscribeOptions.Context); ok {
			waitTimeSeconds = waitTime
		}

		client := sqs.NewFromConfig(cfg)

		url, err := client.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
			QueueName: aws.String(b.options.SubscribeOptions.Group),
		})
		if err != nil {
			return err
		}

		b.sqsClient = &sqsClient{client, url.QueueUrl, visibilityTimeout, waitTimeSeconds}
	}

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
