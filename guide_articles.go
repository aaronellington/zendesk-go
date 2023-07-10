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
	AuthorID          UserID    `json:"author_id"`
	Body              string    `json:"body"`
	CommentsDisabled  bool      `json:"comments_disabled"`
	ContentTagIds     []any     `json:"content_tag_ids"`
	CreatedAt         time.Time `json:"created_at"`
	Draft             bool      `json:"draft"`
	EditedAt          time.Time `json:"edited_at"`
	HtmlURL           string    `json:"html_url"`
	ID                ArticleID `json:"id"`
	LabelNames        []any     `json:"label_names"`
	Locale            string    `json:"locale"`
	Name              string    `json:"name"`
	Outdated          bool      `json:"outdated"`
	PermissionGroupID uint64    `json:"permission_group_id"`
	Position          uint64    `json:"position"`
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
