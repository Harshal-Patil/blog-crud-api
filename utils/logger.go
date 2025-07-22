package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger() {
	Log.SetFormatter(&logrus.JSONFormatter{}) // Structured logs for ELK/Filebeat
	Log.SetOutput(os.Stdout)                  // Default to stdout
	Log.SetLevel(logrus.InfoLevel)            // Adjust level as needed
}
