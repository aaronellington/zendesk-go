package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type GroupsResponse struct {
	Groups []Group
	CursorPaginationResponse
}

type Group struct {
	ID        GroupID   `json:"id"`
	IsPublic  bool      `json:"is_public"`
	Name      string    `json:"name"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/
type GroupsService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/#list-groups
func (s GroupsService) List(
	ctx context.Context,
	pageHandler func(response GroupsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/groups?%s", query.Encode())

	for {
		target := GroupsResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
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

		endpoint = target.Links.Next
	}

	return nil
}
