package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SectionsResponse struct {
	Sections     []Section `json:"sections"`
	Page         int       `json:"page"`
	PreviousPage any       `json:"previous_page"`
	NextPage     string    `json:"next_page"`
	PerPage      int       `json:"per_page"`
	PageCount    int       `json:"page_count"`
	Count        int       `json:"count"`
	SortBy       string    `json:"sort_by"`
	SortOrder    string    `json:"sort_order"`
	cursorPaginationResponse
}

type Section struct {
	CategoryID      CategoryID `json:"category_id"`
	CreatedAt       time.Time  `json:"created_at"`
	Description     string     `json:"description"`
	HTMLURL         string     `json:"html_url"`
	ID              SectionID  `json:"id"`
	Locale          string     `json:"locale"`
	Name            string     `json:"name"`
	Outdated        bool       `json:"outdated"`
	ParentSectionID any        `json:"parent_section_id"`
	Position        uint       `json:"position"`
	Sorting         string     `json:"sorting"`
	SourceLocale    string     `json:"source_locale"`
	ThemeTemplate   string     `json:"theme_template"`
	UpdatedAt       time.Time  `json:"updated_at"`
	URL             string     `json:"url"`
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/sections/
type SectionService struct {
	client *client
}

func (s SectionService) List(
	ctx context.Context,
	pageHandler func(response SectionsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")

	for {
		target := SectionsResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/help_center/sections?%s", query.Encode()),
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
