package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type test struct {
	tt             *testing.T
	app            *fiber.App
	request        *http.Request
	followRedirect bool
}

func NewTest(t *testing.T, s Server, method string, route string, reqBody any) test {
	reqBodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err, "failed to marshal request body")

	tst := test{
		tt:      t,
		app:     s.fiberApp,
		request: httptest.NewRequest(method, route, bytes.NewReader(reqBodyBytes)),
	}

	return tst.WithContentType("application/json")
}

func (t test) WithHeader(key string, value string) test {
	t.request.Header.Set(key, value)
	return t
}

func (t test) WithContentType(value string) test {
	return t.WithHeader("Content-Type", value)
}

func (t test) WithCookie(key string, value string) test {
	t.request.AddCookie(&http.Cookie{Name: key, Value: value})
	return t
}

func (t test) FollowRedirect() test {
	t.followRedirect = true
	return t
}

func (t test) Send() *http.Response {
	res, err := t.app.Test(t.request, -1)
	require.NoError(t.tt, err, "failed to send request")

	if t.followRedirect && res.StatusCode >= http.StatusMultipleChoices && res.StatusCode < http.StatusBadRequest {
		req := httptest.NewRequest(fiber.MethodGet, res.Header.Get("Location"), nil)
		req.Header.Set("Cookie", t.request.Header.Get("Cookie"))
		res, err = t.app.Test(req, -1)
		require.NoError(t.tt, err, "failed to follow redirect")
	}

	return res
}
