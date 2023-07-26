package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ArticleID uint64

type ArticleResponse struct {
	Article Article `json:"article"`
}

type ArticlesResponse struct {
	Articles []Article `json:"articles"`
	CursorPaginationResponse
}

type Article struct {
	AuthorID          UserID    `json:"author_id"`
	Body              string    `json:"body"`
	CommentsDisabled  bool      `json:"comments_disabled"`
	ContentTagIds     []any     `json:"content_tag_ids"`
	CreatedAt         time.Time `json:"created_at"`
	Draft             bool      `json:"draft"`
	EditedAt          time.Time `json:"edited_at"`
	HTMLURL           string    `json:"html_url"`
	ID                ArticleID `json:"id"`
	LabelNames        []any     `json:"label_names"`
	Locale            string    `json:"locale"`
	Name              string    `json:"name"`
	Outdated          bool      `json:"outdated"`
	PermissionGroupID uint64    `json:"permission_group_id"`
	Position          int64     `json:"position"`
	Promoted          bool      `json:"promoted"`
	SectionID         uint64    `json:"section_id"`
	SourceLocale      string    `json:"source_locale"`
	Title             string    `json:"title"`
	UpdatedAt         time.Time `json:"updated_at"`
	URL               string    `json:"url"`
	UserSegmentID     uint64    `json:"user_segment_id"`
	VoteCount         uint64    `json:"vote_count"`
	VoteSum           uint64    `json:"vote_sum"`
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/
type ArticlesService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/#show-article
func (s ArticlesService) Show(ctx context.Context, id ArticleID) (Article, error) {
	target := ArticleResponse{}

	if err := s.client.ZendeskRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/help_center/articles/%d", id),
		http.NoBody,
		&target,
	); err != nil {
		return Article{}, err
	}

	return target.Article, nil
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/#list-articles
func (s ArticlesService) List(
	ctx context.Context,
	pageHandler func(response ArticlesResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")

	for {
		target := ArticlesResponse{}

		if err := s.client.ZendeskRequest(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/help_center/articles?%s", query.Encode()),
			http.NoBody,
			&target,
		); err != nil {
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
