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

type test[RES any] struct {
	tt      *testing.T
	app     *fiber.App
	request *http.Request
}

func NewTest[RES any](t *testing.T, s Server, method string, route string, reqBody any) test[RES] {
	reqBodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err, "failed to marshal request body")

	tst := test[RES]{
		tt:      t,
		app:     s.fiberApp,
		request: httptest.NewRequest(method, route, bytes.NewReader(reqBodyBytes)),
	}

	return tst.WithContentType("application/json")
}

func (t test[RES]) WithHeader(key string, value string) test[RES] {
	t.request.Header.Set(key, value)
	return t
}

func (t test[RES]) WithContentType(value string) test[RES] {
	return t.WithHeader("Content-Type", value)
}

func (t test[RES]) WithAuthorizationHeader(value string) test[RES] {
	return t.WithHeader("Authorization", fmt.Sprintf("Bearer %s", value))
}

func (t test[RES]) Send() *testResult[RES] {
	res, err := t.app.Test(t.request, -1)
	require.NoError(t.tt, err, "failed to send request")

	return &testResult[RES]{tt: t.tt, res: res}
}

type testResult[RES any] struct {
	tt                *testing.T
	res               *http.Response
	statusCodeChecked bool
}

func (t *testResult[RES]) ExpectsStatusCode(statusCode int) *testResult[RES] {
	require.Equal(t.tt, statusCode, t.res.StatusCode)
	t.statusCodeChecked = true
	return t
}

func (t testResult[RES]) Body() []byte {
	if !t.statusCodeChecked {
		t.ExpectsStatusCode(fiber.StatusOK)
	}

	bodyBytes, err := io.ReadAll(t.res.Body)
	require.NoError(t.tt, err, "failed to read response body")
	return bodyBytes
}

func (t testResult[RES]) UnmarshalBody() RES {
	var dst RES
	require.NoError(t.tt, json.Unmarshal(t.Body(), &dst), "failed to unmarshal body")
	return dst
}
