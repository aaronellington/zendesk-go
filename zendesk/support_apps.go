package zendesk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/apps/apps/
type AppService struct {
	client *client
}

type AppsResponse struct {
	Apps []App `json:"apps"`
}

// https://developer.zendesk.com/api-reference/ticketing/apps/apps/#json-format
type App struct {
	ID                          AppID           `json:"id"`
	OwnerID                     AppDeveloperID  `json:"owner_id"`
	AppOrganization             AppOrganization `json:"app_organization"`
	Name                        string          `json:"name"`
	SingleInstall               bool            `json:"single_install"`
	DefaultLocale               string          `json:"default_locale"`
	AuthorName                  string          `json:"author_name"`
	AuthorEmail                 string          `json:"author_email"`
	AuthorURL                   string          `json:"author_url"`
	RemoteInstallationURL       string          `json:"remote_installation_url"`
	ShortDescription            string          `json:"short_description"`
	LongDescription             string          `json:"long_description"`
	RawLongDescription          string          `json:"raw_long_description"`
	InstallationInstructions    string          `json:"installation_instructions"`
	RawInstallationInstructions string          `json:"raw_installation_instructions"`
	SmallIcon                   string          `json:"small_icon"`
	LargeIcon                   string          `json:"large_icon"`
	Screenshots                 []string        `json:"screenshots"`
	Visibility                  string          `json:"visibility"`
	Installable                 bool            `json:"installable"`
	CreatedAt                   time.Time       `json:"created_at"`
	UpdatedAt                   time.Time       `json:"updated_at"`
	FrameworkVersion            string          `json:"framework_version"`
	Featured                    bool            `json:"featured"`
	Promoted                    bool            `json:"promoted"`
	Products                    []string        `json:"products"`
	Categories                  []AppCategory   `json:"categories"`
	Version                     string          `json:"version"`
	MarketingOnly               bool            `json:"marketing_only"`
	Deprecated                  bool            `json:"deprecated"`
	Obsolete                    bool            `json:"obsolete"`
	Locations                   []AppLocationID `json:"locations"`
	Paid                        bool            `json:"paid"`
	Rating                      AppRating       `json:"rating"`
	FeatureColor                string          `json:"feature_color"`
	GoogleAnalyticsCode         string          `json:"google_analytics_code"`
	State                       string          `json:"state"`
	ClosedPreview               bool            `json:"closed_preview"`
	TermsConditionsURL          string          `json:"terms_conditions_url"`
	Collections                 []AppCollection `json:"collections"`
	UploadID                    uint64          `json:"upload_id"`
	Parameters                  []AppParameter  `json:"parameters"`
	InstallationCount           uint64          `json:"installation_count"`
}

type AppInstallation struct {
	ID              AppInstallationID  `json:"id"`
	AppID           AppID              `json:"app_id"`
	Product         string             `json:"product"`
	Settings        map[string]any     `json:"settings"`
	SettingsObjects []AppSettingObject `json:"settings_objects"`
	Enabled         bool               `json:"enabled"`
	Updated         uint64             `json:"updated"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`

	// "app_id": 1024831,
	// "product": "support",
	// "settings": {
	// 	"name": "Sour Owl",
	// 	"title": "Sour Owl"
	// },
	// "settings_objects": [
	// 	{
	// 		"name": "name",
	// 		"value": "Sour Owl"
	// 	},
	// 	{
	// 		"name": "title",
	// 		"value": "Sour Owl"
	// 	}
	// ],
	// "enabled": true,
	// "updated": "20240513225436",
	// "updated_at": "2024-03-30T16:20:08Z",
	// "created_at": "2024-03-30T16:20:08Z",
	// "role_restrictions": null,
	// "recurring_payment": false,
	// "collapsible": true,
	// "plan_information": {
	// 	"name": null
	// },
	// "paid": null,
	// "has_unpaid_subscription": null,
	// "has_incomplete_subscription": false,
	// "stripe_publishable_key": null,
	// "stripe_account": null,
	// "stripe_subscription_id": null,
	// "group_restrictions": []
}

type AppSettingObject struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type AppCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AppCollection struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AppOrganization struct {
	ID              any    `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Website         string `json:"website"`
	CountryCode     string `json:"country_code"`
	StripeAccountID string `json:"stripe_account_id"`
}

type (
	AppParameterID    uint
	AppInstallationID uint
)

// https://developer.zendesk.com/api-reference/ticketing/apps/app_locations/#json-format
type AppLocationID uint

const (
	ZendeskTopBar AppLocationID = iota + 1
	ZendeskNavBar
	ZendeskTicketSidebar
	ZendeskNewTicketSidebar
	ZendeskUserSidebar
	ZendeskOrganizationSidebar
	ZendeskBackground
	ZendeskChatSidebar
	ZendeskModal
	ZendeskTicketEditor
	_
	ZendeskSystemTopBar
	_
	ZopimBackground
	SellDealCard
	SellPersonCard
	SellCompanyCard
	SellLeadCard
	SellBackground
	SellModal
	SellDashboard
	SellNoteEditor
	SellCallLogEditor
	SellEmailEditor
	SellTopBar
	SellVisitEditor
)

type AppRating struct {
	TotalCount uint            `json:"total_count"`
	Average    *uint           `json:"average"`
	Rating     map[string]uint `json:"rating"`
}

type AppParameter struct {
	ID           AppParameterID `json:"id"`
	AppID        AppID          `json:"app_id"`
	Name         string         `json:"name"`
	Kind         string         `json:"kind"`
	Required     bool           `json:"required"`
	Position     uint           `json:"position"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DefaultValue any            `json:"default_value"`
	Secure       bool           `json:"secure"`
}

/*
https://developer.zendesk.com/api-reference/ticketing/apps/apps/#list-all-apps

Does not support pagination.
*/
func (s AppService) ListAllApps(
	ctx context.Context,
) (AppsResponse, error) {
	target := AppsResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/api/v2/apps",
		http.NoBody,
	)
	if err != nil {
		return AppsResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return AppsResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/apps/apps/#get-app-public-key
func (s AppService) GetAppPublicKey(
	ctx context.Context,
	id AppID,
) (string, error) {
	response, err := s.client.ZendeskGetRequest(
		ctx,
		fmt.Sprintf("/api/v2/apps/%d/public_key", id),
	)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	publicKeyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(publicKeyBytes), nil
}

// https://developer.zendesk.com/api-reference/ticketing/apps/apps/#get-information-about-app
func (s AppService) GetInformationAboutApp(
	ctx context.Context,
	id AppID,
) (App, error) {
	target := App{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/apps/%d", id),
		http.NoBody,
	)
	if err != nil {
		return App{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return App{}, err
	}

	return target, nil
}
