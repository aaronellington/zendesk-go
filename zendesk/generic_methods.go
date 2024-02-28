package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type genericService[
	ID comparable,
	SingleResponse any,
	ListResponse paginationResponse,
] struct {
	client  *client
	apiName string
}

func (s genericService[
	ID,
	SingleResponse,
	ListResponse,
]) Show(ctx context.Context, id ID) (SingleResponse, error) {
	target := *new(SingleResponse)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/%s/%v", s.apiName, id),
		http.NoBody,
	)
	if err != nil {
		return target, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return target, err
	}

	return target, nil
}

func (s genericService[
	ID,
	SingleResponse,
	ListResponse,
]) List(
	ctx context.Context,
	pageHandler func(response ListResponse) error,
) error {
	query := url.Values{}

	x := any(new(ListResponse))
	if _, ok := x.(isCursorPagination); ok {
		if !query.Has("page[size]") {
			query.Set("page[size]", "100")
		}
	}

	endpoint := fmt.Sprintf("/api/v2/%s?%s", s.apiName, query.Encode())

	for {
		target := *new(ListResponse)

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

		nextPage := target.nextPage()

		if nextPage == nil {
			break
		}

		endpoint = *nextPage
	}

	return nil
}

func (s genericService[
	ID,
	SingleResponse,
	ListResponse,
]) Update(
	ctx context.Context,
	id ID,
	payload any,
) (SingleResponse, error) {
	target := *new(SingleResponse)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/%s/%v", s.apiName, id),
		structToReader(payload),
	)
	if err != nil {
		return target, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return target, err
	}

	return target, nil
}

func (s genericService[
	ID,
	SingleResponse,
	ListResponse,
]) Create(
	ctx context.Context,
	payload any,
) (SingleResponse, error) {
	target := *new(SingleResponse)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/%s", s.apiName),
		structToReader(payload),
	)
	if err != nil {
		return target, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return target, err
	}

	return target, nil
}

func (s genericService[
	ID,
	SingleResponse,
	ListResponse,
]) Delete(ctx context.Context, id ID) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/%s/%v", s.apiName, id),
		http.NoBody,
	)
	if err != nil {
		return err
	}

	return s.client.ZendeskRequest(request, nil)
}
