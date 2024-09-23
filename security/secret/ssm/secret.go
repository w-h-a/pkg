package ssm

import (
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/w-h-a/pkg/security/secret"
	"github.com/w-h-a/pkg/telemetry/log"
)

type ssmSecret struct {
	options   secret.SecretOptions
	ssmClient SsmClient
}

func (s *ssmSecret) Options() secret.SecretOptions {
	return s.options
}

// TODO: retry
func (s *ssmSecret) GetSecret(name string) (map[string]string, error) {
	value, err := s.ssmClient.GetValue(name)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		name: value,
	}, nil
}

func (s *ssmSecret) String() string {
	return "ssm"
}

func (s *ssmSecret) configure() error {
	if len(s.options.Nodes) == 0 {
		return fmt.Errorf("secret store addresses are required")
	}

	if ssm, ok := GetSsmClientFromContext(s.options.Context); ok {
		s.ssmClient = ssm
	}

	if s.ssmClient != nil {
		return nil
	}

	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion("us-west-2"),
	)
	if err != nil {
		return err
	}

	s.ssmClient = &ssmClient{ssm.NewFromConfig(
		cfg,
		func(o *ssm.Options) {
			o.EndpointResolverV2 = &ssmResolver{s.options.Nodes}
		},
	)}

	return nil
}

func NewSecret(opts ...secret.SecretOption) secret.Secret {
	options := secret.NewSecretOptions(opts...)

	s := &ssmSecret{
		options: options,
	}

	if err := s.configure(); err != nil {
		log.Fatal(err)
	}

	return s
}
