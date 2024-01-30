package zendesk

import (
	"context"
	"net/http"
)

// https://developer.zendesk.com/api-reference/live-chat/chat-api/departments
type DepartmentService struct {
	client *client
}

type Department struct {
	Description *string  `json:"description"`
	Members     *[]int64 `json:"members"`
	Enabled     bool     `json:"enabled"`
	ID          int64    `json:"id"`
	Settings    struct {
		ChatEnabled    bool   `json:"chat_enabled"`
		SupportGroupID *int64 `json:"support_group_id"`
	} `json:"settings"`
	Name *string `json:"name"`
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/departments/#list-departments
func (s *DepartmentService) List(ctx context.Context) ([]Department, error) {
	target := []Department{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/api/v2/departments",
		http.NoBody,
	)
	if err != nil {
		return []Department{}, err
	}

	if err := s.client.ChatsRequest(request, &target); err != nil {
		return []Department{}, err
	}

	return target, nil
}
