package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"example.com/fiber-m3o-validator/mocks"
	"example.com/fiber-m3o-validator/server"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	m3o "go.m3o.com"
)

func sendRequest(t *testing.T, s server.Server, method string, route string, reqBody fiber.Map, resBody any) *http.Response {
	reqBodyBytes, err := json.Marshal(reqBody)
	reqBodyBuf := bytes.NewReader(reqBodyBytes)
	if err != nil {
		t.Fatalf("failed to encode request: %v", err)
	}

	req := httptest.NewRequest(method, route, reqBodyBuf)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.Test(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	switch resBody.(type) {
	case *fiber.Map:
		err := json.NewDecoder(res.Body).Decode(resBody)
		if err != nil {
			t.Fatalf("failed to decode response body to fiber map: %v", err)
		}
	case *string:
		buf, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("failed to decode response body to string: %v", err)
		}
		*(resBody.(*string)) = string(buf)
	default:
		t.Fatalf("failed to decode response body: unsupported type %s", reflect.TypeOf(resBody).String())
	}

	return res
}

func newServer(t *testing.T, setupMocks func(*mocks.M3ODb) error) server.Server {
	m3oDbClient := mocks.NewM3ODb(t)

	err := setupMocks(m3oDbClient)
	if err != nil {
		t.Fatalf("failed to setup mocks: %v", err)
	}

	return server.NewServer(&m3o.Client{Db: m3oDbClient})
}

func TestServer_HealthCheck(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		server := newServer(t, func(mockM3ODbClient *mocks.M3ODb) error {
			return nil
		})

		resBody := ""
		res := sendRequest(t, server, fiber.MethodGet, "/health", nil, &resBody)
		if assert.Equal(t, fiber.StatusOK, res.StatusCode, resBody) {
			assert.Equal(t, "OK", string(resBody))
		}
	})
}
