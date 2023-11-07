package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type AuditLogsResponse struct {
	AuditLogs []AuditLog `json:"audit_logs"`
	CursorPaginationResponse
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/audit_logs/#json-format
type AuditLog struct {
	Action            AuditLogAction `json:"action"`
	ActionLabel       string         `json:"action_label"`
	ActorID           ActorID        `json:"actor_id"`
	ActorName         string         `json:"actor_name"`
	ChangeDescription string         `json:"change_description"`
	CreatedAt         time.Time      `json:"created_at"`
	ID                AuditLogID     `json:"id"`
	IPAddress         *string        `json:"ip_address"`
	SourceID          SourceID       `json:"source_id"`
	SourceLabel       string         `json:"source_label"`
	SourceType        string         `json:"source_type"`
	URL               string         `json:"url"`
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/audit_logs/
type AuditLogService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/audit_logs/#list-audit-logs
func (s AuditLogService) List(
	ctx context.Context,
	pageHandler func(response AuditLogsResponse) error,
	modifiers ...ListAccountConfigurationAuditLogModifier,
) error {
	query := url.Values{}
	// Default values
	query.Set("page[size]", "100")
	query.Set("sort", "created_at")

	for _, modifier := range modifiers {
		modifier.ModifyListAccountConfigurationAuditLogRequest(&query)
	}

	endpoint := fmt.Sprintf("/api/v2/audit_logs.json?%s", query.Encode())

	for {
		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		target := AuditLogsResponse{}

		if err := s.client.ZendeskRequest(request, &target, false); err != nil {
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

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/audit_logs/#parameters
type ListAccountConfigurationAuditLogModifier interface {
	ModifyListAccountConfigurationAuditLogRequest(queryParameters *url.Values)
}

type listAccountConfigurationAuditLogModifier func(queryParameters *url.Values)

func (l listAccountConfigurationAuditLogModifier) ModifyListAccountConfigurationAuditLogRequest(queryParameters *url.Values) {
	l(queryParameters)
}

func WithPageSize(pageSize uint8) listAccountConfigurationAuditLogModifier {
	return listAccountConfigurationAuditLogModifier(func(queryParameters *url.Values) {
		queryParameters.Set("page[size]", fmt.Sprintf("%d", pageSize))
	})
}

func WithSort(field string, direction CursorPaginationSortDirection) listAccountConfigurationAuditLogModifier {
	return listAccountConfigurationAuditLogModifier(func(queryParameters *url.Values) {
		queryParameters.Set("sort", fmt.Sprintf("%s%s", direction, field))
	})
}

func WithFilterForAction(action AuditLogAction) listAccountConfigurationAuditLogModifier {
	return listAccountConfigurationAuditLogModifier(func(queryParameters *url.Values) {
		queryParameters.Add("filter[action]", string(action))
	})
}

func WithFilterForActorID(actorID ActorID) listAccountConfigurationAuditLogModifier {
	return listAccountConfigurationAuditLogModifier(func(queryParameters *url.Values) {
		queryParameters.Add("filter[actor_id]", fmt.Sprintf("%d", actorID))
	})
}

func WithFilterForCreatedAt(startTime time.Time, endTime time.Time) listAccountConfigurationAuditLogModifier {
	return listAccountConfigurationAuditLogModifier(func(queryParameters *url.Values) {
		queryParameters.Add("filter[created_at][]", startTime.Format(timeFormat))
		queryParameters.Add("filter[created_at][]", endTime.Format(timeFormat))
	})
}

func WithFilterForIPAddress(ipAddress string) listAccountConfigurationAuditLogModifier {
	return listAccountConfigurationAuditLogModifier(func(queryParameters *url.Values) {
		queryParameters.Add("filter[ip_address]", ipAddress)
	})
}

func WithFilterForSourceType(sourceType string) listAccountConfigurationAuditLogModifier {
	return listAccountConfigurationAuditLogModifier(func(queryParameters *url.Values) {
		queryParameters.Add("filter[source_type]", sourceType)
	})
}

// Filter audit logs by the source id. Requires filter[source_type] to also be set.
func WithFilterForSourceID(sourceType string, sourceID uint64) listAccountConfigurationAuditLogModifier {
	return listAccountConfigurationAuditLogModifier(func(queryParameters *url.Values) {
		queryParameters.Add("filter[source_type]", sourceType)
		queryParameters.Add("filter[source_id]", fmt.Sprintf("%d", sourceID))
	})
}
