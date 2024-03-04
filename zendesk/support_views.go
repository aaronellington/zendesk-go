package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ViewsResponse struct {
	Views []View `json:"Views"`
	CursorPaginationResponse
}

type View struct {
	ID          ViewID                `json:"id"`
	URL         string                `json:"url"`
	Title       string                `json:"title"`
	Active      bool                  `json:"active"`
	UpdatedAt   time.Time             `json:"updated_at"`
	CreatedAt   time.Time             `json:"created_at"`
	Default     bool                  `json:"default"`
	Position    int                   `json:"position"`
	Description *string               `json:"description"`
	Execution   ViewExecution         `json:"execution"`
	Conditions  BusinessRuleCondition `json:"conditions"`
	Restriction ViewRestriction       `json:"restriction"`
	RawTitle    string                `json:"raw_title"`
}

type ViewRestriction struct {
	ID   uint      `json:"id"`
	IDs  []GroupID `json:"ids"`
	Type string    `json:"type"`
}

type ViewGroup struct {
	ID         any    `json:"id"`
	Title      string `json:"title"`
	Filterable bool   `json:"filterable"`
	Sortable   bool   `json:"sortable"`
	Order      string `json:"order"`
}

type ViewSort struct {
	ID         any    `json:"id"`
	Title      string `json:"title"`
	Filterable bool   `json:"filterable"`
	Sortable   bool   `json:"sortable"`
	Order      string `json:"order"`
}

type ViewColumns struct {
	ID         any    `json:"id"`
	Title      string `json:"title"`
	Filterable bool   `json:"filterable"`
	Sortable   bool   `json:"sortable"`
}

type ViewFields struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Filterable bool   `json:"filterable"`
	Sortable   bool   `json:"sortable"`
}

type ViewExecution struct {
	GroupBy    string        `json:"group_by"`
	GroupOrder string        `json:"group_order"`
	SortBy     string        `json:"sort_by"`
	SortOrder  string        `json:"sort_order"`
	Group      ViewGroup     `json:"group"`
	Sort       ViewSort      `json:"sort"`
	Columns    []ViewColumns `json:"columns"`
	Fields     []ViewFields  `json:"fields"`
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/views/
type ViewService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/views/#list-views
func (s *ViewService) List(
	ctx context.Context,
	pageHandler func(response ViewsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/views?%s", query.Encode())

	for {
		target := ViewsResponse{}

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
