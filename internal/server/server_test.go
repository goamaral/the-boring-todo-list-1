package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/the-boring-to-do-list-1/internal/server"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildReqBodyReader(t *testing.T, reqBody any) io.Reader {
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to encode request: %v", err)
	}
	return bytes.NewReader(reqBodyBytes)
}

type testRequestResponse[T any] struct {
	t   *testing.T
	res *http.Response
}

func (trr testRequestResponse[T]) Test(statusCode int, assertBody func(resBody T)) {
	// Read response body buffer
	resBodyBuf, err := io.ReadAll(trr.res.Body)
	require.NoError(trr.t, err, "failed to read response body buffer")

	// Check status code
	require.Equal(trr.t, statusCode, trr.res.StatusCode, string(resBodyBuf))

	// Unmarshall response body
	var resBody T
	trr.unmarshalResponseBody(resBodyBuf, &resBody)

	// Run body assertions
	if assertBody != nil {
		assertBody(resBody)
	}
}

func (trr testRequestResponse[T]) unmarshalResponseBody(resBodyBuf []byte, resBody any) {
	if _, ok := resBody.(*string); ok {
		*(resBody.(*string)) = string(resBodyBuf)
	} else {
		err := json.Unmarshal(resBodyBuf, resBody)
		require.NoError(trr.t, err, "failed to decode response body")
	}
}

func testRequest[T any](t *testing.T, s server.Server, method string, route string, reqBody io.Reader) testRequestResponse[T] {
	// Send request
	req := httptest.NewRequest(method, route, reqBody)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.Test(req)
	require.NoError(t, err, "failed to send request")

	return testRequestResponse[T]{t, res}
}

func TestServer_HealthCheck(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := server.NewServer(jwt_provider.NewTestProvider(t), nil)

		testRequest[string](t, s, fiber.MethodGet, "/health", nil).
			Test(fiber.StatusOK, func(resBody string) {
				assert.Equal(t, "OK", resBody)
			})
	})
}
