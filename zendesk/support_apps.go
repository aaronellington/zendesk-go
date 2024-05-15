package zendesk

import "time"

// https://developer.zendesk.com/api-reference/ticketing/apps/apps/
type AppService struct {
	client *client
}

type AppAppsResponse struct {
	Apps []App `json:"apps"`
	CursorPaginationResponse
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
	Visibility                  any             `json:"visibility"`
	Installable                 bool            `json:"installable"`
	CreatedAt                   time.Time       `json:"created_at"`
	UpdatedAt                   time.Time       `json:"updated_at"`
	FrameworkVersion            string          `json:"framework_version"`
	Featured                    bool            `json:"featured"`
	Promoted                    bool            `json:"promoted"`
	Products                    []string        `json:"products"`
	Categories                  any             `json:"categories"`
	Version                     any             `json:"version"`
	MarketingOnly               bool            `json:"marketing_only"`
	Deprecated                  bool            `json:"deprecated"`
	Obsolete                    bool            `json:"obsolete"`
	Locations                   any             `json:"locations"`
	Paid                        bool            `json:"paid"`
	Rating                      any             `json:"rating"`
	FeatureColor                any             `json:"feature_color"`
	GoogleAnalyticsCode         any             `json:"google_analytics_code"`
	State                       any             `json:"state"`
	ClosedPreview               any             `json:"closed_preview"`
	TermsConditionsURL          any             `json:"terms_conditions_url"`
	Collections                 any             `json:"collections"`
	UploadID                    any             `json:"upload_id"`
	Parameters                  any             `json:"parameters"`
	InstallationCount           uint64          `json:"installation_count"`
}

type AppOrganization struct {
	ID              any    `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Website         string `json:"website"`
	CountryCode     string `json:"country_code"`
	StripeAccountID string `json:"stripe_account_id"`
}
