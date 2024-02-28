package zendesk

type paginationResponse interface {
	nextPage() *string
}

type isCursorPagination interface {
	isCursorPagination()
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
