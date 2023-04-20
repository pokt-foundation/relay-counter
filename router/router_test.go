package router

import (
	"bytes"
	context "context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pokt-foundation/relay-counter/types"
	"github.com/sirupsen/logrus"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRouter_HealthCheck(t *testing.T) {
	c := require.New(t)

	router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", logrus.New())
	c.NoError(err)

	tests := []struct {
		name               string
		expectedStatusCode int
	}{
		{
			name:               "Success",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		c.NoError(err)

		rr := httptest.NewRecorder()
		router.router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_CreateSession(t *testing.T) {
	c := require.New(t)

	driverMock := &MockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, "8080", logrus.New())
	c.NoError(err)

	rawCountToSend := types.RelayCount{
		AppPublicKey: "21",
		Day:          time.Now(),
		Success:      21,
		Error:        7,
	}

	countToSend, err := json.Marshal(rawCountToSend)
	c.NoError(err)

	tests := []struct {
		name                string
		expectedStatusCode  int
		reqInput            []byte
		errReturnedByDriver error
		apiKey              string
		setMock             bool
	}{
		{
			name:               "Success",
			expectedStatusCode: http.StatusOK,
			reqInput:           countToSend,
			setMock:            true,
		},
		{
			name:               "Wrong input",
			expectedStatusCode: http.StatusBadRequest,
			reqInput:           []byte("wrong"),
		},
		{
			name:                "Failure on driver",
			expectedStatusCode:  http.StatusInternalServerError,
			reqInput:            countToSend,
			errReturnedByDriver: errors.New("dummy"),
			setMock:             true,
		},
		{
			name:               "Not authorized",
			expectedStatusCode: http.StatusUnauthorized,
			apiKey:             "wrong",
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodPost, "/v0/count", bytes.NewBuffer(tt.reqInput))
		c.NoError(err)

		req.Header.Set("Authorization", tt.apiKey)
		rr := httptest.NewRecorder()

		if tt.setMock {
			driverMock.On("WriteRelayCount", mock.Anything, mock.Anything).Return(tt.errReturnedByDriver).Once()
		}

		router.router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_RunServer(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name       string
		ctxTimeout time.Duration
	}{
		{
			name:       "Context finished",
			ctxTimeout: time.Millisecond,
		},
		{
			name:       "Context not finished",
			ctxTimeout: time.Minute,
		},
	}

	for _, tt := range tests {
		router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", logrus.New())
		c.NoError(err)

		ctxTimeout, cancel := context.WithTimeout(context.Background(), tt.ctxTimeout)
		defer cancel()

		go router.RunServer(ctxTimeout)

		time.Sleep(time.Second)
	}
}
