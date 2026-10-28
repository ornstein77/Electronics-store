package logger

import (
	"go.uber.org/zap"
	"os"
)

func New() (*zap.Logger, error) {
	if err := os.MkdirAll("out/logs", 0755); err != nil {
		return nil, err
	}
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout", "out/logs/app.log"}
	return c.Build()
}
