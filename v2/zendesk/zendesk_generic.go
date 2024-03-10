package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type zendeskEntity interface {
	zendeskEntityName() string
}

func showRequest[ID ~uint64, T zendeskEntity](
	ctx context.Context,
	c *client,
	id ID,
) (T, error) {
	return genericRequest[T](
		c,
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"/api/v2/%s/%d",
			(*new(T)).zendeskEntityName(),
			id,
		),
		http.NoBody,
	)
}

func createRequest[T zendeskEntity](
	ctx context.Context,
	c *client,
	payload any,
) (T, error) {
	return genericRequest[T](
		c,
		ctx,
		http.MethodPut,
		fmt.Sprintf(
			"/api/v2/%s",
			(*new(T)).zendeskEntityName(),
		),
		payload,
	)
}

func updateRequest[ID ~uint64, T zendeskEntity](
	ctx context.Context,
	c *client,
	id ID,
	payload any,
) (T, error) {
	return genericRequest[T](
		c,
		ctx,
		http.MethodPost,
		fmt.Sprintf(
			"/api/v2/%s/%d",
			(*new(T)).zendeskEntityName(),
			id,
		),
		payload,
	)
}

func deleteRequest[ID ~uint64, T zendeskEntity](
	ctx context.Context,
	c *client,
	id ID,
) error {
	return c.request(
		ctx,
		http.MethodDelete,
		fmt.Sprintf(
			"/api/v2/%s/%d",
			(*new(T)).zendeskEntityName(),
			id,
		),
		http.NoBody,
		nil,
	)
}

func listRequest[T paginationResponse](
	ctx context.Context,
	c *client,
	pageHandler func(response T) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return paginatedRequest(
		c,
		ctx,
		fmt.Sprintf(
			"/api/v2/%s",
			(*new(T)).zendeskEntityName(),
		),
		pageHandler,
		requestQueryModifiers...,
	)
}

func incrementalExportRequest[T paginationResponse](
	ctx context.Context,
	c *client,
	startTime time.Time,
	pageHandler func(response T) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	requestQueryModifiers = append(requestQueryModifiers, func(query *url.Values) {
		query.Set(
			"start_time",
			fmt.Sprintf("%d", startTime.Unix()),
		)
	})

	return paginatedRequest(
		c,
		ctx,
		fmt.Sprintf(
			"/api/v2/incremental/%s",
			(*new(T)).zendeskEntityName(),
		),
		pageHandler,
		requestQueryModifiers...,
	)
}

func paginatedRequest[T paginationResponse](
	c *client,
	ctx context.Context,
	path string,
	pageHandler func(response T) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	query := url.Values{}
	for _, requestQueryModifier := range requestQueryModifiers {
		requestQueryModifier(&query)
	}

	if _, ok := any(*new(T)).(isCursorPagination); ok {
		if !query.Has("page[size]") {
			query.Set("page[size]", "100")
		}
	}

	nextEndpoint := fmt.Sprintf(
		"%s?%s",
		path,
		query.Encode(),
	)

	for {
		response, err := genericRequest[T](
			c,
			ctx,
			http.MethodGet,
			nextEndpoint,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := pageHandler(response); err != nil {
			return err
		}

		nextEndpoint = response.nextPageEndpoint()
		if nextEndpoint == "" {
			break
		}
	}

	return nil
}

func genericRequest[T any](
	c *client,
	ctx context.Context,
	method string,
	endpoint string,
	payload any,
) (T, error) {
	target := *new(T)

	if err := c.request(ctx, method, endpoint, payload, &target); err != nil {
		return *new(T), err
	}

	return target, nil
}
