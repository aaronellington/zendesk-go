package zendesk

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type ConfigOption func(s *internalConfig)

type RequestPreProcessor interface {
	ProcessRequest(r *http.Request) error
}

type RequestPreProcessorFunc func(*http.Request) error

func (p RequestPreProcessorFunc) ProcessRequest(r *http.Request) error {
	return p(r)
}

type internalConfig struct {
	roundTripper                    http.RoundTripper
	userAgent                       string
	timeout                         time.Duration
	realTimeChatWebsocketHost       string
	requestPreProcessors            []RequestPreProcessor
}

func WithRoundTripper(roundTripper http.RoundTripper) ConfigOption {
	return func(s *internalConfig) {
		s.roundTripper = roundTripper
	}
}

func WithUserAgent(userAgent string) ConfigOption {
	return func(s *internalConfig) {
		s.userAgent = userAgent
	}
}

func SetTimeout(timeout time.Duration) ConfigOption {
	return func(s *internalConfig) {
		s.timeout = timeout
	}
}

func SetRealTimeChatWebsocketHost(websocketHost string) ConfigOption {
	return func(s *internalConfig) {
		websocketHost = strings.Replace(websocketHost, "http", "ws", 1)
		s.realTimeChatWebsocketHost = websocketHost
	}
}

func WithRequestPreProcessor(requestPreProcessor RequestPreProcessor) ConfigOption {
	return func(s *internalConfig) {
		s.requestPreProcessors = append(s.requestPreProcessors, requestPreProcessor)
	}
}

func WithLogger(logger *log.Logger) ConfigOption {
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

func WithSlogger(logger *slog.Logger) ConfigOption {
	return WithRequestPreProcessor(&sloggerWrapper{
		logger: logger,
	})
}

type sloggerWrapper struct {
	logger *slog.Logger
}

func (s *sloggerWrapper) ProcessRequest(r *http.Request) error {
	s.logger.Debug(fmt.Sprintf("Request: %s %s", r.Method, r.URL.String()))

	return nil
}
