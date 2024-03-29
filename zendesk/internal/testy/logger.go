package testy

import (
	"log"
	"testing"
)

func NewLogger(t *testing.T) *log.Logger {
	return log.New(testWriter{t}, t.Name()+" - ", log.LstdFlags)
}

type testWriter struct {
	t *testing.T
}

func (tw testWriter) Write(p []byte) (int, error) {
	tw.t.Log(string(p))

	return len(p), nil
}
