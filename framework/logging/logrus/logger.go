package logrus

import (
	"series/adapter/logger"

	"github.com/sirupsen/logrus"
)

func New(level string) (logger.Logger, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(lvl)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
	return logger, nil
}
