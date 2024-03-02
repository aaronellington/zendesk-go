package zendesk

import (
	"context"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/brands/
type BrandService struct {
	client  *client
	generic genericService[
		BrandID,
		BrandResponse,
		BrandsResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/brands/#json-format
type Brand struct {
	Active            bool                `json:"active"`
	BrandURL          string              `json:"brand_url"`
	CreatedAt         time.Time           `json:"created_at"`
	Default           bool                `json:"default"`
	HasHelpCenter     bool                `json:"has_help_center"`
	HelpCenterState   string              `json:"help_center_state"`
	HostMapping       string              `json:"host_mapping"`
	ID                BrandID             `json:"id"`
	IsDeleted         bool                `json:"is_deleted"`
	Logo              BrandLogoAttachment `json:"logo"`
	Name              string              `json:"name"`
	SignatureTemplate string              `json:"signature_template"`
	Subdomain         string              `json:"subdomain"`
	TicketFormIDs     []TicketFormID      `json:"ticket_form_ids"`
	UpdatedAt         *time.Time          `json:"updated_at"`
	URL               string              `json:"url"`
}

type BrandsResponse struct {
	Brands []Brand `json:"brands"`
	cursorPaginationResponse
}

type BrandResponse struct {
	Brand Brand `json:"brand"`
}

type BrandLogoAttachment struct {
	ContentType      string       `json:"content_type"`
	ContentURL       string       `json:"content_url"`
	Deleted          bool         `json:"deleted"`
	FileName         string       `json:"file_name"`
	Height           uint64       `json:"height"`
	ID               AttachmentID `json:"id"`
	Inline           bool         `json:"inline"`
	MappedContentURL string       `json:"mapped_content_url"`
	Size             uint64       `json:"size"`
	URL              string       `json:"url"`
	Width            uint64       `json:"width"`
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/brands/#show-a-brand
func (s BrandService) Show(
	ctx context.Context,
	id BrandID,
) (BrandResponse, error) {
	return s.generic.Show(ctx, id)
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/brands/#list-brands
func (s BrandService) List(
	ctx context.Context,
	pageHandler func(response BrandsResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}
