package zendesk_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/aaronellington/zendesk-go/v2/zendesk"
	"github.com/aaronellington/zendesk-go/v2/zendesk/internal/testy"
)

func getTestInstance(
	t *testing.T,
	queue []http.RoundTripper,
) *zendesk.Service {
	return zendesk.New(
		"example",
		zendesk.AuthenticationToken{
			Email: "user@example.com",
			Token: "not-a-real-token",
		},
		zendesk.WithLogger(testy.NewLogger(t)),
		zendesk.WithRoundTripper(testy.RoundTripperQueue(t, queue)),
	)
}

func createBaseRequest(t *testing.T, ctx context.Context, method string, url string, payloadFilePath string) *http.Request {
	var payload io.Reader = http.NoBody

	if payloadFilePath != "" {
		base := filepath.Base(payloadFilePath)
		dir := filepath.Dir(payloadFilePath)

		abs, err := filepath.Abs(dir)
		if err != nil {
			t.Fatal(err)
		}

		file, err := os.Open(payloadFilePath)
		if err != nil {
			t.Fatalf("test payload file not found: %s/%s", abs, base)
		}
		defer file.Close()

		payloadBytes, err := io.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}

		payloadBytes = bytes.TrimSuffix(payloadBytes, []byte("\n"))

		payload = bytes.NewReader(payloadBytes)
	}

	targetRequest, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		t.Fatal(err)
	}

	targetRequest.Host = "example.zendesk.com"
	targetRequest.Header.Set("User-Agent", "aaronellington/zendesk-go")
	targetRequest.Header.Set("Accept", "application/json")
	targetRequest.Header.Set("Content-Type", "application/json")
	targetRequest.SetBasicAuth("user@example.com/token", "not-a-real-token")

	return targetRequest
}

func createResponse(t *testing.T, statusCode int, filePath string) *http.Response {
	base := filepath.Base(filePath)
	dir := filepath.Dir(filePath)

	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("test response file not found: %s/%s", abs, base)
	}

	response := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(file),
		Header:     make(http.Header),
	}

	return response
}
