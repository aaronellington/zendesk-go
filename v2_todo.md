func (s UserService) Create(ctx context.Context, payload UserPayload) (UserResponse, error)
func (s UserService) IncrementalExport(ctx context.Context, startTime int64, pageHandler func(response UsersIncrementalExportResponse) error) error
func (s UserService) IncrementalExportWithSideloads(ctx context.Context, startTime int64, sideloads []UserSideload, pageHandler func(response UsersIncrementalExportResponse) error) error
func (s UserService) Search(ctx context.Context, query string) (UsersResponse, error)
func (s UserService) SearchWithSideloads(ctx context.Context, query string, sideloads []UserSideload, pageHandler func(response UserSearchResponse) error) error
func (s UserService) Show(ctx context.Context, id UserID) (User, error)
func (s UserService) ShowSelf(ctx context.Context) (User, error)
func (s UserService) ShowWithSideloads(ctx context.Context, id UserID, sideloads []UserSideload) (UserResponse, error)
func (s UserService) Update(ctx context.Context, id UserID, payload UserPayload) (UserResponse, error)

func (s GroupMembershipService) Create(ctx context.Context, userID UserID, groupID GroupID) (GroupMembershipResponse, error)
func (s GroupMembershipService) Delete(ctx context.Context, userID UserID, groupMembershipID GroupMembershipID) error
func (s GroupMembershipService) List(ctx context.Context, pageHandler func(response GroupMembershipsResponse) error) error
func (s GroupMembershipService) ListByGroup(ctx context.Context, groupID GroupID, pageHandler func(response GroupMembershipsResponse) error) error
func (s GroupMembershipService) ListByUser(ctx context.Context, userID UserID, pageHandler func(response GroupMembershipsResponse) error) error
func (s GroupMembershipService) SetDefault(ctx context.Context, userID UserID, groupMembershipID GroupMembershipID) (GroupMembershipsResponse, error)
func (s GroupMembershipService) Show(ctx context.Context, id GroupMembershipID) (GroupMembership, error)

func (s OrganizationService) Autocomplete(ctx context.Context, term string, pageHandler func(response OrganizationAutocompleteResponse) error) error
func (s OrganizationService) Create(ctx context.Context, payload OrganizationPayload) (OrganizationResponse, error)
func (s OrganizationService) IncrementalExport(ctx context.Context, startTime int64, pageHandler func(response OrganizationsIncrementalExportResponse) error) error
func (s OrganizationService) Show(ctx context.Context, id OrganizationID) (Organization, error)
func (s OrganizationService) Update(ctx context.Context, id OrganizationID, payload OrganizationPayload) (OrganizationResponse, error)

func (s *SuspendedTicketService) Delete(ctx context.Context, id SuspendedTicketID) error
func (s *SuspendedTicketService) RecoverMultiple(ctx context.Context, ids []SuspendedTicketID) error
func (s SuspendedTicketService) List(ctx context.Context, pageHandler func(response SuspendedTicketsResponse) error) error

func (s TicketService) Import(ctx context.Context, payload TicketPayload) (TicketResponse, error)
func (s TicketService) Merge(ctx context.Context, destination TicketID, payload MergeRequestPayload) (JobStatusResponse, error)

func (s TicketAttachmentService) DownloadToFile(ctx context.Context, contentURL string, filePath string) error
func (s TicketAttachmentService) Download(ctx context.Context, contentURL string, writer io.Writer) error
func (s TicketAttachmentService) Show(ctx context.Context, attachmentID AttachmentID) (TicketAttachment, error)
func (s TicketAttachmentService) Upload(ctx context.Context, localFilePath string, uploadToken UploadToken) (TicketAttachmentUploadResponse, error)
func (s TicketAttachmentService) UploadWithFilename(ctx context.Context, localFilePath string, filename string, uploadToken UploadToken) (TicketAttachmentUploadResponse, error)

func (s *AgentEventService) GetAgentStates(ctx context.Context) AgentStates
func (s *AgentEventService) IncrementalExport(ctx context.Context, startTime time.Time, pageHandler func(response AgentEventExportResponse) error) error
func (s *AgentEventService) UpdateAgentStates(ctx context.Context, defaultStateTime time.Time) error
func (s *ChatsService) IncrementalExport(ctx context.Context, startTime time.Time, pageHandler func(response ChatsIncrementalExportResponse) error) error
func (s *ChatsService) List(ctx context.Context, pageHandler func(page ChatsResponse) error) error
func (s *ChatsService) Search(ctx context.Context, query string, pageHandler func(page ChatsSearchResponse) error) error
func (s *ChatsService) Show(ctx context.Context, id ChatID) (Chat, error)
func (s *DepartmentService) List(ctx context.Context) ([]Department, error)
func (s *DepartmentService) Show(ctx context.Context, id GroupID) (Department, error)
func (s *MacroService) List(ctx context.Context, pageHandler func(response MacrosResponse) error) error
func (s *RealTimeChatRestService) GetAllChatMetrics(ctx context.Context) (ChatsStreamResponse, error)
func (s *RealTimeChatRestService) GetAllChatMetricsForDepartment(ctx context.Context, departmentID GroupID) (ChatsStreamResponse, error)
func (s *RealTimeChatRestService) GetAllChatMetricsForSpecificTimeWindow(ctx context.Context, timeWindow LiveChatTimeWindow) (ChatsStreamResponse, error)
func (s *RealTimeChatRestService) GetSingleChatMetric(ctx context.Context, chatMetric LiveChatMetricKeyChat) (ChatsStreamResponse, error)
func (s *RealTimeChatRestService) GetSingleChatMetricForDepartment(ctx context.Context, chatMetric LiveChatMetricKeyChat, departmentID GroupID) (ChatsStreamResponse, error)
func (s *RealTimeChatRestService) GetSingleChatMetricForSpecificTimeWindow(ctx context.Context, chatMetric LiveChatMetricKeyChat, timeWindow LiveChatTimeWindow) (ChatsStreamResponse, error)
func (s *ScheduleService) List(ctx context.Context) (SchedulesResponse, error)
func (s *SideConversationService) Create(ctx context.Context, ticketID TicketID, payload SideConversationCreatePayload) error
func (s *TriggerService) List(ctx context.Context, pageHandler func(response TriggersResponse) error) error
func (s *UserIdentityService) Create(ctx context.Context, userID UserID, payload UserIdentityPayload) (UserIdentityResponse, error)
func (s *UserIdentityService) List(ctx context.Context, userID UserID, pageHandler func(response UserIdentitiesResponse) error) error
func (s *ViewService) List(ctx context.Context, pageHandler func(response ViewsResponse) error) error
func (s *WebhookService) HandleWebhookEvent(eventHandlers WebhookEventHandlers, webhookSigningSecret string) http.Handler
func (s *WebhookService) HandleWebhookTrigger(handler func(ctx context.Context, webhookBody []byte) error, webhookSigningSecret string) http.Handler
func (s ArticleService) List(ctx context.Context, pageHandler func(response ArticlesResponse) error) error
func (s ArticleService) Show(ctx context.Context, id ArticleID) (Article, error)
func (s AuditLogService) List(ctx context.Context, pageHandler func(response AuditLogsResponse) error, modifiers ...ListAccountConfigurationAuditLogModifier) error
func (s AutomationService) List(ctx context.Context, pageHandler func(response AutomationsResponse) error) error
func (s BrandService) List(ctx context.Context, pageHandler func(response BrandsResponse) error) error
func (s BrandService) Show(ctx context.Context, id BrandID) (Brand, error)
func (s CategoryService) List(ctx context.Context, pageHandler func(response CategoriesResponse) error) error
func (s CustomRoleService) List(ctx context.Context, pageHandler func(response CustomRolesResponse) error) error
func (s CustomRoleService) Show(ctx context.Context, id CustomRoleID) (CustomRole, error)
func (s CustomStatusService) Create(ctx context.Context, payload CustomStatusPayload) (CustomStatusResponse, error)
func (s CustomStatusService) List(ctx context.Context, pageHandler func(response CustomStatusesResponse) error) error
func (s CustomStatusService) Show(ctx context.Context, id CustomStatusID) (CustomStatus, error)
func (s GroupsService) Create(ctx context.Context, payload GroupPayload) (GroupResponse, error)
func (s GroupsService) List(ctx context.Context, pageHandler func(response GroupsResponse) error) error
func (s GroupsService) Show(ctx context.Context, id GroupID) (Group, error)
func (s OAuthClientService) Create(ctx context.Context, payload any) (OAuthClientConfiguration, error)
func (s OAuthClientService) Delete(ctx context.Context, id LiveChatOAuthClientID) (OAuthClientConfiguration, error)
func (s OAuthClientService) GenerateOAuthClientSecret(ctx context.Context, id LiveChatOAuthClientID) (OAuthClientConfiguration, error)
func (s OAuthClientService) List(ctx context.Context) ([]OAuthClientConfiguration, error)
func (s OAuthClientService) Show(ctx context.Context, id LiveChatOAuthClientID) (OAuthClientConfiguration, error)
func (s OAuthClientService) Update(ctx context.Context, id LiveChatOAuthClientID, payload any) (OAuthClientConfiguration, error)
func (s OrganizationFieldService) Create(ctx context.Context, payload OrganizationFieldPayload) (OrganizationFieldConfigurationResponse, error)
func (s OrganizationFieldService) Delete(ctx context.Context, id OrganizationFieldID) error
func (s OrganizationFieldService) List(ctx context.Context, pageHandler func(response OrganizationFieldsConfigurationResponse) error) error
func (s OrganizationFieldService) Show(ctx context.Context, id OrganizationFieldID) (OrganizationFieldConfiguration, error)
func (s OrganizationFieldService) Update(ctx context.Context, id OrganizationFieldID, payload OrganizationFieldPayload) (OrganizationFieldConfigurationResponse, error)
func (s OrganizationMembershipService) Create(ctx context.Context, payload OrganizationMembershipPayload) (OrganizationMembershipResponse, error)
func (s OrganizationMembershipService) List(ctx context.Context, pageHandler func(response OrganizationMembershipsResponse) error) error
func (s OrganizationMembershipService) Show(ctx context.Context, id OrganizationMembershipID) (OrganizationMembership, error)
func (s SatisfactionRatingService) List(ctx context.Context, pageHandler func(response SatisfactionRatingsResponse) error) error
func (s SatisfactionRatingService) ListWithModifiers(ctx context.Context, pageHandler func(response SatisfactionRatingsResponse) error, modifiers ...ListTicketSatisfactionRatingModifier) error
func (s SatisfactionRatingService) Show(ctx context.Context, id SatisfactionRatingID) (SatisfactionRating, error)
func (s SectionService) List(ctx context.Context, pageHandler func(response SectionsResponse) error) error
func (s SideConversationTargetChildTicket) SideConversationTarget()
func (s TicketAuditService) ListForTicket(ctx context.Context, ticketID TicketID, pageHandler func(response TicketAuditsResponse) error) error
func (s TicketCommentService) ListByTicketID(ctx context.Context, ticketID TicketID, pageHandler func(response TicketCommentResponse) error) error
func (s TicketCommentService) ListByTicketIDWithSideload(ctx context.Context, ticketID TicketID, sideloads []TicketCommentSideload, pageHandler func(response TicketCommentResponse) error) error
func (s TicketFieldService) List(ctx context.Context, pageHandler func(response TicketFieldsConfigurationResponse) error) error
func (s TicketFieldService) Show(ctx context.Context, id TicketFieldID) (TicketFieldConfiguration, error)
func (s TicketFormService) List(ctx context.Context) ([]TicketForm, error)
func (s TicketFormService) Show(ctx context.Context, ticketFormID TicketFormID) (TicketForm, error)
func (s TicketTagService) List(ctx context.Context, pageHandler func(response TagsResponse) error) error
func (s TicketTagService) Search(ctx context.Context, searchTerm string, pageHandler func(response TagSearchResponse) error) error
func (s UserFieldService) List(ctx context.Context, pageHandler func(response UserFieldsConfigurationResponse) error) error
func (s UserFieldService) Show(ctx context.Context, id UserFieldID) (UserFieldConfiguration, error)
func (schedule Schedule) Active(now time.Time) (bool, error)
func (schedule Schedule) Location() (*time.Location, error)
