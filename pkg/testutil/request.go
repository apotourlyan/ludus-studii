package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/httputil/content"
	"github.com/apotourlyan/ludus-studii/pkg/httputil/header"
)

func SuccessRequest[TResp any](
	t *testing.T,
	method, path, correlationID string,
	body any,
) (*TResp, int) {
	resp, _, code := Request[TResp, any](t, method, path, correlationID, body)
	return resp, code
}

func ErrorRequest[TErr any](
	t *testing.T,
	method, path, correlationID string,
	body any,
) (*TErr, int) {
	_, err, code := Request[any, TErr](t, method, path, correlationID, body)
	return err, code
}

func Request[TResp any, TErr any](
	t *testing.T,
	method, path, correlationID string,
	body any,
) (*TResp, *TErr, int) {
	t.Helper()

	// Marshal request body
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create request
	req, err := http.NewRequest(
		method,
		path,
		bodyReader,
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set(header.ContentType, content.ApplicationJson)
	req.Header.Set(header.CorrelationID, correlationID)

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	// Decode response if 2xx
	if statusCode >= 200 && statusCode < 300 {
		var result TResp
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		return &result, nil, statusCode
	}

	// Decode response if 4xx/5xx
	if statusCode >= 400 && statusCode < 600 {
		var errResult TErr
		if err := json.NewDecoder(resp.Body).Decode(&errResult); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		return nil, &errResult, statusCode
	}

	// Non-2xx/4xx/5xx response - return nil
	return nil, nil, statusCode
}
