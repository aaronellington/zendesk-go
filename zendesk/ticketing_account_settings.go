package zendesk

import (
	"context"
	"net/http"
)

type AccountID uint64

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/account_settings/#json-format
type AccountSettings struct {
	Settings struct {
		// https://developer.zendesk.com/api-reference/ticketing/account-configuration/account_settings/#branding
		Branding struct {
			HeaderColor         string `json:"header_color"`
			TabBackgroundColor  string `json:"tab_background_color"`
			PageBackgroundColor string `json:"page_background_color"`
			TextColor           string `json:"text_color"`
			HeaderLogoURL       string `json:"header_logo_url"`
			FaviconURL          string `json:"favicon_url"`
		} `json:"branding"`
	} `json:"settings"`
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/account_settings/
type TicketingAccountSettingsService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/account_settings/#show-settings
func (s *TicketingAccountSettingsService) Show(ctx context.Context) (AccountSettings, error) {
	return genericRequest[AccountSettings](
		s.c,
		ctx,
		http.MethodGet,
		"/api/v2/account/settings",
		http.NoBody,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/account_settings/#update-account-settings
func (s *TicketingAccountSettingsService) Update(ctx context.Context, payload AccountSettingsPayload) (AccountSettings, error) {
	return genericRequest[AccountSettings](
		s.c,
		ctx,
		http.MethodPut,
		"/api/v2/account/settings",
		payload,
	)
}

type AccountSettingsPayload struct {
	Settings any `json:"settings"`
}
