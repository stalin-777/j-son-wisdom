package app

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	quote string
}

func (m *MockStorage) GetRandomQuote() string {
	return m.quote
}

type MockPow struct {
	challenge    string
	verifyNonce  string
	verifyResult bool
}

func (m *MockPow) GenerateChallenge() string {
	return m.challenge
}

func (m *MockPow) Verify(challenge, nonce string) bool {
	return challenge == m.challenge && nonce == m.verifyNonce && m.verifyResult
}

// Тест для метода Stop
func TestApp_Stop(t *testing.T) {
	cfg := &Cfg{
		Network:    "tcp",
		Address:    "localhost:0",
		Difficulty: 5,
	}
	mockPow := &MockPow{}
	mockStorage := &MockStorage{}
	wg := &sync.WaitGroup{}

	application, err := New(cfg, mockPow, mockStorage, wg, slog.Default())
	assert.NoError(t, err)

	go application.Run(context.Background())

	time.Sleep(time.Millisecond * 100)
	application.Stop()

	_, err = net.Dial("tcp", application.listener.Addr().String())
	assert.Error(t, err, "Listener should be closed")
}

// Табличный тест для handleConnection
func TestApp_handleConnection(t *testing.T) {
	tests := []struct {
		name              string
		challenge         string
		expectedChallenge string
		nonce             string
		verifyResult      bool
		expectedQuote     string
		expectedMsg       string
	}{
		{
			name:              "Successful PoW and returns quote",
			challenge:         "testChallenge",
			expectedChallenge: "testChallenge:mock:5",
			nonce:             "validNonce",
			verifyResult:      true,
			expectedQuote:     "Here is a random quote.",
			expectedMsg:       "Here is a random quote.\n",
		},
		{
			name:              "Failed PoW, returns error message",
			challenge:         "testChallenge",
			expectedChallenge: "testChallenge:mock:5",
			nonce:             "invalidNonce",
			verifyResult:      false,
			expectedMsg:       "PoW solution is incorrect!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPow := &MockPow{
				challenge:    tt.challenge,
				verifyNonce:  tt.nonce,
				verifyResult: tt.verifyResult,
			}
			mockStorage := &MockStorage{
				quote: tt.expectedQuote,
			}

			cfg := &Cfg{
				Network:    "tcp",
				Address:    "localhost:0",
				Algorithm:  "mock",
				Difficulty: 5,
			}
			wg := &sync.WaitGroup{}
			application, err := New(cfg, mockPow, mockStorage, wg, slog.Default())
			assert.NoError(t, err)

			clientConn, serverConn := net.Pipe()
			defer clientConn.Close()
			defer serverConn.Close()

			application.wg.Add(1)
			go application.handleConnection(context.Background(), serverConn)

			buf := make([]byte, 1024)
			n, err := clientConn.Read(buf)
			assert.NoError(t, err)

			receivedChallenge := string(buf[:n])
			assert.Equal(t, tt.expectedChallenge+"\n", receivedChallenge)

			_, err = clientConn.Write([]byte(tt.nonce))
			assert.NoError(t, err)

			n, err = clientConn.Read(buf)
			assert.NoError(t, err)

			response := string(buf[:n])
			assert.Equal(t, tt.expectedMsg, response)

			wg.Wait()
		})
	}
}
