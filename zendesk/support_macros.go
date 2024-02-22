package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type MacrosResponse struct {
	Macros []Macro `json:"macros"`
	CursorPaginationResponse
}

type Macro struct {
	URL         string               `json:"url"`
	ID          int64                `json:"id"`
	Title       string               `json:"title"`
	Active      bool                 `json:"active"`
	UpdatedAt   time.Time            `json:"updated_at"`
	CreatedAt   time.Time            `json:"created_at"`
	Default     bool                 `json:"default"`
	Position    int                  `json:"position"`
	Description any                  `json:"description"`
	Actions     []BusinessRuleAction `json:"actions"`
	Restriction any                  `json:"restriction"`
	RawTitle    string               `json:"raw_title"`
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/macros/
type MacroService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/macros/#list-macros
func (s *MacroService) List(
	ctx context.Context,
	pageHandler func(response MacrosResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/macros?%s", query.Encode())

	for {
		target := MacrosResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := s.client.ZendeskRequest(request, &target); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}

		if !target.Meta.HasMore {
			break
		}

		endpoint = target.Links.Next
	}

	return nil
}
