package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type AutomationsResponse struct {
	Automations []Automation `json:"automations"`
	CursorPaginationResponse
}

type Automation struct {
	ID         AutomationID         `json:"id"`
	URL        string               `json:"url"`
	Title      string               `json:"title"`
	Active     bool                 `json:"active"`
	UpdatedAt  time.Time            `json:"updated_at"`
	CreatedAt  time.Time            `json:"created_at"`
	Default    bool                 `json:"default"`
	Actions    []AutomationActions  `json:"actions"`
	Conditions AutomationConditions `json:"conditions"`
	Position   int                  `json:"position"`
	RawTitle   string               `json:"raw_title"`
}

type AutomationConditions struct {
	All []AutomationCondition `json:"all"`
	Any []AutomationCondition `json:"any"`
}

type AutomationCondition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type AutomationActions struct {
	Field string `json:"field"`
	Value any    `json:"value"`
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/automations/
type AutomationService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/automations/#list-automations
func (s AutomationService) List(
	ctx context.Context,
	pageHandler func(response AutomationsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/automations?%s", query.Encode())

	for {
		target := AutomationsResponse{}

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
