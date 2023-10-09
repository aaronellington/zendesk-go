package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type CategoriesResponse struct {
	Categories []Category `json:"categories"`
	CursorPaginationResponse
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
	client *client
}

func (s CategoryService) List(ctx context.Context, pageHandler func(response CategoriesResponse) error) error {
	query := url.Values{}
	query.Set("page[size]", "100")

	for {
		target := CategoriesResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/help_center/categories?%s", query.Encode()),
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

		query.Set("page[after]", target.Meta.AfterCursor)
	}

	return nil
}
