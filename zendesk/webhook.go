package zendesk

import "time"

type WebhookEventType string

const (
	WebhookEventUserActive         WebhookEventType = "zen:event-type:user.active_changed"
	WebhookEventOmnichannelRouting WebhookEventType = "zen:event-type:omnichannel_config.omnichannel_routing_feature_changed"
	// Other webhook events...
)

// https://support.zendesk.com/hc/en-us/articles/4408839108378-Creating-webhooks-to-interact-with-third-party-systems#ariaid-title4
// NOTE: For Webhook Trigger or Automation Payloads, any structure can be defined by a Zendesk Administrator
type WebhookTriggerEvent any

// https://developer.zendesk.com/api-reference/webhooks/event-types/webhook-event-types/
type WebhookEvent struct {
	Type                WebhookEventType `json:"type"`
	AccountID           AccountID        `json:"account_id"`
	ID                  WebhookEventID   `json:"id"`
	Time                time.Time        `json:"time"`
	ZendeskEventVersion string           `json:"zendesk_event_version"`
	Subject             string           `json:"subject"`
	// Detail              any              `json:"detail"`
	// Event               any              `json:"event"`
}

type WebhookEventArticle struct {
	Detail WebhookEventDetailArticle `json:"detail"`

	WebhookEvent
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/article-events/#detail-object-properties
type WebhookEventDetailArticle struct {
	BrandID BrandID   `json:"brand_id"`
	ID      ArticleID `json:"id"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/article-events/#detail-object-properties
type WebhookEventDetailCommunityPost struct {
	BrandID BrandID `json:"brand_id"`
	PostID  PostID  `json:"post_id"`
	TopicID TopicID `json:"topic_id"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/article-events/#detail-object-properties
type WebhookEventDetailOrganization struct {
	CreatedAt      time.Time `json:"created_at"`
	ExternalID     *string   `json:"external_id"`
	GroupID        *string   `json:"group_id"`
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	SharedComments bool      `json:"shared_comments"`
	SharedTickets  bool      `json:"shared_tickets"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/article-events/#detail-object-properties
type WebhookEventDetailUser struct {
	CreatedAt      time.Time    `json:"created_at"`
	Email          string       `json:"email"`
	ExternalID     string       `json:"external_id"`
	DefaultGroupID string       `json:"default_group_id"`
	ID             string       `json:"id"`
	OrganizationID string       `json:"organization_id"`
	Role           CustomRoleID `json:"role"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/article-events/#detail-object-properties
type WebhookEventDetailAgentAvailability struct {
	AccountID AccountID `json:"account_id"`
	AgentID   string    `json:"agent_id"`
	Version   string    `json:"version"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/article-events/#detail-object-properties
type WebhookEventDetailOmnichannelRoutingConfiguration struct {
	AccountID AccountID `json:"account_id"`
}

type WebhookEventEvent struct {
}
