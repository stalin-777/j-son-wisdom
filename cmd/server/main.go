package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/stalin-777/j-son-wisdom/internal/storage"
	"github.com/stalin-777/j-son-wisdom/pkg/hash"
	"github.com/stalin-777/j-son-wisdom/pkg/pow"

	"github.com/stalin-777/j-son-wisdom/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := app.NewCfg()
	if err != nil {
		log.Fatal("Failed to initialize storage", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)
	if !cfg.IsProduction {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		logger.Info("Log level: " + slog.LevelDebug.String())
	}

	storage, err := storage.New()
	if err != nil {
		log.Fatal("Failed to initialize storage", err)
	}

	hasher, ok := hash.GetAlgorithm(cfg.Algorithm)
	if !ok {
		log.Fatal("Wrong hash algorithm", err)
	}

	pow := pow.New(cfg.Difficulty, hasher, nil)

	wg := &sync.WaitGroup{}
	a, err := app.New(cfg, pow, storage, wg, logger)
	if err != nil {
		log.Fatal("Failed to start application: ", err)
	}

	a.Run(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	cancel()
	wg.Wait()

	a.Stop()

	logger.Info("The application was gracefully shutdown")
}
