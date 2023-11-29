package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/brands/
type BrandService struct {
	client *client
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
	CursorPaginationResponse
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
) (Brand, error) {
	target := BrandResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/brands/%d", id),
		http.NoBody,
	)
	if err != nil {
		return Brand{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return Brand{}, err
	}

	return target.Brand, nil
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/brands/#list-brands
func (s BrandService) List(
	ctx context.Context,
	pageHandler func(response BrandsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf(
		"/api/v2/brands?%s",
		query.Encode(),
	)

	for {
		target := BrandsResponse{}

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
