package ssm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SsmClient interface {
	GetValue(prefix, name string) (string, error)
}

type ssmClient struct {
	*ssm.Client
}

func (c *ssmClient) GetValue(prefix, name string) (string, error) {
	input := &ssm.GetParameterInput{
		Name:           aws.String(prefix + name),
		WithDecryption: aws.Bool(true),
	}

	output, err := c.GetParameter(context.Background(), input)
	if err != nil {
		return "", err
	}

	if output.Parameter.Value == nil {
		return "", nil
	}

	return *output.Parameter.Value, nil
}
