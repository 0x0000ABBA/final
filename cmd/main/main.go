package main

import (
	"context"
	"final/internal/app"
	"final/internal/config"
	"final/internal/logger"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	conf, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

	l, err := logger.New(conf.Mode)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to initialize logger: %w", err))
	}

	ctx := context.Background()

	a, err := app.New(conf, l)
	if err != nil {
		l.Fatalf("failed to init app: %v\n", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		if err := a.Run(ctx); err != nil {
			l.Fatalf("failed to run app: %v\n", err)
		}
	}()

	<-sig

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	if err := a.Shutdown(ctx); err != nil {
		l.Fatalf("failed to gracefully shutdown app: %v\n", err)
	}

	cancel()
	close(sig)
}
