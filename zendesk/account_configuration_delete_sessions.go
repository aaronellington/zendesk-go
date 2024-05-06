package zendesk

import (
	"context"
	"fmt"
	"net/http"
)

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/sessions
type SessionService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/sessions/#bulk-delete-sessions
func (s SessionService) BulkDelete(
	ctx context.Context,
	id UserID,
) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/users/%d/sessions", id),
		nil,
	)
	if err != nil {
		return err
	}

	if err := s.client.ZendeskRequest(request, nil); err != nil {
		return err
	}

	return nil
}
