package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
)

// storage interface defines methods for retrieving random quotes
type storage interface {
	GetRandomQuote() string
}

type pow interface {
	GenerateChallenge() string
	Verify(challenge, nonce string) bool
}

// App represents the main application structure
type App struct {
	listener   net.Listener
	pow        pow
	hash       string
	difficulty int
	storage    storage
	wg         *sync.WaitGroup
	log        *slog.Logger
}

// New creates a new instance of the App
func New(cfg *Cfg, p pow, s storage, wg *sync.WaitGroup, log *slog.Logger) (*App, error) {
	listener, err := net.Listen(cfg.Network, cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("Failed to create server: %w", err)
	}

	log.Info(fmt.Sprintf("Server started on address: %s", cfg.Address))

	return &App{
		listener:   listener,
		pow:        p,
		difficulty: cfg.Difficulty,
		hash:       cfg.Algorithm,
		storage:    s,
		wg:         &sync.WaitGroup{},
		log:        log,
	}, nil
}

// Stop gracefully shuts down the application
func (a *App) Stop() {
	err := a.listener.Close()
	if err != nil {
		a.log.Error("failed to gracefully shutdown application", slog.Any("error", err))
	}
}

// Run starts the main loop for accepting connections
func (a *App) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			conn, err := a.listener.Accept()
			if err != nil {
				a.log.Error("failed to accept connection", slog.Any("error", err))
				continue
			}

			a.wg.Add(1)
			go a.handleConnection(ctx, conn)
		}
	}()
}

func (a *App) handleConnection(_ context.Context, conn net.Conn) {
	defer conn.Close()
	defer a.wg.Done()

	// 1. Sending a PoW task and difficulty to the client
	challenge := a.pow.GenerateChallenge()
	a.log.Debug("Generated challenge: " + challenge)

	task := fmt.Sprintf("%s:%s:%d\n", challenge, a.hash, a.difficulty)
	_, err := conn.Write([]byte(task))
	if err != nil {
		a.log.Error("failed to write challenge", slog.Any("error", err))
		return
	}

	// 2. Receiving a response with nonce from the client
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		a.log.Error("failed to read nonce", slog.Any("error", err))
		return
	}

	nonce := string(buffer[:n])
	a.log.Debug("Got nonce: " + nonce)

	// 3. PoW Verification
	if a.pow.Verify(challenge, nonce) {
		// 4. If PoW is correct, send a random quote
		quote := a.storage.GetRandomQuote()

		_, err := conn.Write([]byte(quote + "\n"))
		if err != nil {
			a.log.Error("failed to write quote:", slog.Any("error", err))
		}
	} else {
		_, err := conn.Write([]byte("PoW solution is incorrect!\n"))
		if err != nil {
			a.log.Error("failed to write error message", slog.Any("error", err))
		}
	}
}
