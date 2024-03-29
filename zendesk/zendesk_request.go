package zendesk

import (
	"fmt"
	"net/url"
	"strings"
	"time"
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

func WithPageSize(pageSize uint8) RequestQueryModifiers {
	return func(query *url.Values) {
		query.Set("page[size]", fmt.Sprintf("%d", pageSize))
	}
}

func WithSort(field string, direction string) RequestQueryModifiers {
	return func(query *url.Values) {
		query.Set("sort", fmt.Sprintf("%s%s", direction, field))
	}
}

func WithFilterForAction(action AuditLogAction) RequestQueryModifiers {
	return func(query *url.Values) {
		query.Add("filter[action]", string(action))
	}
}

func WithFilterForActorID(actorID UserID) RequestQueryModifiers {
	return func(query *url.Values) {
		query.Add("filter[actor_id]", fmt.Sprintf("%d", actorID))
	}
}

func WithFilterForCreatedAt(startTime time.Time, endTime time.Time) RequestQueryModifiers {
	return func(query *url.Values) {
		query.Add("filter[created_at][]", startTime.Format(timeFormat))
		query.Add("filter[created_at][]", endTime.Format(timeFormat))
	}
}

func WithFilterForIPAddress(ipAddress string) RequestQueryModifiers {
	return func(query *url.Values) {
		query.Add("filter[ip_address]", ipAddress)
	}
}

func WithFilterForSourceType(sourceType string) RequestQueryModifiers {
	return func(query *url.Values) {
		query.Add("filter[source_type]", sourceType)
	}
}

// Filter audit logs by the source id. Requires filter[source_type] to also be set.
func WithFilterForSourceID(sourceType string, sourceID uint64) RequestQueryModifiers {
	return func(query *url.Values) {
		query.Add("filter[source_type]", sourceType)
		query.Add("filter[source_id]", fmt.Sprintf("%d", sourceID))
	}
}
