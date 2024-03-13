package common

import (
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01-02 15:04:05",
	})
	Logger.SetLevel(logrus.DebugLevel)
}
