package zendesk

type Tag string

type TagsPayload struct {
	Tags []Tag `json:"tags"`
}

type TagsResponse struct {
	Tags []Tag `json:"tags"`
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/
type TicketingTagsService struct {
	c *client
}
