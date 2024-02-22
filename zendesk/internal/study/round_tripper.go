package study

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	Query     url.Values
	Validator func(r *http.Request) error
}

type TestResponse interface {
	CreateResponse() (*http.Response, error)
}

type TestResponseFile struct {
	StatusCode        int
	FilePath          string
	ResponseModifiers []ResponseModifier
}

type ResponseModifier interface {
	ModifyResponse(r *http.Response)
}

type ResponseModifierFunc func(*http.Response)

func (r ResponseModifierFunc) ModifyResponse(response *http.Response) {
	r(response)
}

func WithResponseHeaders(headers map[string][]string) ResponseModifierFunc {
	return func(r *http.Response) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}

		for headerName, headerValue := range headers {
			for _, individualValue := range headerValue {
				r.Header.Add(headerName, individualValue)
			}
		}
	}
}

func (f *TestResponseFile) CreateResponse() (*http.Response, error) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return nil, fmt.Errorf("response body file not found: %s", f.FilePath)
	}
	// defer file.Close()

	response := &http.Response{
		StatusCode: f.StatusCode,
		Body:       io.NopCloser(file),
		Header:     make(http.Header),
	}

	for _, responseModifier := range f.ResponseModifiers {
		responseModifier.ModifyResponse(response)
	}

	headers := response.Header

	// Check if the Content-Type header has been set in the Header map. If not - default to application/json
	if _, ok := headers["Content-Type"]; !ok {
		response.Header.Set("Content-Type", "application/json")
	}

	return response, nil
}

type TestResponseNoContent struct {
	StatusCode int
}

func (f *TestResponseNoContent) CreateResponse() (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}

type TestResponseURLError struct {
	URLError *url.Error
}

func (f *TestResponseURLError) CreateResponse() (*http.Response, error) {
	return nil, f.URLError
}

func ServeAndValidate(t *testing.T, r TestResponse, expected ExpectedTestRequest) RoundTripFunc {
	return func(t *testing.T, request *http.Request) (*http.Response, error) {
		if err := Assert(expected.Method, request.Method); err != nil {
			t.Fatal(err)
		}

		if err := Assert(expected.Path, request.URL.Path); err != nil {
			t.Fatal(err)
		}

		if expected.Query == nil {
			expected.Query = url.Values{}
		}

		if err := Assert(expected.Query, request.URL.Query()); err != nil {
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
