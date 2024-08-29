package snssqs

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/w-h-a/pkg/broker"
	"github.com/w-h-a/pkg/telemetry/log"
)

type SnsClient interface {
	ProduceToTopic(bs []byte, topic string) error
}

type snsClient struct {
	*sns.Client
}

func (c *snsClient) ProduceToTopic(bs []byte, topic string) error {
	input := &sns.PublishInput{
		Message:  aws.String(string(bs)),
		TopicArn: aws.String(topic),
	}

	if _, err := c.Publish(context.Background(), input); err != nil {
		return err
	}

	return nil
}

type SqsClient interface {
	ConsumeFromGroup(sub broker.Subscriber)
}

type sqsClient struct {
	*sqs.Client
	queueUrl          *string
	visibilityTimeout int32
	waitTimeSeconds   int32
}

type sqsMsg struct {
	Message string `json:"message"`
}

func (c *sqsClient) ConsumeFromGroup(sub broker.Subscriber) {
	result, err := c.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl:              c.queueUrl,
		MaxNumberOfMessages:   1,
		VisibilityTimeout:     c.visibilityTimeout,
		WaitTimeSeconds:       c.waitTimeSeconds,
		MessageAttributeNames: []string{"All"},
	})
	if err != nil {
		log.Errorf("failed to receive sqs message from group %s: %s", sub.Options().Group, err.Error())
		return
	}

	if len(result.Messages) == 0 {
		return
	}

	for _, msg := range result.Messages {
		body := msg.Body

		var sqsMsg sqsMsg

		if err := json.Unmarshal([]byte(*body), &sqsMsg); err != nil {
			log.Errorf("failed to unmarshal message from group %s: %s", sub.Options().Group, err)
			continue
		}

		if err := sub.Handler([]byte(sqsMsg.Message)); err != nil {
			log.Errorf("failed to handle message from group %s: %s", sub.Options().Group, err)
			continue
		}

		msgHandle := msg.ReceiptHandle
		c.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(sub.Options().Group),
			ReceiptHandle: msgHandle,
		})
	}
}
