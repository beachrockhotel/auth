package main

import (
	"context"
	"github.com/beachrockhotel/auth/internal/logger"
	"go.uber.org/zap/zapcore"
	"log"

	"github.com/beachrockhotel/auth/internal/app"
)

func main() {
	logger.InitDefault(zapcore.InfoLevel)
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
