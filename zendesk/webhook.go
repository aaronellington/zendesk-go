package zendesk

import "time"

type WebhookEventType string

const (
	WebhookEventUserActive WebhookEventType = "zen:event-type:user.active_changed"
	// Other webhook events...
)

type WebhookEvent struct {
	Type                WebhookEventType `json:"type"`
	AccountID           AccountID        `json:"account_id"`
	ID                  WebhookEventID   `json:"id"`
	Time                time.Time        `json:"time"`
	ZendeskEventVersion string           `json:"zendesk_event_version"`
	Subject             string           `json:"subject"`
	Detail              any              `json:"detail"`
	Event               any              `json:"event"`
}

type WebhookEventDetailArticle struct {
	BrandID BrandID   `json:"brand_id"`
	ID      ArticleID `json:"id"`
}

type WebhookEventEvent struct {
}
