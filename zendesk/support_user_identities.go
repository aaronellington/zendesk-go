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

type UserIdentityResponse struct {
	Identity UserIdentity `json:"identity"`
}

type UserIdentity struct {
	URL                string     `json:"url"`
	UserID             UserID     `json:"user_id"`
	ID                 IdentityID `json:"id"`
	Type               string     `json:"type"`
	Verified           bool       `json:"verified"`
	Primary            bool       `json:"primary"`
	UndeliverableCount uint64     `json:"undeliverable_count"`
	DeliverableState   string     `json:"deliverable_state"`
	Value              string     `json:"value"`
}

type UserIdentityPayload struct {
	Identity any `json:"identity"`
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

// https://developer.zendesk.com/api-reference/ticketing/users/user_identities/#create-identity
func (s *UserIdentityService) Create(
	ctx context.Context,
	userID UserID,
	payload UserIdentityPayload,
) (UserIdentityResponse, error) {
	target := UserIdentityResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/users/%d/identities", userID),
		structToReader(payload),
	)
	if err != nil {
		return UserIdentityResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return UserIdentityResponse{}, err
	}

	return target, nil
}
