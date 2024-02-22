package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TriggersResponse struct {
	Triggers []Trigger `json:"triggers"`
	CursorPaginationResponse
}

type Trigger struct {
	ID          TriggerID              `json:"id"`
	URL         string                 `json:"url"`
	Title       string                 `json:"title"`
	Active      bool                   `json:"active"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CreatedAt   time.Time              `json:"created_at"`
	Default     bool                   `json:"default"`
	Actions     []BusinessRuleAction   `json:"actions"`
	Conditions  BusinessRuleConditions `json:"conditions"`
	Description *string                `json:"description"`
	Position    int                    `json:"position"`
	RawTitle    string                 `json:"raw_title"`
	CategoryID  string                 `json:"category_id"`
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/triggers/
type TriggerService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/triggers/#list-triggers
func (s *TriggerService) List(
	ctx context.Context,
	pageHandler func(response TriggersResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/triggers?%s", query.Encode())

	for {
		target := TriggersResponse{}

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
