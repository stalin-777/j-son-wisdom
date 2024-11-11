package client

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"

	"github.com/stalin-777/j-son-wisdom/pkg/hash"
)

// Client represents the client part of the connection.
type Client struct {
	conn net.Conn
}

// NewClient initializes a new connection
func NewClient(network string, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

// Close closes the connection
func (c *Client) Close() error {
	return c.conn.Close()
}

func solvePoW(challenge string, hash hash.Hasher, difficulty int) string {
	nonce := 0
	for {
		guess := fmt.Sprintf("%s%d", challenge, nonce)
		hash := hash.Sum(guess)

		if strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
			return fmt.Sprintf("%d", nonce)
		}

		nonce++
	}
}

// SolveAndSend solves a PoW problem, sends the result and receives a response
func (c *Client) SolveAndSend() (string, error) {
	message, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading task: %w", err)
	}

	message = strings.TrimSpace(message)
	slog.Debug("Received message: " + message)

	parts := strings.Split(message, ":")
	if len(parts) != 3 {
		return "", fmt.Errorf("wrong message len")
	}

	challenge := parts[0]
	hasher, ok := hash.GetAlgorithm(parts[1])
	if !ok {
		return "", fmt.Errorf("wrong hash algorithm")
	}

	difficulty, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("failed to convert difficulty: %w", err)
	}

	nonce := solvePoW(challenge, hasher, difficulty)
	slog.Debug("PoW solution: " + nonce)

	_, err = fmt.Fprint(c.conn, nonce)
	if err != nil {
		return "", fmt.Errorf("error sending solution: %w", err)
	}

	reply, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return reply, nil
}
