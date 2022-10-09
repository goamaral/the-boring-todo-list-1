package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/fiber-m3o-validator/server"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func sendRequest(t *testing.T, s server.Server, method string, route string, reqBody any, resBody any) *http.Response {
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

	if _, ok := resBody.(*string); ok {
		buf, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("failed to decode response body to string: %v", err)
		}
		*(resBody.(*string)) = string(buf)
	} else {
		err = json.NewDecoder(res.Body).Decode(resBody)
		if err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
	}

	return res
}

func TestServer_HealthCheck(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := server.NewServer(nil)

		resBody := ""
		res := sendRequest(t, s, fiber.MethodGet, "/health", nil, &resBody)
		if assert.Equal(t, fiber.StatusOK, res.StatusCode, resBody) {
			assert.Equal(t, "OK", string(resBody))
		}
	})
}
