package zendesk

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/account_settings/
type AccountConfigurationService struct {
	auditLogService   *AuditLogService
	customRoleService *CustomRoleService
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/audit_logs/
func (s *AccountConfigurationService) AuditLogs() *AuditLogService {
	return s.auditLogService
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/
func (s *AccountConfigurationService) CustomRoles() *CustomRoleService {
	return s.customRoleService
}
