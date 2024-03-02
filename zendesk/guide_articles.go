package zendesk

import (
	"context"
	"time"
)

type ArticleResponse struct {
	Article Article `json:"article"`
}

type ArticlesResponse struct {
	Articles []Article `json:"articles"`
	cursorPaginationResponse
}

type Article struct {
	AuthorID          UserID            `json:"author_id"`
	Body              string            `json:"body"`
	CommentsDisabled  bool              `json:"comments_disabled"`
	ContentTagIds     []any             `json:"content_tag_ids"`
	CreatedAt         time.Time         `json:"created_at"`
	Draft             bool              `json:"draft"`
	EditedAt          time.Time         `json:"edited_at"`
	HTMLURL           string            `json:"html_url"`
	ID                ArticleID         `json:"id"`
	LabelNames        []any             `json:"label_names"`
	Locale            string            `json:"locale"`
	Name              string            `json:"name"`
	Outdated          bool              `json:"outdated"`
	PermissionGroupID PermissionGroupID `json:"permission_group_id"`
	Position          int64             `json:"position"`
	Promoted          bool              `json:"promoted"`
	SectionID         SectionID         `json:"section_id"`
	SourceLocale      string            `json:"source_locale"`
	Title             string            `json:"title"`
	UpdatedAt         time.Time         `json:"updated_at"`
	URL               string            `json:"url"`
	UserSegmentID     UserSegmentID     `json:"user_segment_id"`
	VoteCount         uint64            `json:"vote_count"`
	VoteSum           int64             `json:"vote_sum"`
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/
type ArticleService struct {
	client  *client
	generic genericService[
		ArticleID,
		ArticleResponse,
		ArticlesResponse,
	]
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/#show-article
func (s ArticleService) Show(ctx context.Context, id ArticleID) (ArticleResponse, error) {
	return s.generic.Show(ctx, id)
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/#list-articles
func (s ArticleService) List(
	ctx context.Context,
	pageHandler func(response ArticlesResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}
