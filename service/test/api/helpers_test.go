package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func request(
	t *testing.T,
	handler http.Handler,
	method string,
	path string,
	body io.Reader,
) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, body)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	return res
}

func jsonRequest(
	t *testing.T,
	handler http.Handler,
	method string,
	path string,
	body any,
) *httptest.ResponseRecorder {
	t.Helper()

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to encode request body: %v", err)
	}

	return request(
		t,
		handler,
		method,
		path,
		bytes.NewReader(data),
	)
}

func decodeJSON[T any](
	t *testing.T,
	res *httptest.ResponseRecorder,
) T {
	t.Helper()

	var got T

	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return got
}

func getJSON[T any](
	t *testing.T,
	handler http.Handler,
	path string,
) T {
	t.Helper()

	res := request(
		t,
		handler,
		http.MethodGet,
		path,
		nil,
	)

	requireStatus(t, res, http.StatusOK)

	return decodeJSON[T](t, res)
}

func requireStatus(
	t *testing.T,
	res *httptest.ResponseRecorder,
	want int,
) {
	t.Helper()

	if res.Code != want {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			want,
			res.Code,
		)
	}
}
