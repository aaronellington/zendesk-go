package zendesk

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/account_settings/
type AccountConfigurationService struct {
	auditLogService *AuditLogService
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/audit_logs/
func (s *AccountConfigurationService) AuditLogs() *AuditLogService {
	return s.auditLogService
}
