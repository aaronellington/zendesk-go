package zendesk

import (
	"context"
	"time"
)

type ViewResponse struct {
	View View `json:"view"`
}

type ViewsResponse struct {
	Views []View `json:"views"`
	cursorPaginationResponse
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
	client  *client
	generic genericService[
		ViewID,
		ViewResponse,
		ViewsResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/views/#list-views
func (s *ViewService) List(
	ctx context.Context,
	pageHandler func(response ViewsResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}
