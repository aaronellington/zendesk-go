package study

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

type RoundTripper struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

func (r RoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return r.RoundTripFunc(request)
}

type RoundTripFunc func(t *testing.T, request *http.Request) (*http.Response, error)

func RoundTripperQueue(t *testing.T, queue []RoundTripFunc) http.RoundTripper {
	runNumber := 0

	return RoundTripper{
		RoundTripFunc: func(r *http.Request) (*http.Response, error) {
			defer func() {
				runNumber++
			}()

			if len(queue) <= runNumber {
				return nil, errors.New("empty queue")
			}

			return queue[runNumber](t, r)
		},
	}
}

type ExpectedTestRequest struct {
	Method    string
	Path      string
	Validator func(r *http.Request) error
}

type TestResponse interface {
	CreateResponse() (*http.Response, error)
}

type TestResponseFile struct {
	StatusCode int
	FilePath   string
}

func (f *TestResponseFile) CreateResponse() (*http.Response, error) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return nil, fmt.Errorf("response body file not found: %s", f.FilePath)
	}
	// defer file.Close()

	return &http.Response{
		StatusCode: f.StatusCode,
		Body:       io.NopCloser(file),
	}, nil
}

func ServeAndValidate(t *testing.T, r TestResponse, expected ExpectedTestRequest) RoundTripFunc {
	return func(t *testing.T, request *http.Request) (*http.Response, error) {
		if err := Assert(expected.Method, request.Method); err != nil {
			t.Fatal(err)
		}

		if err := Assert(expected.Path, request.URL.Path); err != nil {
			t.Fatal(err)
		}

		if expected.Validator != nil {
			if err := expected.Validator(request); err != nil {
				t.Fatal(err)
			}
		}

		return r.CreateResponse()
	}
}
