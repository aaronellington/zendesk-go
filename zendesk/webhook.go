package zendesk

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type WebhookService struct{}

func (s *WebhookService) HandleEvent(
	secret string,
	eventHandlers WebhookHandlers,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !verifyZendeskWebhookSignatureIsValid(r, secret) {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		baseTarget := WebhookEvent[any, any]{}
		if err := json.Unmarshal(bodyBytes, &baseTarget); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if baseTarget.Type == "" {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		var handlerError error
		switch baseTarget.Type {
		case "zen:event-type:article.published":
			handlerError = callIt[any, any](r, bodyBytes, eventHandlers.ArticlePublished)
		case "zen:event-type:article.publishedxxxxx":
			handlerError = callIt[any, any](r, bodyBytes, eventHandlers.ArticlePublished)
		case "zen:event-type:article.publishedx":
			handlerError = callIt[any, any](r, bodyBytes, eventHandlers.ArticlePublished)
		case "zen:event-type:article.publishedxx":
			handlerError = callIt[any, any](r, bodyBytes, eventHandlers.ArticlePublished)
		case "zen:event-type:article.publishedxxx":
			handlerError = callIt[any, any](r, bodyBytes, eventHandlers.ArticlePublished)
		}

		if handlerError != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (s *WebhookService) HandleTrigger(
	secret string,
	handler func(r *http.Request) error,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !verifyZendeskWebhookSignatureIsValid(r, secret) {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if err := handler(r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	})
}

func callIt[Event any, Detail any](
	r *http.Request,
	bodyBytes []byte,
	it func(ctx context.Context, eventData WebhookEvent[Event, Detail]) error,
) error {
	if it == nil {
		return nil
	}

	target := WebhookEvent[Event, Detail]{}

	if err := json.Unmarshal(bodyBytes, &target); err != nil {
		return err
	}

	return it(r.Context(), target)
}

// https://developer.zendesk.com/documentation/webhooks/verifying/
func verifyZendeskWebhookSignatureIsValid(
	r *http.Request,
	secret string,
) bool {
	expectedZendeskSignature := r.Header.Get("X-Zendesk-Webhook-Signature")
	zendeskSignatureTimestamp := r.Header.Get("X-Zendesk-Webhook-Signature-Timestamp")

	if expectedZendeskSignature == "" || zendeskSignatureTimestamp == "" {
		return false
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}

	// Replace the content on the body
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	contentToHash := []byte(zendeskSignatureTimestamp)
	contentToHash = append(contentToHash, bodyBytes...)

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write(contentToHash)

	actualZendeskSignature := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	return expectedZendeskSignature == actualZendeskSignature
}

type WebhookHandlers struct {
	AgentChannelCreated                    func(ctx context.Context, eventData WebhookEvent[any, any]) error
	AgentChannelDeleted                    func(ctx context.Context, eventData WebhookEvent[any, any]) error
	AgentMaxCapacityChanged                func(ctx context.Context, eventData WebhookEvent[any, any]) error
	AgentStateChanged                      func(ctx context.Context, eventData WebhookEvent[any, any]) error
	AgentUnifiedStateChanged               func(ctx context.Context, eventData WebhookEvent[any, any]) error
	AgentWorkItemAdded                     func(ctx context.Context, eventData WebhookEvent[any, any]) error
	AgentWorkItemRemoved                   func(ctx context.Context, eventData WebhookEvent[any, any]) error
	ArticleCommentChanged                  func(ctx context.Context, eventData WebhookEvent[any, any]) error
	ArticleCommentCreated                  func(ctx context.Context, eventData WebhookEvent[any, any]) error
	ArticleCommentPublished                func(ctx context.Context, eventData WebhookEvent[any, any]) error
	ArticleCommentUnpublished              func(ctx context.Context, eventData WebhookEvent[any, any]) error
	ArticlePublished                       func(ctx context.Context, eventData WebhookEvent[any, any]) error
	ArticleSubscriptionCreated             func(ctx context.Context, eventData WebhookEvent[any, any]) error
	ArticleUnpublished                     func(ctx context.Context, eventData WebhookEvent[any, any]) error
	ArticleVoteChanged                     func(ctx context.Context, eventData WebhookEvent[any, any]) error
	ArticleVoteCreated                     func(ctx context.Context, eventData WebhookEvent[any, any]) error
	ArticleVoteRemoved                     func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostChanged                   func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostCommentChanged            func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostCommentCreated            func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostCommentPublished          func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostCommentUnpublished        func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostCommentVoteChanged        func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostCommentVoteCreated        func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostCreated                   func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostPublished                 func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostSubscriptionCreated       func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostUnpublished               func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostVoteChanged               func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostVoteCreated               func(ctx context.Context, eventData WebhookEvent[any, any]) error
	CommunityPostVoteRemoved               func(ctx context.Context, eventData WebhookEvent[any, any]) error
	OmniChannelRoutingConfigFeatureChanged func(ctx context.Context, eventData WebhookEvent[any, any]) error
	OrganizationCreated                    func(ctx context.Context, eventData WebhookEvent[any, any]) error
	OrganizationCustomFieldChanged         func(ctx context.Context, eventData WebhookEvent[any, any]) error
	OrganizationDeleted                    func(ctx context.Context, eventData WebhookEvent[any, any]) error
	OrganizationExternalIDChanged          func(ctx context.Context, eventData WebhookEvent[any, any]) error
	OrganizationNameChanged                func(ctx context.Context, eventData WebhookEvent[any, any]) error
	OrganizationTagsChanged                func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserActiveChanged                      func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserAliasChanged                       func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserCreated                            func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserCustomFieldChanged                 func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserCustomRoleChanged                  func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserDefaultGroupChanged                func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserDeleted                            func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserDetailsChanged                     func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserExternalIDChanged                  func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserGroupMembershipCreated             func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserGroupMembershipDeleted             func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserIdentityChanged                    func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserIdentityCreated                    func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserIdentityDeleted                    func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserLastLoginChanged                   func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserMerged                             func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserNameChanged                        func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserNotesChanged                       func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserOnlyPrivateCommentsChanged         func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserOrganizationMembershipCreated      func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserOrganizationMembershipDeleted      func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserPasswordChanged                    func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserPhotoChanged                       func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserRoleChanged                        func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserSuspendedChanged                   func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserTagsChanged                        func(ctx context.Context, eventData WebhookEvent[any, any]) error
	UserTimeZoneChanged                    func(ctx context.Context, eventData WebhookEvent[any, any]) error
}

type (
	WebhookEventID   string
	WebhookEventType string
)

// https://developer.zendesk.com/api-reference/webhooks/event-types/webhook-event-types/#event-schema
type WebhookEvent[Event any, Detail any] struct {
	ID                  WebhookEventID   `json:"id"`
	Type                WebhookEventType `json:"type"`
	AccountID           AccountID        `json:"account_id"`
	Time                time.Time        `json:"time"`
	ZendeskEventVersion string           `json:"zendesk_event_version"`
	Subject             string           `json:"subject"`
	Event               Event            `json:"event"`
	Detail              Detail           `json:"detail"`
}
