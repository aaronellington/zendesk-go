package zendesk

import (
	"context"
	"fmt"
	"net/http"
)

// https://developer.zendesk.com/api-reference/live-chat/chat-api/oauth_clients/
type OAuthClientService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/oauth_clients/#json-format
type OAuthClientConfiguration struct {
	AgentID          UserID                `json:"agent_id"`
	ClientIdentifier string                `json:"client_identifier"`
	ClientSecret     string                `json:"client_secret"`
	ClientType       string                `json:"client_type"`
	Company          string                `json:"company"`
	CreateDate       string                `json:"create_date"`
	ID               LiveChatOAuthClientID `json:"id"`
	Name             string                `json:"name"`
	RedirectURIs     string                `json:"redirect_uris"`
	Scopes           string                `json:"scopes"`
	UpdateDate       *string               `json:"update_date"`
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/oauth_clients/#list-oauth-clients
func (s OAuthClientService) List(
	ctx context.Context,
) ([]OAuthClientConfiguration, error) {
	target := []OAuthClientConfiguration{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/api/v2/oauth/clients",
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	if err := s.client.ChatRequest(request, &target); err != nil {
		return nil, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/oauth_clients/#show-oauth-client
func (s OAuthClientService) Show(
	ctx context.Context,
	id LiveChatOAuthClientID,
) (OAuthClientConfiguration, error) {
	target := OAuthClientConfiguration{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/oauth/clients/%d", id),
		http.NoBody,
	)
	if err != nil {
		return OAuthClientConfiguration{}, err
	}

	if err := s.client.ChatRequest(request, &target); err != nil {
		return OAuthClientConfiguration{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/oauth_clients/#create-oauth-client
func (s OAuthClientService) Create(
	ctx context.Context,
	payload any,
) (OAuthClientConfiguration, error) {
	target := OAuthClientConfiguration{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/api/v2/oauth/clients",
		structToReader(payload),
	)
	if err != nil {
		return OAuthClientConfiguration{}, err
	}

	if err := s.client.ChatRequest(request, &target); err != nil {
		return OAuthClientConfiguration{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/oauth_clients/#update-oauth-client
func (s OAuthClientService) Update(
	ctx context.Context,
	id LiveChatOAuthClientID,
	payload any,
) (OAuthClientConfiguration, error) {
	target := OAuthClientConfiguration{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/oauth/clients/%d", id),
		structToReader(payload),
	)
	if err != nil {
		return OAuthClientConfiguration{}, err
	}

	if err := s.client.ChatRequest(request, &target); err != nil {
		return OAuthClientConfiguration{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/oauth_clients/#delete-oauth-client
func (s OAuthClientService) Delete(
	ctx context.Context,
	id LiveChatOAuthClientID,
) (OAuthClientConfiguration, error) {
	target := OAuthClientConfiguration{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/oauth/clients/%d", id),
		http.NoBody,
	)
	if err != nil {
		return OAuthClientConfiguration{}, err
	}

	if err := s.client.ChatRequest(request, &target); err != nil {
		return OAuthClientConfiguration{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/oauth_clients/#generate-oauth-client-secret
func (s OAuthClientService) GenerateOAuthClientSecret(
	ctx context.Context,
	id LiveChatOAuthClientID,
) (OAuthClientConfiguration, error) {
	target := OAuthClientConfiguration{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/oauth/clients/%d/client_secret", id),
		http.NoBody,
	)
	if err != nil {
		return OAuthClientConfiguration{}, err
	}

	if err := s.client.ChatRequest(request, &target); err != nil {
		return OAuthClientConfiguration{}, err
	}

	return target, nil
}
