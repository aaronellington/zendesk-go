package zendesk

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/
type OrganizationFieldService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#json-format
type OrganizationField struct {
	Active      bool                `json:"active"`
	Description string              `json:"description"`
	ID          OrganizationFieldID `json:"id"`
}

type OrganizationFieldsResponse struct {
	OrganizationFields []OrganizationField `json:"organization_fields"`
	CursorPaginationResponse
}

type OrganizationFieldResponse struct {
	OrganizationField
}
