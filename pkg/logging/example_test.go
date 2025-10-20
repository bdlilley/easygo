package logging_test

import (
	"context"
	"os"

	"github.com/bdlilley/easygo"
	"github.com/bdlilley/easygo/pkg/logging"
	"github.com/sirupsen/logrus"
)

// ExampleLogger demonstrates how to use Logrus with the Logger interface
func ExampleLogger() {
	// Create a new Logrus logger instance
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	// The Logger type is directly compatible with logrus.FieldLogger
	var logger logging.Logger = log

	// Use it with structured logging
	logger.WithField("key", "value").Info("structured log message")

	// Use it with multiple fields
	logger.WithFields(logrus.Fields{
		"component": "aws-client",
		"region":    "us-east-1",
	}).Debug("initializing AWS client")

	// Use it with error context
	err := someFunction()
	if err != nil {
		logger.WithError(err).Error("operation failed")
	}
}

// ExampleLogger_withAWSClient demonstrates using Logrus with the AWS client
func ExampleLogger_withAWSClient() {
	// Create a Logrus logger
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Create AWS client with Logrus logger
	client, err := easygo.NewAwsClient(context.Background(), &easygo.NewEGAwsClientArgs{
		Logger: log,
		Region: "us-west-2",
	})
	if err != nil {
		log.WithError(err).Fatal("failed to create AWS client")
	}

	// Use the client
	identity, err := client.GetCallerIdentity(context.Background())
	if err != nil {
		log.WithError(err).Error("failed to get caller identity")
		return
	}

	log.WithFields(logrus.Fields{
		"account": *identity.Account,
		"arn":     *identity.Arn,
	}).Info("authenticated with AWS")
}

// ExampleLogger_customImplementation shows that any logrus.FieldLogger implementation works
func ExampleLogger_customImplementation() {
	// You can use logrus.Entry as well (returned by WithField/WithFields)
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	// Create a logger with default fields that will be included in all logs
	var logger logging.Logger = log.WithFields(logrus.Fields{
		"service": "my-service",
		"version": "1.0.0",
	})

	// All subsequent logs will include the service and version fields
	logger.Info("service started")
	logger.Debug("processing request")
	logger.Error("an error occurred")
}

func someFunction() error {
	return nil
}
