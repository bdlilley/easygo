package logging

import "github.com/sirupsen/logrus"

// Logger is a generic logging interface that uses logrus.FieldLogger as the base
// This provides full compatibility with Logrus while allowing for other implementations
type Logger = logrus.FieldLogger
