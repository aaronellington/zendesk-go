package zendesk

type paginationResponse interface {
	zendeskEntity
	nextPageEndpoint() string
}

type isCursorPagination interface {
	isCursorPagination()
}

type cursorPaginationResponse struct {
	Meta struct {
		HasMore      bool   `json:"has_more"`
		AfterCursor  string `json:"after_cursor"`
		BeforeCursor string `json:"before_cursor"`
	} `json:"meta"`
	Links struct {
		First string `json:"first"`
		Last  string `json:"last"`
		Next  string `json:"next"`
	} `json:"links"`
}

func (r cursorPaginationResponse) nextPageEndpoint() string {
	if !r.Meta.HasMore {
		return ""
	}

	return r.Links.Next
}

func (r cursorPaginationResponse) isCursorPagination() {}

type incrementalExportResponse struct {
	Count       int    `json:"count"`
	EndOfStream bool   `json:"end_of_stream"`
	EndTime     int64  `json:"end_time"`
	NextPage    string `json:"next_page"`
}

func (r incrementalExportResponse) nextPageEndpoint() string {
	if r.EndOfStream {
		return ""
	}

	return r.NextPage
}
