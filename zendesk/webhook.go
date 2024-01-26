package zendesk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"time"
)

const (
	WebhookHeaderSignature          string = "X-Zendesk-Webhook-Signature"
	WebhookHeaderSignatureTimestamp string = "X-Zendesk-Webhook-Signature-Timestamp"
)

// https://developer.zendesk.com/api-reference/webhooks/webhooks-api/webhooks/
type WebhookService struct {
	client *client
}

// type UserEventHandler[K WebhookUserEventData] func(e WebhookEventUser[K]) error

type WebhookEventType string

const (
	WebhookEventUserActive               WebhookEventType = "zen:event-type:user.active_changed"
	WebhookEventOmnichannelConfigFeature WebhookEventType = "zen:event-type:omnichannel_config.omnichannel_routing_feature_changed"
	// Other webhook events...
)

// https://support.zendesk.com/hc/en-us/articles/4408839108378-Creating-webhooks-to-interact-with-third-party-systems#ariaid-title4
// NOTE: For Webhookss connected to Triggers or Automations, any structure can be defined by a Zendesk Administrator for the payload
type WebhookTriggerEvent any

func (s WebhookService) HandleWebhook(
	processor func(requestBody []byte) error,
	webhookSigningSecret string,
) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			// Verifying webhook requests is optional
			if webhookSigningSecret != "" {
				s.verifyZendeskWebhookSignature(w, r, webhookSigningSecret)
			}

			// TODO: Include authentication methods
			// https://developer.zendesk.com/documentation/webhooks/webhook-security-and-authentication/#webhook-authentication

			body, err := readWebhookBody(r)
			if err != nil {
				respondToWebhookRequest(
					w,
					http.StatusInternalServerError,
					"Could not read Webhook Request body",
				)

				return
			}

			if err := processor(body); err != nil {
				respondToWebhookRequest(
					w,
					http.StatusInternalServerError,
					"Error occurred while processing webhook request",
				)

				return
			}

			respondToWebhookRequest(
				w,
				http.StatusOK,
				"Successfully handled Webhook Request",
			)
		},
	)
}

// https://developer.zendesk.com/documentation/webhooks/verifying/
func (s WebhookService) verifyZendeskWebhookSignature(w http.ResponseWriter, r *http.Request, webhookSigningSecret string) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		respondToWebhookRequest(
			w,
			http.StatusBadRequest,
			"Bad Request",
		)
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	expectedZendeskSignature := r.Header.Get(WebhookHeaderSignature)
	zendeskSignatureTimestamp := r.Header.Get(WebhookHeaderSignatureTimestamp)
	if expectedZendeskSignature == "" || zendeskSignatureTimestamp == "" {
		respondToWebhookRequest(
			w,
			http.StatusBadRequest,
			"Bad Request",
		)

		return
	}

	actualZendeskSignature := buildZendeskSignature(zendeskSignatureTimestamp, bodyBytes, webhookSigningSecret)
	if expectedZendeskSignature != actualZendeskSignature {
		respondToWebhookRequest(
			w,
			http.StatusBadRequest,
			"Bad Request",
		)

		return
	}
}

func buildZendeskSignature(
	timestamp string,
	bodyBytes []byte,
	webhookSigningSecret string,
) string {
	content := []byte(timestamp)
	content = append(content, bodyBytes...)

	hash := hmac.New(sha256.New, []byte(webhookSigningSecret))
	hash.Write(content)

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func respondToWebhookRequest(w http.ResponseWriter, status int, message string) {
	encoder := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := encoder.Encode(message); err != nil {
		return
	}
}

func readWebhookBody(r *http.Request) ([]byte, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes, nil
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/webhook-event-types/
type WebhookEvent struct {
	Type                WebhookEventType `json:"type"`
	AccountID           AccountID        `json:"account_id"`
	ID                  WebhookEventID   `json:"id"`
	Time                time.Time        `json:"time"`
	ZendeskEventVersion string           `json:"zendesk_event_version"`
	Subject             string           `json:"subject"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/article-events/
type WebhookEventArticle struct {
	Detail WebhookEventDetailArticle `json:"detail"`
	Event  any                       `json:"event"`
	WebhookEvent
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/article-events/#detail-object-properties
type WebhookEventDetailArticle struct {
	BrandID BrandID   `json:"brand_id"`
	ID      ArticleID `json:"id"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/user-events/
type WebhookEventUser struct {
	WebhookEvent
	Event  any                    `json:"event"`
	Detail WebhookEventDetailUser `json:"detail"`
}

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

// type WebhookEventUser[Event WebhookUserEventData] struct {
// 	WebhookEvent
// 	Event  Event                  `json:"event"`
// 	Detail WebhookEventDetailUser `json:"detail"`
// }

// type WebhookUserEventData interface {
// 	WebhookEventUserActiveStatusChanged | WebhookEventUserDetailsChanged
// }

// // https://developer.zendesk.com/api-reference/webhooks/event-types/user-events/#detail-object-properties

// type WebhookEventUserActiveStatusChanged struct {
// 	Current  bool `json:"current"`
// 	Previous bool `json:"previous"`
// }

// type WebhookEventUserDetailsChanged struct {
// 	Current  string `json:"current"`
// 	Previous string `json:"previous"`
// }

// https://developer.zendesk.com/api-reference/webhooks/event-types/community-events/
type WebhookEventCommunityPost struct {
	Detail WebhookEventDetailCommunityPost `json:"detail"`
	Event  any                             `json:"event"`
	WebhookEvent
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/community-events/#detail-object-properties
type WebhookEventDetailCommunityPost struct {
	BrandID BrandID `json:"brand_id"`
	PostID  PostID  `json:"post_id"`
	TopicID TopicID `json:"topic_id"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/organization-events/
type WebhookEventOrganization struct {
	Detail WebhookEventDetailOrganization `json:"detail"`
	Event  any                            `json:"event"`
	WebhookEvent
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/organization-events/#detail-object-properties
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

// https://developer.zendesk.com/api-reference/webhooks/event-types/omnichannel-routing-configuration-events/
type WebhookEventOmnichannelRoutingConfig struct {
	Detail WebhookEventDetailOmnichannelRoutingConfig `json:"detail"`
	Event  any                                        `json:"event"`
	WebhookEvent
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/omnichannel-routing-configuration-events/#detail-object-properties
type WebhookEventDetailOmnichannelRoutingConfig struct {
	AccountID AccountID `json:"account_id"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/agent-availability-events/
type WebhookEventAgentAvailability struct {
	Detail WebhookEventDetailAgentAvailability `json:"detail"`
	Event  any                                 `json:"event"`
	WebhookEvent
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/agent-availability-events/#detail-object-properties
type WebhookEventDetailAgentAvailability struct {
	AccountID AccountID `json:"account_id"`
	AgentID   string    `json:"agent_id"`
	Version   string    `json:"version"`
}
