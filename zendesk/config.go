package zendesk

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type configOption func(s *internalConfig)

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
	realTimeChatWebsocketConnection *net.Conn
	requestPreProcessors            []RequestPreProcessor
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

func SetTimeout(timeout time.Duration) configOption {
	return func(s *internalConfig) {
		s.timeout = timeout
	}
}

func SetRealTimeChatWebsocketConnection(websocketConnection *net.Conn) configOption {
	return func(s *internalConfig) {
		s.realTimeChatWebsocketConnection = websocketConnection
	}
}

func SetRealTimeChatWebsocketHost(websocketHost string) configOption {
	return func(s *internalConfig) {
		s.realTimeChatWebsocketHost = websocketHost
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

func WithSlogger(logger *slog.Logger) configOption {
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
