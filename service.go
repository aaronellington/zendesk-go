package zendesk

import "net/http"

func New(
	subDomain string,
	auth authentication,
	opts ...configOption,
) *Service {
	config := &internalConfig{}
	for _, opt := range opts {
		opt(config)
	}

	c := &client{
		httpClient: &http.Client{
			Transport: config.roundTripper,
		},
		subdomain:            subDomain,
		auth:                 auth,
		requestPreProcessors: config.requestPreProcessors,
	}

	return &Service{
		supportService: &SupportService{
			ticketService: TicketService{
				client: c,
			},
		},
	}
}

type Service struct {
	supportService *SupportService
}

func (s *Service) Support() *SupportService {
	return s.supportService
}
