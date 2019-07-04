package main

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func BootstrapLogger(verbosity_level int) {
	formatter := &prefixed.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		QuoteEmptyFields: true,
	}
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:     "cyan",
		TimestampStyle:  "black+h",
		DebugLevelStyle: "blue+h:",
		InfoLevelStyle:  "green+h:",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red+b",
		PanicLevelStyle: "red+B",
	})
	logrus.SetFormatter(formatter)
	switch {
	case -2 >= verbosity_level:
		logrus.SetLevel(logrus.PanicLevel)
		break
	case -1 == verbosity_level:
		logrus.SetLevel(logrus.FatalLevel)
		break
	case 0 == verbosity_level:
		logrus.SetLevel(logrus.ErrorLevel)
		break
	case 1 == verbosity_level:
		logrus.SetLevel(logrus.WarnLevel)
		break
	case 2 == verbosity_level:
		logrus.SetLevel(logrus.InfoLevel)
		break
	case 3 == verbosity_level:
		logrus.SetLevel(logrus.TraceLevel)
		break
	default:
		logrus.SetLevel(logrus.DebugLevel)
		break
	}
	logrus.GetLevel()
}
