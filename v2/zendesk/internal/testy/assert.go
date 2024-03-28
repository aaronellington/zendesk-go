package testy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-test/deep"
)

func Assert[T any](expected T, actual T) error {
	if diff := deep.Equal(expected, actual); len(diff) > 0 {
		diff = append([]string{""}, diff...)

		return fmt.Errorf("assert failed: %s", strings.Join(diff, "\n- "))
	}

	return nil
}

type RequestResponseTester struct {
	Request       *http.Request
	Response      *http.Response
	ResponseError error
}

func (i RequestResponseTester) RoundTrip(request *http.Request) (*http.Response, error) {
	if err := Assert(request, i.Request); err != nil {
		return nil, err
	}

	return i.Response, i.ResponseError
}
