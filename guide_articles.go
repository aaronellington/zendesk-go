package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ArticleID uint64

type ArticleResponse struct {
	Article Article `json:"article"`
}

type Article struct {
	ID                ArticleID `json:"id"`
	URL               string    `json:"url"`
	HtmlURL           string    `json:"html_url"`
	AuthorID          uint64    `json:"author_id"`
	CommentsDisabled  bool      `json:"comments_disabled"`
	Draft             bool      `json:"draft"`
	Promoted          bool      `json:"promoted"`
	Position          uint64    `json:"position"`
	VoteSum           uint64    `json:"vote_sum"`
	VoteCount         uint64    `json:"vote_count"`
	SectionID         uint64    `json:"section_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Name              string    `json:"name"`
	Title             string    `json:"title"`
	SourceLocale      string    `json:"source_locale"`
	Locale            string    `json:"locale"`
	Outdated          bool      `json:"outdated"`
	EditedAt          time.Time `json:"edited_at"`
	UserSegmentID     uint64    `json:"user_segment_id"`
	PermissionGroupID uint64    `json:"permission_group_id"`
	ContentTagIds     []any     `json:"content_tag_ids"`
	LabelNames        []any     `json:"label_names"`
	Body              string    `json:"body"`
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/
type ArticlesService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/#show-article
func (s ArticlesService) Show(ctx context.Context, id ArticleID) (Article, error) {
	target := ArticleResponse{}

	if err := s.client.jsonRequest(
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
