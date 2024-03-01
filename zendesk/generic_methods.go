package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
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
]) Show(
	ctx context.Context,
	id ID,
) (
	SingleResponse,
	error,
) {
	return s.getSingle(
		ctx,
		fmt.Sprintf("/api/v2/%s/%v", s.apiName, id),
	)
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
	if _, ok := any(new(ListResponse)).(isCursorPagination); ok {
		if !query.Has("page[size]") {
			query.Set("page[size]", "100")
		}
	}

	return genericList(
		ctx,
		s.client,
		fmt.Sprintf("/api/v2/%s?%s", s.apiName, query.Encode()),
		pageHandler,
	)
}

func (s genericService[
	ID,
	SingleResponse,
	ListResponse,
]) Search(
	ctx context.Context,
	query string,
	pageHandler func(response ListResponse) error,
) error {
	q := url.Values{}
	q.Set("query", query)

	return genericList(
		ctx,
		s.client,
		fmt.Sprintf("/api/v2/%s/search?%s", s.apiName, q.Encode()),
		pageHandler,
	)
}

func (s genericService[
	ID,
	SingleResponse,
	ListResponse,
]) IncrementalExport(
	ctx context.Context,
	startTime time.Time,
	perPage uint,
	sideloads []string,
	pageHandler func(response ListResponse) error,
) error {
	query := url.Values{}
	query.Set("start_time", fmt.Sprintf("%d", startTime.Unix()))
	query.Set("per_page", fmt.Sprintf("%d", perPage))

	if len(sideloads) > 0 {
		query.Set("include", strings.Join(sideloads, ","))
	}

	return genericList(
		ctx,
		s.client,
		fmt.Sprintf("/api/v2/incremental/%s.json?%s", s.apiName, query.Encode()),
		pageHandler,
	)
}

func (s genericService[
	ID,
	SingleResponse,
	ListResponse,
]) Update(
	ctx context.Context,
	id ID,
	payload any,
) (
	SingleResponse,
	error,
) {
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

	return genericRequest[SingleResponse](s.client, request)
}

func (s genericService[
	ID,
	SingleResponse,
	ListResponse,
]) Create(
	ctx context.Context,
	payload any,
) (SingleResponse, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/%s", s.apiName),
		structToReader(payload),
	)
	if err != nil {
		return *new(SingleResponse), err
	}

	return genericRequest[SingleResponse](s.client, request)
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

func (s genericService[
	ID,
	SingleResponse,
	ListResponse,
]) getSingle(
	ctx context.Context,
	endpoint string,
) (SingleResponse, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		http.NoBody,
	)
	if err != nil {
		return *new(SingleResponse), err
	}

	return genericRequest[SingleResponse](s.client, request)
}

func genericRequest[T any](
	client *client,
	request *http.Request,
) (T, error) {
	target := *new(T)

	if err := client.ZendeskRequest(request, &target); err != nil {
		return *new(T), err
	}

	return target, nil
}

func genericList[T paginationResponse](
	ctx context.Context,
	client *client,
	endpoint string,
	pageHandler func(page T) error,
) error {
	for {
		// Create the Request
		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		// Make the Request
		page, err := genericRequest[T](client, request)
		if err != nil {
			return err
		}

		// Give the caller the page data
		if err := pageHandler(page); err != nil {
			return err
		}

		// Calculate the next page
		nextPage := page.nextPage()

		// Calculate the next page for incremental exports
		if strings.Contains(request.URL.Path, "/incremental/") {
			if incrementalExportResponse, ok := any(page).(isIncrementalExport); ok {
				nextPage = incrementalExportResponse.isIncrementalExportNextPage(request.URL)
			}
		}

		if nextPage == "" {
			break
		}

		endpoint = nextPage
	}

	return nil
}
