// easygo aws helpers
package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/bdlilley/easygo/pkg/logging"
	"github.com/rotisserie/eris"
)

type EGAwsClient struct {
	cfg           aws.Config
	stsClient     *sts.Client
	secretsClient *secretsmanager.Client
}

type NewEGAwsClientArgs struct {
	Logger        logging.Logger
	Region        string
	AssumeRoleArn string
	// RetryMaxAttempts sets the maximum number of attempts (default: 3)
	// Set to 0 to use AWS default behavior
	RetryMaxAttempts int
	// RetryMode sets the retry mode (Standard, Adaptive, or Legacy)
	// If empty, defaults to Standard
	RetryMode aws.RetryMode
	// HTTPClient allows providing a custom HTTP client with custom timeout/retry logic
	// If nil, the default HTTP client will be used
	HTTPClient *http.Client
}

func NewEGAwsClient(ctx context.Context, args *NewEGAwsClientArgs) (*EGAwsClient, error) {
	// Build config options
	configOpts := []func(*config.LoadOptions) error{
		config.WithRegion(args.Region),
	}

	// Configure retry behavior
	if args.RetryMaxAttempts > 0 {
		configOpts = append(configOpts, config.WithRetryMaxAttempts(args.RetryMaxAttempts))
		args.Logger.Debug("configured retry max attempts", "maxAttempts", args.RetryMaxAttempts)
	}

	if args.RetryMode != "" {
		configOpts = append(configOpts, config.WithRetryMode(args.RetryMode))
		args.Logger.Debug("configured retry mode", "mode", args.RetryMode)
	}

	// Configure custom HTTP client if provided
	if args.HTTPClient != nil {
		configOpts = append(configOpts, config.WithHTTPClient(args.HTTPClient))
		args.Logger.Debug("using custom HTTP client")
	}

	cfg, err := config.LoadDefaultConfig(ctx, configOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	args.Logger.Debug("loaded AWS config from default credentials chain")

	stsClient := sts.NewFromConfig(cfg)

	if args.AssumeRoleArn != "" {
		args.Logger.Debug("AssumeRoleArn is set; assuming role", "roleArn", args.AssumeRoleArn)
		result, err := stsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn: aws.String(args.AssumeRoleArn),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to assume role: %w", err)
		}
		cfg.Credentials = aws.NewCredentialsCache(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     *result.Credentials.AccessKeyId,
				SecretAccessKey: *result.Credentials.SecretAccessKey,
				SessionToken:    *result.Credentials.SessionToken,
				Expires:         *result.Credentials.Expiration,
			}, nil
		}))
		args.Logger.Debug("assume role successful", "roleArn", args.AssumeRoleArn)
		stsClient = sts.NewFromConfig(cfg)
	}

	id, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to get caller identity: %w", err)
	}
	args.Logger.Debug("caller identity", "identity", id)

	secretsClient := secretsmanager.NewFromConfig(cfg)

	return &EGAwsClient{
		cfg:           cfg,
		stsClient:     stsClient,
		secretsClient: secretsClient,
	}, nil
}

// GetCallerIdentity retrieves information about the current AWS identity
func (c *EGAwsClient) GetCallerIdentity(ctx context.Context) (*sts.GetCallerIdentityOutput, error) {
	input := &sts.GetCallerIdentityInput{}
	return c.stsClient.GetCallerIdentity(ctx, input)
}

// GetConfig returns the AWS config
func (c *EGAwsClient) GetConfig() aws.Config {
	return c.cfg
}

// GetSTSClient returns the STS client
func (c *EGAwsClient) GetSTSClient() *sts.Client {
	return c.stsClient
}

// GetSecretsClient returns the SecretsManager client
func (c *EGAwsClient) GetSecretsClient() *secretsmanager.Client {
	return c.secretsClient
}

type GetLatestSecretValueResult struct {
	ByteValue []byte
}

type ByteTransformer[T any] struct {
	ByteValue []byte
}

// Gets the latest value of secretNameOrArn and unmarshals it into result
func (c *EGAwsClient) GetLatestJsonSecretValue(ctx context.Context, secretNameOrArn string, result any) error {
	output, err := c.secretsClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretNameOrArn),
	})
	if err != nil {
		return eris.Wrap(err, "failed to get secret value")
	}

	var byteValue []byte
	if output.SecretString != nil {
		byteValue = []byte(*output.SecretString)
	} else if output.SecretBinary != nil {
		byteValue = output.SecretBinary
	} else {
		return eris.New("secret found but value is empty")
	}

	err = json.Unmarshal(byteValue, result)
	if err != nil {
		return eris.Wrap(err, "failed to unmarshal byte value")
	}

	return nil
}
