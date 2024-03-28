package zendesk

import "time"

type (
	AuditLogID     uint64
	AuditLogAction string
	SourceID       int64
)

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/audit_logs/#json-format
type AuditLog struct {
	Action            AuditLogAction `json:"action"`
	ActionLabel       string         `json:"action_label"`
	ActorID           UserID         `json:"actor_id"`
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

const (
	Create   AuditLogAction = "create"
	Destroy  AuditLogAction = "destroy"
	Exported AuditLogAction = "exported"
	Login    AuditLogAction = "login"
	Update   AuditLogAction = "update"
)

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/audit_logs/
type TicketingAuditLogsService struct {
	c *client
}
