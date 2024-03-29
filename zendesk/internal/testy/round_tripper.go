package testy

import (
	"errors"
	"net/http"
	"testing"
)

func RoundTripperQueue(t *testing.T, queue []http.RoundTripper) http.RoundTripper {
	runNumber := 0

	t.Cleanup(func() {
		if len(queue) > runNumber {
			t.Fatal("queue not empty at end of test, less requests were made than expected")
		}
	})

	return RoundTripperFunc(
		func(request *http.Request) (*http.Response, error) {
			defer func() {
				runNumber++
			}()

			if len(queue) <= runNumber {
				return nil, errors.New("empty queue before end of test, more requests made than expected")
			}

			return queue[runNumber].RoundTrip(request)
		},
	)
}

type RoundTripperFunc func(request *http.Request) (*http.Response, error)

func (r RoundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return r(request)
}
