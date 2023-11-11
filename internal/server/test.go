package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type test struct {
	app     *fiber.App
	request *http.Request

	statusCode int
}

func (s Server) NewTest(t *testing.T, method string, route string, reqBody any) *test {
	reqBodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err, "failed to marshal request body")

	tst := &test{
		app:        s.fiberApp,
		request:    httptest.NewRequest(method, route, bytes.NewReader(reqBodyBytes)),
		statusCode: http.StatusOK,
	}

	return tst.WithContentType("application/json")
}

func (t *test) WithHeader(key string, value string) *test {
	t.request.Header.Set(key, value)
	return t
}

func (t *test) WithContentType(value string) *test {
	return t.WithHeader("Content-Type", value)
}

func (t *test) WithAuthorizationHeader(value string) *test {
	return t.WithHeader("Authorization", fmt.Sprintf("Bearer %s", value))
}

func (t *test) WithStatusCode(code int) *test {
	t.statusCode = code
	return t
}

func (t *test) Send(resBody any) error {
	res, err := t.app.Test(t.request, -1)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if t.statusCode == res.StatusCode {
		if resBody != nil {
			err := json.Unmarshal(bodyBytes, resBody)
			if err != nil {
				return fmt.Errorf("failed to unmarshal response body: %w", err)
			}
		}
	} else {
		return fmt.Errorf("unexpected status code (Got: %d, Expected: %d) with body (%s)", res.StatusCode, t.statusCode, string(bodyBytes))
	}

	return nil
}
