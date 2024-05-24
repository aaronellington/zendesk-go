package zendesk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"unicode/utf8"
)

const minimumTagNameForSearch = int(2)

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags
type TicketTagService struct {
	client *client
}

type Tag string

type TagsPayload struct {
	Tags Tags `json:"tags"`
}

type Tags []Tag

func (tags Tags) HasTag(targetTag Tag) bool {
	for _, tag := range tags {
		if tag == targetTag {
			return true
		}
	}

	return false
}

type TagMeta struct {
	Name  Tag    `json:"name"`
	Count uint64 `json:"count"`
}

type TagsSearchResponse struct {
	Tags Tags `json:"tags"`
}

type TagsResponse struct {
	Tags []TagMeta `json:"tags"`
	CursorPaginationResponse
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#list-tags
func (s TicketTagService) List(
	ctx context.Context,
	pageHandler func(response TagsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf(
		"/api/v2/tags?%s",
		query.Encode(),
	)

	for {
		target := TagsResponse{}

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

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#search-tags
func (s TicketTagService) Search(
	ctx context.Context,
	searchTerm string,
) (Tags, error) {
	target := TagsSearchResponse{}

	if validSearchTermSize := utf8.RuneCountInString(searchTerm) >= minimumTagNameForSearch; !validSearchTermSize {
		return Tags{}, errors.New("invalid request - 'searchTerm' must be at least 2 characters")
	}

	query := url.Values{}
	query.Set("name", searchTerm)
	endpoint := fmt.Sprintf(
		"/api/v2/autocomplete/tags?%s",
		query.Encode(),
	)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		http.NoBody,
	)
	if err != nil {
		return Tags{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return Tags{}, err
	}

	return target.Tags, nil
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#set-tags
func (s TicketService) SetTags(ctx context.Context, ticketID TicketID, tags Tags) (Tags, error) {
	target := TagsPayload{}

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(TagsPayload{
		Tags: tags,
	}); err != nil {
		return Tags{}, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/tickets/%d/tags", ticketID),
		payloadBuf,
	)
	if err != nil {
		return Tags{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return Tags{}, err
	}

	return target.Tags, nil
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#add-tags
func (s TicketService) AddTags(ctx context.Context, ticketID TicketID, tags Tags) (Tags, error) {
	target := TagsPayload{}

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(TagsPayload{
		Tags: tags,
	}); err != nil {
		return Tags{}, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/tickets/%d/tags", ticketID),
		payloadBuf,
	)
	if err != nil {
		return Tags{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return Tags{}, err
	}

	return target.Tags, nil
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#remove-tags
func (s TicketService) RemoveTags(ctx context.Context, ticketID TicketID, tags Tags) (Tags, error) {
	target := TagsPayload{}

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(TagsPayload{
		Tags: tags,
	}); err != nil {
		return Tags{}, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/tickets/%d/tags", ticketID),
		payloadBuf,
	)
	if err != nil {
		return Tags{}, err
	}

	if err := s.client.ZendeskRequest(
		request,
		&target,
	); err != nil {
		return Tags{}, err
	}

	return target.Tags, nil
}
