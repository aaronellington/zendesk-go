package zendesk

type LiveChatService struct {
	accounts               *LiveChatAccountsService
	agents                 *LiveChatAgentsService
	bans                   *LiveChatBansService
	chats                  *LiveChatChatsService
	departments            *LiveChatDepartmentsService
	goals                  *LiveChatGoalsService
	incrementalAgentEvents *LiveChatIncrementalAgentEventsService
	incrementalExports     *LiveChatIncrementalExportsService
	oauthClients           *LiveChatOauthClientsService
	oauthTokens            *LiveChatOauthTokensService
	roles                  *LiveChatRolesService
	routingSettings        *LiveChatRoutingSettingsService
	shortcuts              *LiveChatShortcutsService
	skills                 *LiveChatSkillsService
	triggers               *LiveChatTriggersService
	visitors               *LiveChatVisitorsService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/accounts/
func (s *LiveChatService) Accounts() *LiveChatAccountsService {
	return s.accounts
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/agents/
func (s *LiveChatService) Agents() *LiveChatAgentsService {
	return s.agents
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/bans/
func (s *LiveChatService) Bans() *LiveChatBansService {
	return s.bans
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/
func (s *LiveChatService) Chats() *LiveChatChatsService {
	return s.chats
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/departments/
func (s *LiveChatService) Departments() *LiveChatDepartmentsService {
	return s.departments
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/goals/
func (s *LiveChatService) Goals() *LiveChatGoalsService {
	return s.goals
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/incremental_agent_events_api/
func (s *LiveChatService) IncrementalAgentEvents() *LiveChatIncrementalAgentEventsService {
	return s.incrementalAgentEvents
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/incremental_export/
func (s *LiveChatService) IncrementalExports() *LiveChatIncrementalExportsService {
	return s.incrementalExports
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/oauth_clients/
func (s *LiveChatService) OauthClients() *LiveChatOauthClientsService {
	return s.oauthClients
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/oauth_tokens/
func (s *LiveChatService) OauthTokens() *LiveChatOauthTokensService {
	return s.oauthTokens
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/roles/
func (s *LiveChatService) Roles() *LiveChatRolesService {
	return s.roles
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/routing_settings/
func (s *LiveChatService) RoutingSettings() *LiveChatRoutingSettingsService {
	return s.routingSettings
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/shortcuts/
func (s *LiveChatService) Shortcuts() *LiveChatShortcutsService {
	return s.shortcuts
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/skills/
func (s *LiveChatService) Skills() *LiveChatSkillsService {
	return s.skills
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/triggers/
func (s *LiveChatService) Triggers() *LiveChatTriggersService {
	return s.triggers
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/visitors/
func (s *LiveChatService) Visitors() *LiveChatVisitorsService {
	return s.visitors
}
