package zendesk

type (
	AuditLogID     uint64
	AuditLogAction string
	SourceID       int64
)

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
