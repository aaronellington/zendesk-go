package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type UserFieldResponse struct {
	UserField UserField `json:"user_field"`
}

type UserFieldsResponse struct {
	UserFields []UserField `json:"user_fields"`
	CursorPaginationResponse
}

type UserField struct {
	Active             bool                `json:"active"`
	ID                 UserFieldID         `json:"id"`
	CreatedAt          time.Time           `json:"created_at"`
	System             bool                `json:"system"`
	Description        *string             `json:"description"`
	Key                string              `json:"key"`
	Position           uint64              `json:"position"`
	RawDescription     *string             `json:"raw_description"`
	RawTitle           *string             `json:"raw_title"`
	Title              *string             `json:"title"`
	Type               string              `json:"type"`
	UpdatedAt          *time.Time          `json:"updated_at"`
	URL                string              `json:"url"`
	CustomFieldOptions []CustomFieldOption `json:"custom_field_options"`
}

type CustomFieldOption struct {
	ID       CustomFieldOptionID `json:"id"`
	Name     string              `json:"name"`
	Position uint64              `json:"position"`
	RawName  string              `json:"raw_name"`
	URL      string              `json:"url"`
	Value    string              `json:"value"`
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_fields/
type UserFieldService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_fields/#list-user-fields
func (s UserFieldService) List(
	ctx context.Context,
	pageHandler func(response UserFieldsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/user_fields?%s", query.Encode())

	for {
		target := UserFieldsResponse{}

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

func (s UserFieldService) Show(
	ctx context.Context,
	id UserFieldID,
) (UserField, error) {
	target := UserFieldResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/user_fields/%d", id),
		http.NoBody,
	)
	if err != nil {
		return UserField{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return UserField{}, err
	}

	return target.UserField, nil
}
