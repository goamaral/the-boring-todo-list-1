package server_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/the-boring-to-do-list-1/internal/server"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func sendRequest(t *testing.T, s server.Server, method string, route string, reqBody any, resBody any) (*http.Response, error) {
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

	// Read response body
	resBodyBuf, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	resBodyStr := string(resBodyBuf)

	// Check if error occured
	if res.StatusCode == fiber.StatusInternalServerError {
		return res, errors.New(string(resBodyBuf))
	}

	if _, ok := resBody.(*string); ok {
		*(resBody.(*string)) = resBodyStr
	} else {
		err := json.Unmarshal(resBodyBuf, resBody)
		if err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
	}

	return res, nil
}

func TestServer_HealthCheck(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := server.NewServer(nil)

		resBody := ""
		res, err := sendRequest(t, s, fiber.MethodGet, "/health", nil, &resBody)
		if assert.NoError(t, err, err) && assert.Equal(t, fiber.StatusOK, res.StatusCode, resBody) {
			assert.Equal(t, "OK", string(resBody))
		}
	})
}
