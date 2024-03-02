package zendesk

import (
	"context"
	"time"
)

type CategoriesResponse struct {
	Categories []Category `json:"categories"`
	cursorPaginationResponse
}

type CategoryResponse struct {
	Category Category `json:"category"`
}

type Category struct {
	CreatedAt    time.Time  `json:"created_at"`
	Description  string     `json:"description"`
	HTMLURL      string     `json:"html_url"`
	ID           CategoryID `json:"id"`
	Locale       string     `json:"locale"`
	Name         string     `json:"name"`
	Outdated     bool       `json:"outdated"`
	Position     int64      `json:"position"`
	SourceLocale string     `json:"source_locale"`
	UpdatedAt    time.Time  `json:"updated_at"`
	URL          string     `json:"url"`
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/categories/
type CategoryService struct {
	client  *client
	generic genericService[
		CategoryID,
		CategoryResponse,
		CategoriesResponse,
	]
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/categories/#list-categories
func (s CategoryService) List(
	ctx context.Context,
	pageHandler func(response CategoriesResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}
