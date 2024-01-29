package zendesk

import (
	"context"
	"net/http"
)

type ChatDepartmentsService struct {
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

func (s *ChatDepartmentsService) List(ctx context.Context, pageHandler func(page []Department) error) error {
	requestURL := "https://www.zopim.com/api/v2/departments"

	for {
		target := []Department{}
		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			requestURL,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := s.client.LiveChatRequest(request, &target); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}
	}
}
