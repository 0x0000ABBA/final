package logger

import "go.uber.org/zap"

func New(logLevel string) (*zap.SugaredLogger, error) {

	var logger *zap.Logger
	var err error

	switch logLevel {
	case "development":
		logger, err = zap.NewDevelopment()
	case "production":
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
