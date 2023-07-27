package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type UserIdentitiesResponse struct {
	Identities []UserIdentity `json:"identities"`
	CursorPaginationResponse
}

type UserIdentity struct {
	ID    uint64 `json:"id"`
	Value string `json:"value"`
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_identities/
type UserIdentityService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_identities/#list-identities
func (s *UserIdentityService) List(
	ctx context.Context,
	userID UserID,
	pageHandler func(response UserIdentitiesResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/users/%d/identities?%s", userID, query.Encode())

	for {
		target := UserIdentitiesResponse{}

		if err := s.client.ZendeskRequest(
			ctx,
			http.MethodGet,
			endpoint,
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

		endpoint = target.Links.Next
	}

	return nil
}
