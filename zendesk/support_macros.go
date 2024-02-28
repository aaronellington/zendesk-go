package zendesk

import (
	"context"
	"time"
)

type MacroResponse struct {
	Macro Macro `json:"macro"`
}

type MacrosResponse struct {
	Macros []Macro `json:"macros"`
	cursorPaginationResponse
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
	client  *client
	generic genericService[
		MacroID,
		MacroResponse,
		MacrosResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/macros/#list-macros
func (s *MacroService) List(
	ctx context.Context,
	pageHandler func(response MacrosResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}
