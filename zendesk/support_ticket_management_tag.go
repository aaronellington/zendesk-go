package zendesk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"unicode/utf8"
)

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags
type TicketTagService struct {
	client *client
}

type Tag string

type TagMeta struct {
	Name  Tag    `json:"name"`
	Count uint64 `json:"count"`
}

type TagsResponse struct {
	Tags []TagMeta `json:"tags"`
	CursorPaginationResponse
}

type TagSearchResponse struct {
	Tags []Tag `json:"tags"`
	OffsetPaginationResponse
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

/*
https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#search-tags

Does not support cursor pagination.
*/
func (s TicketTagService) Search(
	ctx context.Context,
	searchTerm string,
	pageHandler func(response TagSearchResponse) error,
) error {
	if utf8.RuneCountInString(searchTerm) < 3 {
		return errors.New("invalid searchterm - searchterm must be at least 2 characters")
	}

	query := url.Values{}
	query.Set("name", searchTerm)
	endpoint := fmt.Sprintf(
		"/api/v2/autocomplete/tags?%s",
		query.Encode(),
	)

	for {
		target := TagSearchResponse{}

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

		if target.NextPage != nil {
			endpoint = *target.NextPage

			continue
		}

		break
	}

	return nil
}
