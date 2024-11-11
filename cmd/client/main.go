package main

import (
	"log/slog"

	"github.com/stalin-777/j-son-wisdom/pkg/client"

	"github.com/spf13/pflag"
)

func main() {
	protocol := pflag.StringP("network", "n", "tcp", "Connection network(protocol). Default: tcp")
	address := pflag.StringP("address", "a", "localhost:8080", "Connection address. Default: localhost:8080")

	pflag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)
	
	c, err := client.NewClient(*protocol, *address)
	if err != nil {
		slog.Error("Connection error: ", slog.Any("error", err))
		return
	}
	defer c.Close()

	res, err := c.SolveAndSend()
	if err != nil {
		slog.Error("", slog.Any("error", err))
	}

	slog.Info("Response: " + res)
}
