package zendesk

import (
	"context"
	"fmt"
	"net/http"
)

// https://developer.zendesk.com/api-reference/live-chat/chat-api/departments
type DepartmentService struct {
	client *client
}

type Department struct {
	Description string             `json:"description"`
	Members     []UserID           `json:"members"`
	Enabled     bool               `json:"enabled"`
	ID          ChatDepartmentID   `json:"id"`
	Settings    DepartmentSettings `json:"settings"`
	Name        string             `json:"name"`
}

type DepartmentSettings struct {
	ChatEnabled    bool    `json:"chat_enabled"`
	SupportGroupID GroupID `json:"support_group_id"`
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

	if err := s.client.ChatRequest(request, &target); err != nil {
		return []Department{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/departments/#show-department
func (s *DepartmentService) Show(ctx context.Context, id ChatDepartmentID) (Department, error) {
	target := Department{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/departments/%d", id),
		http.NoBody,
	)
	if err != nil {
		return Department{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return Department{}, err
	}

	return target, nil
}
