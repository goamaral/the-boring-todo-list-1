package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/fiber-m3o-validator/server"
	"github.com/gofiber/fiber/v2"
)

func sendRequest(t *testing.T, s server.Server, method string, route string, reqBody fiber.Map) (fiber.Map, *http.Response) {
	reqBodyBytes, err := json.Marshal(reqBody)
	reqBodyBuf := bytes.NewReader(reqBodyBytes)
	if err != nil {
		t.Fatalf("failed to encode request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/tasks", reqBodyBuf)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.Test(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	resBody := fiber.Map{}
	json.NewDecoder(res.Body).Decode(&resBody)

	return resBody, res
}
