package zendesk

// https://developer.zendesk.com/api-reference/introduction/pagination/#using-cursor-pagination
type CursorPaginationResponse struct {
	Meta  CursorPaginationMeta  `json:"meta"`
	Links CursorPaginationLinks `json:"links"`
}

type CursorPaginationMeta struct {
	HasMore      bool   `json:"has_more"`
	AfterCursor  string `json:"after_cursor"`
	BeforeCursor string `json:"before_cursor"`
}

type CursorPaginationLinks struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Next  string `json:"next"`
}

type CursorPaginationSortDirection string

const (
	Asc  CursorPaginationSortDirection = ""
	Desc CursorPaginationSortDirection = "-"
)

// https://developer.zendesk.com/api-reference/introduction/pagination/#using-offset-pagination
type OffsetPaginationResponse struct {
	NextPage     *string `json:"next_page"`
	PreviousPage *string `json:"previous_page"`
	Count        uint64  `json:"count"`
}
