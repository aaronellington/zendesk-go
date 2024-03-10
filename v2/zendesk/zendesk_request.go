package zendesk

import (
	"fmt"
	"net/url"
	"strings"
)

type RequestQueryModifiers func(query *url.Values)

func WithCursorPaginationPageSize(pageSize int) RequestQueryModifiers {
	return func(query *url.Values) {
		query.Set(
			"page[size]",
			fmt.Sprintf("%d", pageSize),
		)
	}
}

func WithTimeBasedIncrementalExportPageSize(pageSize int) RequestQueryModifiers {
	return func(query *url.Values) {
		query.Set(
			"per_page",
			fmt.Sprintf("%d", pageSize),
		)
	}
}

func WithSideloads(sideloads []string) RequestQueryModifiers {
	return func(query *url.Values) {
		query.Set(
			"include",
			strings.Join(sideloads, ","),
		)
	}
}
