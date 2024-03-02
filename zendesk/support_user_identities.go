package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type UserIdentitiesResponse struct {
	Identities []UserIdentity `json:"identities"`
	cursorPaginationResponse
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

	return genericList(
		ctx,
		s.client,
		endpoint,
		pageHandler,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_identities/#create-identity
func (s *UserIdentityService) Create(
	ctx context.Context,
	userID UserID,
	payload UserIdentityPayload,
) (UserIdentityResponse, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/users/%d/identities", userID),
		structToReader(payload),
	)
	if err != nil {
		return UserIdentityResponse{}, err
	}

	return genericRequest[UserIdentityResponse](s.client, request)
}
