package zendesk

import (
	"fmt"
	"net/url"
	"time"
)

type paginationResponse interface {
	nextPage() *string
}

type isCursorPagination interface {
	isCursorPagination()
}

type isIncrementalExport interface {
	isIncrementalExportNextPage(pre *url.URL) *string
}

// https://developer.zendesk.com/api-reference/introduction/pagination/#using-cursor-pagination
type cursorPaginationResponse struct {
	Meta  cursorPaginationMeta  `json:"meta"`
	Links cursorPaginationLinks `json:"links"`
}

func (r cursorPaginationResponse) nextPage() *string {
	if !r.Meta.HasMore {
		return nil
	}

	return &r.Links.Next
}

func (r cursorPaginationResponse) isCursorPagination() {}

type cursorPaginationMeta struct {
	HasMore      bool   `json:"has_more"`
	AfterCursor  string `json:"after_cursor"`
	BeforeCursor string `json:"before_cursor"`
}

type cursorPaginationLinks struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Next  string `json:"next"`
}

type cursorPaginationSortDirection string

const (
	Asc  cursorPaginationSortDirection = ""
	Desc cursorPaginationSortDirection = "-"
)

// https://developer.zendesk.com/api-reference/introduction/pagination/#using-offset-pagination
type offsetPaginationResponse struct {
	NextPage     *string `json:"next_page"`
	PreviousPage *string `json:"previous_page"`
	Count        uint64  `json:"count"`
}

func (r offsetPaginationResponse) nextPage() *string {
	return r.NextPage
}

type incrementalExportResponse struct {
	EndTimeUnix int64 `json:"end_time"`
	EndOfStream bool  `json:"end_of_stream"`
}

func (response incrementalExportResponse) EndTime() time.Time {
	return time.Unix(response.EndTimeUnix, 0)
}

func (response incrementalExportResponse) isIncrementalExportNextPage(previousEndpoint *url.URL) *string {
	if response.EndOfStream {
		return nil
	}

	q := previousEndpoint.Query()
	q.Set("start_time", fmt.Sprintf("%d", response.EndTime().Unix()))

	previousEndpoint.RawQuery = q.Encode()

	newURL := previousEndpoint.String()

	return &newURL
}
