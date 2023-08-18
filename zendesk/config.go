package zendesk

import (
	"log"
	"net/http"
)

type configOption func(s *internalConfig)

type RequestPreProcessor interface {
	ProcessRequest(*http.Request) error
}

type RequestPreProcessorFunc func(*http.Request) error

func (p RequestPreProcessorFunc) ProcessRequest(r *http.Request) error {
	return p(r)
}

type internalConfig struct {
	roundTripper         http.RoundTripper
	userAgent            string
	requestPreProcessors []RequestPreProcessor
}

func WithRoundTripper(roundTripper http.RoundTripper) configOption {
	return func(s *internalConfig) {
		s.roundTripper = roundTripper
	}
}

func WithUserAgent(userAgent string) configOption {
	return func(s *internalConfig) {
		s.userAgent = userAgent
	}
}

func WithRequestPreProcessor(requestPreProcessor RequestPreProcessor) configOption {
	return func(s *internalConfig) {
		s.requestPreProcessors = append(s.requestPreProcessors, requestPreProcessor)
	}
}

func WithLogger(logger *log.Logger) configOption {
	return WithRequestPreProcessor(&loggerWrapper{
		logger: logger,
	})
}

type loggerWrapper struct {
	logger *log.Logger
}

func (l *loggerWrapper) ProcessRequest(r *http.Request) error {
	l.logger.Printf("Request: %s %s", r.Method, r.URL.String())

	return nil
}
