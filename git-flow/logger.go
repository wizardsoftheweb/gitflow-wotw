package main

import (
	"github.com/sirupsen/logrus"
)

func BootstrapLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		QuoteEmptyFields: true,
	})
	logrus.SetLevel(logrus.TraceLevel)
}
