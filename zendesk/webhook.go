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

const (
	WebhookHeaderSignature          string = "X-Zendesk-Webhook-Signature"
	WebhookHeaderSignatureTimestamp string = "X-Zendesk-Webhook-Signature-Timestamp"
)

// https://developer.zendesk.com/api-reference/webhooks/webhooks-api/webhooks/
type WebhookService struct {
	client *client
}

type WebhookEventType string

// https://developer.zendesk.com/api-reference/webhooks/event-types/article-events/
const (
	WebhookEventArticleAuthorChanged       WebhookEventType = "zen:event-type:article.author_changed"
	WebhookEventArticlePublished           WebhookEventType = "zen:event-type:article.published"
	WebhookEventArticleSubscriptionCreated WebhookEventType = "zen:event-type:article.subscription_created"
	WebhookEventArticleUnpublished         WebhookEventType = "zen:event-type:article.unpublished"
	WebhookEventArticleVoteCreated         WebhookEventType = "zen:event-type:article.vote_created"
	WebhookEventArticleVoteChanged         WebhookEventType = "zen:event-type:article.vote_changed"
	WebhookEventArticleVoteRemoved         WebhookEventType = "zen:event-type:article.vote_removed"
	WebhookEventArticleCommentCreated      WebhookEventType = "zen:event-type:article.comment_created"
	WebhookEventArticleCommentChanged      WebhookEventType = "zen:event-type:article.comment_changed"
	WebhookEventArticleCommentPublished    WebhookEventType = "zen:event-type:article.comment_published"
	WebhookEventArticleCommentUnpublished  WebhookEventType = "zen:event-type:article.comment_unpublished"
)

// https://developer.zendesk.com/api-reference/webhooks/event-types/community-events/
const (
	WebhookEventCommunityPostCreated             WebhookEventType = "zen:event-type:community_post.created"
	WebhookEventCommunityPostChanged             WebhookEventType = "zen:event-type:community_post.changed"
	WebhookEventCommunityPostPublished           WebhookEventType = "zen:event-type:community_post.published"
	WebhookEventCommunityPostUnpublished         WebhookEventType = "zen:event-type:community_post.unpublished"
	WebhookEventCommunityPostSubscriptionCreated WebhookEventType = "zen:event-type:community_post.subscription_created"
	WebhookEventCommunityPostVoteCreated         WebhookEventType = "zen:event-type:community_post.vote_created"
	WebhookEventCommunityPostVoteChanged         WebhookEventType = "zen:event-type:community_post.vote_changed"
	WebhookEventCommunityPostVoteRemoved         WebhookEventType = "zen:event-type:community_post.vote_removed"
	WebhookEventCommunityPostCommentCreated      WebhookEventType = "zen:event-type:community_post.comment_created"
	WebhookEventCommunityPostCommentChanged      WebhookEventType = "zen:event-type:community_post.comment_changed"
	WebhookEventCommunityPostCommentPublished    WebhookEventType = "zen:event-type:community_post.comment_published"
	WebhookEventCommunityPostCommentUnpublished  WebhookEventType = "zen:event-type:community_post.comment_unpublished"
	WebhookEventCommunityPostCommentVoteCreated  WebhookEventType = "zen:event-type:community_post.comment_vote_created"
	WebhookEventCommunityPostCommentVoteChanged  WebhookEventType = "zen:event-type:community_post.comment_vote_changed"
)

// https://developer.zendesk.com/api-reference/webhooks/event-types/organization-events/
const (
	WebhookEventOrganizationCreated            WebhookEventType = "zen:event-type:organization.created"
	WebhookEventOrganizationCustomFieldChanged WebhookEventType = "zen:event-type:organization.custom_field_changed"
	WebhookEventOrganizationDeleted            WebhookEventType = "zen:event-type:organization.deleted"
	WebhookEventOrganizationExternalIDChanged  WebhookEventType = "zen:event-type:organization.external_id_changed"
	WebhookEventOrganizationNameChanged        WebhookEventType = "zen:event-type:organization.name_changed"
	WebhookEventOrganizationTagsChanged        WebhookEventType = "zen:event-type:organization.tags_changed"
)

// https://developer.zendesk.com/api-reference/webhooks/event-types/user-events
const (
	WebhookEventUserAliasChanged                  WebhookEventType = "zen:event-type:user.alias_changed"
	WebhookEventUserCreated                       WebhookEventType = "zen:event-type:user.created"
	WebhookEventUserCustomFieldChanged            WebhookEventType = "zen:event-type:user.custom_field_changed"
	WebhookEventUserCustomRoleChanged             WebhookEventType = "zen:event-type:user.custom_role_changed"
	WebhookEventUserDefaultGroupChanged           WebhookEventType = "zen:event-type:user.default_group_changed"
	WebhookEventUserDetailsChanged                WebhookEventType = "zen:event-type:user.details_changed"
	WebhookEventUserExternalIDChanged             WebhookEventType = "zen:event-type:user.external_id_changed"
	WebhookEventUserGroupMembershipCreated        WebhookEventType = "zen:event-type:user.group_membership_created"
	WebhookEventUserGroupMembershipDeleted        WebhookEventType = "zen:event-type:user.group_membership_deleted"
	WebhookEventUserIdentityChanged               WebhookEventType = "zen:event-type:user.identity_changed"
	WebhookEventUserIdentityCreated               WebhookEventType = "zen:event-type:user.identity_created"
	WebhookEventUserIdentityDeleted               WebhookEventType = "zen:event-type:user.identity_deleted"
	WebhookEventUserActiveChanged                 WebhookEventType = "zen:event-type:user.active_changed"
	WebhookEventUserLastLoginChanged              WebhookEventType = "zen:event-type:user.last_login_changed"
	WebhookEventUserMerged                        WebhookEventType = "zen:event-type:user.merged"
	WebhookEventUserNameChanged                   WebhookEventType = "zen:event-type:user.name_changed"
	WebhookEventUserNotesChanged                  WebhookEventType = "zen:event-type:user.notes_changed"
	WebhookEventUserOnlyPrivateCommentsChanged    WebhookEventType = "zen:event-type:user.only_private_comments_changed"
	WebhookEventUserOrganizationMembershipCreated WebhookEventType = "zen:event-type:user.organization_membership_created"
	WebhookEventUserOrganizationMembershipDeleted WebhookEventType = "zen:event-type:user.organization_membership_deleted"
	WebhookEventUserPasswordChanged               WebhookEventType = "zen:event-type:user.password_changed" // #nosec G101 -- This is a false positive
	WebhookEventUserPhotoChanged                  WebhookEventType = "zen:event-type:user.photo_changed"
	WebhookEventUserRoleChanged                   WebhookEventType = "zen:event-type:user.role_changed"
	WebhookEventUserDeleted                       WebhookEventType = "zen:event-type:user.deleted"
	WebhookEventUserSuspendedChanged              WebhookEventType = "zen:event-type:user.suspended_changed"
	WebhookEventUserTagsChanged                   WebhookEventType = "zen:event-type:user.tags_changed"
	WebhookEventUserTimeZoneChanged               WebhookEventType = "zen:event-type:user.time_zone_changed"
)

// https://developer.zendesk.com/api-reference/webhooks/event-types/agent-availability-events/
const (
	WebhookEventAgentStateChanged        WebhookEventType = "zen:event-type:agent.state_changed"
	WebhookEventAgentWorkItemAdded       WebhookEventType = "zen:event-type:agent.work_item_added"
	WebhookEventAgentWorkItemRemoved     WebhookEventType = "zen:event-type:agent.work_item_removed"
	WebhookEventAgentMaxCapacityChanged  WebhookEventType = "zen:event-type:agent.max_capacity_changed"
	WebhookEventAgentUnifiedStateChanged WebhookEventType = "zen:event-type:agent.unified_state_changed"
	WebhookEventAgentChannelCreated      WebhookEventType = "zen:event-type:agent.channel_created"
	WebhookEventAgentChannelDeleted      WebhookEventType = "zen:event-type:agent.channel_deleted"
)

// https://developer.zendesk.com/api-reference/webhooks/event-types/omnichannel-routing-configuration-events/
const (
	WebhookEventOmnichannelRoutingConfigFeatureChanged WebhookEventType = "zen:event-type:omnichannel_config.omnichannel_routing_feature_changed"
)

const (
	WebhookEventTypePrefixAgentState        string = "zen:event-type:agent"
	WebhookEventTypePrefixUser              string = "zen:event-type:user"
	WebhookEventTypePrefixOrganization      string = "zen:event-type:organization"
	WebhookEventTypePrefixCommunityPost     string = "zen:event-type:community_post"
	WebhookEventTypePrefixOmnichannelConfig string = "zen:event-type:omnichannel_config"
	WebhookEventTypePrefixArticle           string = "zen:event-type:article"
)

type WebhookEventHandlers struct {
	WebhookEventArticlePublished           func(ctx context.Context, eventData WebhookEventArticlePublishedPayload) error
	WebhookEventArticleSubscriptionCreated func(ctx context.Context, eventData WebhookEventArticleSubscriptionCreatedPayload) error
	WebhookEventArticleUnpublished         func(ctx context.Context, eventData WebhookEventArticleUnpublishedPayload) error
	WebhookEventArticleVoteCreated         func(ctx context.Context, eventData WebhookEventArticleVoteCreatedPayload) error
	WebhookEventArticleVoteChanged         func(ctx context.Context, eventData WebhookEventArticleVoteChangedPayload) error
	WebhookEventArticleVoteRemoved         func(ctx context.Context, eventData WebhookEventArticleVoteRemovedPayload) error
	WebhookEventArticleCommentCreated      func(ctx context.Context, eventData WebhookEventArticleCommentCreatedPayload) error
	WebhookEventArticleCommentChanged      func(ctx context.Context, eventData WebhookEventArticleCommentChangedPayload) error
	WebhookEventArticleCommentPublished    func(ctx context.Context, eventData WebhookEventArticleCommentPublishedPayload) error
	WebhookEventArticleCommentUnpublished  func(ctx context.Context, eventData WebhookEventArticleCommentUnpublishedPayload) error
	WebhookEventArticleDefault             func(ctx context.Context, webhookBody []byte) error

	WebhookEventOrganizationCreated            func(ctx context.Context, eventData WebhookEventOrganizationCreatedPayload) error
	WebhookEventOrganizationCustomFieldChanged func(ctx context.Context, eventData WebhookEventOrganizationCustomFieldChangedPayload) error
	WebhookEventOrganizationDeleted            func(ctx context.Context, eventData WebhookEventOrganizationDeletedPayload) error
	WebhookEventOrganizationExternalIDChanged  func(ctx context.Context, eventData WebhookEventOrganizationExternalIDChangedPayload) error
	WebhookEventOrganizationNameChanged        func(ctx context.Context, eventData WebhookEventOrganizationNameChangedPayload) error
	WebhookEventOrganizationTagsChanged        func(ctx context.Context, eventData WebhookEventOrganizationTagsChangedPayload) error
	WebhookEventOrganizationDefault            func(ctx context.Context, webhookBody []byte) error

	WebhookEventUserAliasChanged                  func(ctx context.Context, eventData WebhookEventUserAliasChangedPayload) error
	WebhookEventUserCreated                       func(ctx context.Context, eventData WebhookEventUserCreatedPayload) error
	WebhookEventUserCustomFieldChanged            func(ctx context.Context, eventData WebhookEventUserCustomFieldChangedPayload) error
	WebhookEventUserCustomRoleChanged             func(ctx context.Context, eventData WebhookEventUserCustomRoleChangedPayload) error
	WebhookEventUserDefaultGroupChanged           func(ctx context.Context, eventData WebhookEventUserDefaultGroupChangedPayload) error
	WebhookEventUserDetailsChanged                func(ctx context.Context, eventData WebhookEventUserDetailsChangedPayload) error
	WebhookEventUserExternalIDChanged             func(ctx context.Context, eventData WebhookEventUserExternalIDChangedPayload) error
	WebhookEventUserGroupMembershipCreated        func(ctx context.Context, eventData WebhookEventUserGroupMembershipCreatedPayload) error
	WebhookEventUserGroupMembershipDeleted        func(ctx context.Context, eventData WebhookEventUserGroupMembershipDeletedPayload) error
	WebhookEventUserIdentityChanged               func(ctx context.Context, eventData WebhookEventUserIdentityChangedPayload) error
	WebhookEventUserIdentityCreated               func(ctx context.Context, eventData WebhookEventUserIdentityCreatedPayload) error
	WebhookEventUserIdentityDeleted               func(ctx context.Context, eventData WebhookEventUserIdentityDeletedPayload) error
	WebhookEventUserActiveChanged                 func(ctx context.Context, eventData WebhookEventUserActiveChangedPayload) error
	WebhookEventUserLastLoginChanged              func(ctx context.Context, eventData WebhookEventUserLastLoginChangedPayload) error
	WebhookEventUserMerged                        func(ctx context.Context, eventData WebhookEventUserMergedPayload) error
	WebhookEventUserNameChanged                   func(ctx context.Context, eventData WebhookEventUserNameChangedPayload) error
	WebhookEventUserNotesChanged                  func(ctx context.Context, eventData WebhookEventUserNotesChangedPayload) error
	WebhookEventUserOnlyPrivateCommentsChanged    func(ctx context.Context, eventData WebhookEventUserOnlyPrivateCommentsChangedPayload) error
	WebhookEventUserOrganizationMembershipCreated func(ctx context.Context, eventData WebhookEventUserOrganizationMembershipCreatedPayload) error
	WebhookEventUserOrganizationMembershipDeleted func(ctx context.Context, eventData WebhookEventUserOrganizationMembershipDeletedPayload) error
	WebhookEventUserPasswordChanged               func(ctx context.Context, eventData WebhookEventUserPasswordChangedPayload) error
	WebhookEventUserPhotoChanged                  func(ctx context.Context, eventData WebhookEventUserPhotoChangedPayload) error
	WebhookEventUserRoleChanged                   func(ctx context.Context, eventData WebhookEventUserRoleChangedPayload) error
	WebhookEventUserDeleted                       func(ctx context.Context, eventData WebhookEventUserDeletedPayload) error
	WebhookEventUserSuspendedChanged              func(ctx context.Context, eventData WebhookEventUserSuspendedChangedPayload) error
	WebhookEventUserTagsChanged                   func(ctx context.Context, eventData WebhookEventUserTagsChangedPayload) error
	WebhookEventUserTimeZoneChanged               func(ctx context.Context, eventData WebhookEventUserTimeZoneChangedPayload) error
	WebhookEventUserDefault                       func(ctx context.Context, webhookBody []byte) error

	WebhookEventCommunityPostCreated             func(ctx context.Context, eventData WebhookEventCommunityPostCreatedPayload) error
	WebhookEventCommunityPostChanged             func(ctx context.Context, eventData WebhookEventCommunityPostChangedPayload) error
	WebhookEventCommunityPostPublished           func(ctx context.Context, eventData WebhookEventCommunityPostPublishedPayload) error
	WebhookEventCommunityPostUnpublished         func(ctx context.Context, eventData WebhookEventCommunityPostUnpublishedPayload) error
	WebhookEventCommunityPostSubscriptionCreated func(ctx context.Context, eventData WebhookEventCommunityPostSubscriptionCreatedPayload) error
	WebhookEventCommunityPostVoteCreated         func(ctx context.Context, eventData WebhookEventCommunityPostVoteCreatedPayload) error
	WebhookEventCommunityPostVoteChanged         func(ctx context.Context, eventData WebhookEventCommunityPostVoteChangedPayload) error
	WebhookEventCommunityPostVoteRemoved         func(ctx context.Context, eventData WebhookEventCommunityPostVoteRemovedPayload) error
	WebhookEventCommunityPostCommentCreated      func(ctx context.Context, eventData WebhookEventCommunityPostCommentCreatedPayload) error
	WebhookEventCommunityPostCommentChanged      func(ctx context.Context, eventData WebhookEventCommunityPostCommentChangedPayload) error
	WebhookEventCommunityPostCommentPublished    func(ctx context.Context, eventData WebhookEventCommunityPostCommentPublishedPayload) error
	WebhookEventCommunityPostCommentUnpublished  func(ctx context.Context, eventData WebhookEventCommunityPostCommentUnpublishedPayload) error
	WebhookEventCommunityPostCommentVoteCreated  func(ctx context.Context, eventData WebhookEventCommunityPostCommentVoteCreatedPayload) error
	WebhookEventCommunityPostCommentVoteChanged  func(ctx context.Context, eventData WebhookEventCommunityPostCommentVoteChangedPayload) error
	WebhookEventCommunityPostDefault             func(ctx context.Context, webhookBody []byte) error

	WebhookEventAgentStateChanged        func(ctx context.Context, eventData WebhookEventAgentStateChangedPayload) error
	WebhookEventAgentWorkItemAdded       func(ctx context.Context, eventData WebhookEventAgentWorkItemAddedPayload) error
	WebhookEventAgentWorkItemRemoved     func(ctx context.Context, eventData WebhookEventAgentWorkItemRemovedPayload) error
	WebhookEventAgentMaxCapacityChanged  func(ctx context.Context, eventData WebhookEventAgentMaxCapacityChangedPayload) error
	WebhookEventAgentUnifiedStateChanged func(ctx context.Context, eventData WebhookEventAgentUnifiedStateChangedPayload) error
	WebhookEventAgentChannelCreated      func(ctx context.Context, eventData WebhookEventAgentChannelCreatedPayload) error
	WebhookEventAgentChannelDeleted      func(ctx context.Context, eventData WebhookEventAgentChannelDeletedPayload) error
	WebhookEventAgentDefault             func(ctx context.Context, webhookBody []byte) error

	WebhookEventOmnichannelRoutingConfigFeatureChanged func(ctx context.Context, eventData WebhookEventOmnichannelRoutingConfigFeatureChangedPayload) error
}

// Base webhookevent for Event Based webhooks.
type WebhookEvent struct {
	Type                WebhookEventType `json:"type"`
	AccountID           AccountID        `json:"account_id"`
	ID                  WebhookEventID   `json:"id"`
	Time                time.Time        `json:"time"`
	ZendeskEventVersion string           `json:"zendesk_event_version"`
	Subject             string           `json:"subject"`
	Event               any              `json:"event"`
	Detail              any              `json:"detail"`
}

//gocyclo:ignore
func (s *WebhookService) HandleWebhookEvent(
	eventHandlers WebhookEventHandlers,
	webhookSigningSecret string,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webhookBody, err := readWebhookBody(r)
		if err != nil {
			respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

			return
		}

		baseTarget := WebhookEvent{}
		if err := json.Unmarshal(webhookBody, &baseTarget); err != nil {
			respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

			return
		}

		if baseTarget.Type == "" {
			respondToWebhookRequest(w, http.StatusBadRequest, "Bad Request")

			return
		}

		if !s.verifyZendeskWebhookSignatureIsValid(r, webhookBody, webhookSigningSecret) {
			respondToWebhookRequest(w, http.StatusBadRequest, "Bad Request")

			return
		}

		switch baseTarget.Type {
		// Article Events
		case WebhookEventArticlePublished:
			if eventHandlers.WebhookEventArticlePublished != nil {
				target := WebhookEventArticlePublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventArticlePublished(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventArticleSubscriptionCreated:
			if eventHandlers.WebhookEventArticleSubscriptionCreated != nil {
				target := WebhookEventArticleSubscriptionCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventArticleSubscriptionCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventArticleUnpublished:
			if eventHandlers.WebhookEventArticleUnpublished != nil {
				target := WebhookEventArticleUnpublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventArticleUnpublished(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventArticleVoteCreated:
			if eventHandlers.WebhookEventArticleVoteCreated != nil {
				target := WebhookEventArticleVoteCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventArticleVoteCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventArticleVoteChanged:
			if eventHandlers.WebhookEventArticleVoteChanged != nil {
				target := WebhookEventArticleVoteChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventArticleVoteChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventArticleVoteRemoved:
			if eventHandlers.WebhookEventArticleVoteRemoved != nil {
				target := WebhookEventArticleVoteRemovedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventArticleVoteRemoved(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventArticleCommentCreated:
			if eventHandlers.WebhookEventArticleCommentCreated != nil {
				target := WebhookEventArticleCommentCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventArticleCommentCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventArticleCommentChanged:
			if eventHandlers.WebhookEventArticleCommentChanged != nil {
				target := WebhookEventArticleCommentChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventArticleCommentChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventArticleCommentPublished:
			if eventHandlers.WebhookEventArticleCommentPublished != nil {
				target := WebhookEventArticleCommentPublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventArticleCommentPublished(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventArticleCommentUnpublished:
			if eventHandlers.WebhookEventArticleCommentUnpublished != nil {
				target := WebhookEventArticleCommentUnpublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventArticleCommentUnpublished(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}

		// Organization Events
		case WebhookEventOrganizationCreated:
			if eventHandlers.WebhookEventOrganizationCreated != nil {
				target := WebhookEventOrganizationCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventOrganizationCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventOrganizationCustomFieldChanged:
			if eventHandlers.WebhookEventOrganizationCustomFieldChanged != nil {
				target := WebhookEventOrganizationCustomFieldChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventOrganizationCustomFieldChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventOrganizationDeleted:
			if eventHandlers.WebhookEventOrganizationDeleted != nil {
				target := WebhookEventOrganizationDeletedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventOrganizationDeleted(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventOrganizationExternalIDChanged:
			if eventHandlers.WebhookEventOrganizationExternalIDChanged != nil {
				target := WebhookEventOrganizationExternalIDChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventOrganizationExternalIDChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventOrganizationNameChanged:
			if eventHandlers.WebhookEventOrganizationNameChanged != nil {
				target := WebhookEventOrganizationNameChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventOrganizationNameChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventOrganizationTagsChanged:
			if eventHandlers.WebhookEventOrganizationTagsChanged != nil {
				target := WebhookEventOrganizationTagsChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventOrganizationTagsChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}

		// User Events
		case WebhookEventUserAliasChanged:
			if eventHandlers.WebhookEventUserAliasChanged != nil {
				target := WebhookEventUserAliasChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserAliasChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserCreated:
			if eventHandlers.WebhookEventUserCreated != nil {
				target := WebhookEventUserCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserCustomFieldChanged:
			if eventHandlers.WebhookEventUserCustomFieldChanged != nil {
				target := WebhookEventUserCustomFieldChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserCustomFieldChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserCustomRoleChanged:
			if eventHandlers.WebhookEventUserCustomRoleChanged != nil {
				target := WebhookEventUserCustomRoleChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserCustomRoleChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserDefaultGroupChanged:
			if eventHandlers.WebhookEventUserDefaultGroupChanged != nil {
				target := WebhookEventUserDefaultGroupChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserDefaultGroupChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserDetailsChanged:
			if eventHandlers.WebhookEventUserDetailsChanged != nil {
				target := WebhookEventUserDetailsChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserDetailsChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserExternalIDChanged:
			if eventHandlers.WebhookEventUserExternalIDChanged != nil {
				target := WebhookEventUserExternalIDChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserExternalIDChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserGroupMembershipCreated:
			if eventHandlers.WebhookEventUserGroupMembershipCreated != nil {
				target := WebhookEventUserGroupMembershipCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserGroupMembershipCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserGroupMembershipDeleted:
			if eventHandlers.WebhookEventUserGroupMembershipDeleted != nil {
				target := WebhookEventUserGroupMembershipDeletedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserGroupMembershipDeleted(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserIdentityChanged:
			if eventHandlers.WebhookEventUserIdentityChanged != nil {
				target := WebhookEventUserIdentityChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserIdentityChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserIdentityCreated:
			if eventHandlers.WebhookEventUserIdentityCreated != nil {
				target := WebhookEventUserIdentityCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserIdentityCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserIdentityDeleted:
			if eventHandlers.WebhookEventUserIdentityDeleted != nil {
				target := WebhookEventUserIdentityDeletedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserIdentityDeleted(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserActiveChanged:
			if eventHandlers.WebhookEventUserActiveChanged != nil {
				target := WebhookEventUserActiveChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserActiveChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserLastLoginChanged:
			if eventHandlers.WebhookEventUserLastLoginChanged != nil {
				target := WebhookEventUserLastLoginChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserLastLoginChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserMerged:
			if eventHandlers.WebhookEventUserMerged != nil {
				target := WebhookEventUserMergedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserMerged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserNameChanged:
			if eventHandlers.WebhookEventUserNameChanged != nil {
				target := WebhookEventUserNameChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserNameChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserNotesChanged:
			if eventHandlers.WebhookEventUserNotesChanged != nil {
				target := WebhookEventUserNotesChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserNotesChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserOnlyPrivateCommentsChanged:
			if eventHandlers.WebhookEventUserOnlyPrivateCommentsChanged != nil {
				target := WebhookEventUserOnlyPrivateCommentsChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserOnlyPrivateCommentsChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserOrganizationMembershipCreated:
			if eventHandlers.WebhookEventUserOrganizationMembershipCreated != nil {
				target := WebhookEventUserOrganizationMembershipCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserOrganizationMembershipCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserOrganizationMembershipDeleted:
			if eventHandlers.WebhookEventUserOrganizationMembershipDeleted != nil {
				target := WebhookEventUserOrganizationMembershipDeletedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserOrganizationMembershipDeleted(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserPasswordChanged:
			if eventHandlers.WebhookEventUserPasswordChanged != nil {
				target := WebhookEventUserPasswordChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserPasswordChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserPhotoChanged:
			if eventHandlers.WebhookEventUserPhotoChanged != nil {
				target := WebhookEventUserPhotoChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserPhotoChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserRoleChanged:
			if eventHandlers.WebhookEventUserRoleChanged != nil {
				target := WebhookEventUserRoleChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserRoleChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserDeleted:
			if eventHandlers.WebhookEventUserDeleted != nil {
				target := WebhookEventUserDeletedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserDeleted(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}

		case WebhookEventUserSuspendedChanged:
			if eventHandlers.WebhookEventUserSuspendedChanged != nil {
				target := WebhookEventUserSuspendedChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserSuspendedChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserTagsChanged:
			if eventHandlers.WebhookEventUserTagsChanged != nil {
				target := WebhookEventUserTagsChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserTagsChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventUserTimeZoneChanged:
			if eventHandlers.WebhookEventUserTimeZoneChanged != nil {
				target := WebhookEventUserTimeZoneChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventUserTimeZoneChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}

		// Community Post Events
		case WebhookEventCommunityPostCreated:
			if eventHandlers.WebhookEventCommunityPostCreated != nil {
				target := WebhookEventCommunityPostCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostChanged:
			if eventHandlers.WebhookEventCommunityPostChanged != nil {
				target := WebhookEventCommunityPostChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostPublished:
			if eventHandlers.WebhookEventCommunityPostPublished != nil {
				target := WebhookEventCommunityPostPublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostPublished(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostUnpublished:
			if eventHandlers.WebhookEventCommunityPostUnpublished != nil {
				target := WebhookEventCommunityPostUnpublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostUnpublished(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostSubscriptionCreated:
			if eventHandlers.WebhookEventCommunityPostSubscriptionCreated != nil {
				target := WebhookEventCommunityPostSubscriptionCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostSubscriptionCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostVoteCreated:
			if eventHandlers.WebhookEventCommunityPostVoteCreated != nil {
				target := WebhookEventCommunityPostVoteCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostVoteCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostVoteChanged:
			if eventHandlers.WebhookEventCommunityPostVoteChanged != nil {
				target := WebhookEventCommunityPostVoteChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostVoteChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostVoteRemoved:
			if eventHandlers.WebhookEventCommunityPostVoteRemoved != nil {
				target := WebhookEventCommunityPostVoteRemovedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostVoteRemoved(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentCreated:
			if eventHandlers.WebhookEventCommunityPostCommentCreated != nil {
				target := WebhookEventCommunityPostCommentCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostCommentCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentChanged:
			if eventHandlers.WebhookEventCommunityPostCommentChanged != nil {
				target := WebhookEventCommunityPostCommentChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostCommentChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentPublished:
			if eventHandlers.WebhookEventCommunityPostCommentPublished != nil {
				target := WebhookEventCommunityPostCommentPublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostCommentPublished(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentUnpublished:
			if eventHandlers.WebhookEventCommunityPostCommentUnpublished != nil {
				target := WebhookEventCommunityPostCommentUnpublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostCommentUnpublished(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentVoteCreated:
			if eventHandlers.WebhookEventCommunityPostCommentVoteCreated != nil {
				target := WebhookEventCommunityPostCommentVoteCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostCommentVoteCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentVoteChanged:
			if eventHandlers.WebhookEventCommunityPostCommentVoteChanged != nil {
				target := WebhookEventCommunityPostCommentVoteChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventCommunityPostCommentVoteChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}

		// Agent State Events
		case WebhookEventAgentStateChanged:
			if eventHandlers.WebhookEventAgentStateChanged != nil {
				target := WebhookEventAgentStateChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventAgentStateChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventAgentWorkItemAdded:
			if eventHandlers.WebhookEventAgentWorkItemAdded != nil {
				target := WebhookEventAgentWorkItemAddedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventAgentWorkItemAdded(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventAgentWorkItemRemoved:
			if eventHandlers.WebhookEventAgentWorkItemRemoved != nil {
				target := WebhookEventAgentWorkItemRemovedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventAgentWorkItemRemoved(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventAgentMaxCapacityChanged:
			if eventHandlers.WebhookEventAgentMaxCapacityChanged != nil {
				target := WebhookEventAgentMaxCapacityChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventAgentMaxCapacityChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventAgentUnifiedStateChanged:
			if eventHandlers.WebhookEventAgentUnifiedStateChanged != nil {
				target := WebhookEventAgentUnifiedStateChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventAgentUnifiedStateChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventAgentChannelCreated:
			if eventHandlers.WebhookEventAgentChannelCreated != nil {
				target := WebhookEventAgentChannelCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventAgentChannelCreated(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		case WebhookEventAgentChannelDeleted:
			if eventHandlers.WebhookEventAgentChannelDeleted != nil {
				target := WebhookEventAgentChannelDeletedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventAgentChannelDeleted(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"There was a server error processing the request",
					)

					return
				}
			}

		// Omnichannel Routing Configuration Events
		case WebhookEventOmnichannelRoutingConfigFeatureChanged:
			if eventHandlers.WebhookEventOmnichannelRoutingConfigFeatureChanged != nil {
				target := WebhookEventOmnichannelRoutingConfigFeatureChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := eventHandlers.WebhookEventOmnichannelRoutingConfigFeatureChanged(r.Context(), target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						"Server failed to process Webhook Request correctly",
					)

					return
				}
			}
		default:
			respondToWebhookRequest(
				w,
				http.StatusBadRequest,
				"Unknown webhook event type",
			)

			return
		}

		respondToWebhookRequest(w, http.StatusOK, "Success")
	},
	)
}

func (s *WebhookService) HandleWebhookTrigger(
	handler func(ctx context.Context, webhookBody []byte) error,
	webhookSigningSecret string,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webhookBody, err := readWebhookBody(r)
		if err != nil {
			respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

			return
		}

		if !s.verifyZendeskWebhookSignatureIsValid(r, webhookBody, webhookSigningSecret) {
			respondToWebhookRequest(w, http.StatusBadRequest, "Bad Request")

			return
		}

		if handler != nil {
			if err := handler(r.Context(), webhookBody); err != nil {
				respondToWebhookRequest(w, http.StatusInternalServerError, "Server failed to process Webhook Request correctly")

				return
			}
		}

		respondToWebhookRequest(w, http.StatusOK, "Success")
	})
}

type WebhookEventOrganization[EventData WebhookOrganizationEventData] struct {
	Type                WebhookEventType               `json:"type"`
	AccountID           AccountID                      `json:"account_id"`
	ID                  WebhookEventID                 `json:"id"`
	Time                time.Time                      `json:"time"`
	ZendeskEventVersion string                         `json:"zendesk_event_version"`
	Subject             string                         `json:"subject"`
	Event               EventData                      `json:"event"`
	Detail              WebhookEventOrganizationDetail `json:"detail"`
}

type WebhookEventOrganizationDetail struct {
	CreatedAt      time.Time `json:"created_at"`
	ExternalID     string    `json:"external_id"`
	GroupID        string    `json:"group_id"`
	Email          string    `json:"email"`
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	SharedComments bool      `json:"shared_comments"`
	SharedTickets  bool      `json:"shared_tickets"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/webhook-event-types/
type WebhookEventArticle[EventData WebhookArticleEventData] struct {
	Type                WebhookEventType          `json:"type"`
	AccountID           AccountID                 `json:"account_id"`
	ID                  WebhookEventID            `json:"id"`
	Time                time.Time                 `json:"time"`
	ZendeskEventVersion string                    `json:"zendesk_event_version"`
	Subject             string                    `json:"subject"`
	Event               EventData                 `json:"event"`
	Detail              WebhookEventArticleDetail `json:"detail"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/article-events/
type WebhookEventArticleDetail struct {
	BrandID BrandID   `json:"brand_id"`
	ID      ArticleID `json:"id"`
}

type WebhookEventCommunityPost[EventData WebhookCommunityPostEventData] struct {
	Type                WebhookEventType                `json:"type"`
	AccountID           AccountID                       `json:"account_id"`
	ID                  WebhookEventID                  `json:"id"`
	Time                time.Time                       `json:"time"`
	ZendeskEventVersion string                          `json:"zendesk_event_version"`
	Subject             string                          `json:"subject"`
	Event               EventData                       `json:"event"`
	Detail              WebhookEventCommunityPostDetail `json:"detail"`
}

type WebhookEventCommunityPostDetail struct {
	BrandID string `json:"brand_id"`
	ID      string `json:"id"`
	PostID  string `json:"post_id"`
}

type WebhookEventAgentState[EventData WebhookAgentStateEventData] struct {
	Type                WebhookEventType             `json:"type"`
	AccountID           AccountID                    `json:"account_id"`
	ID                  WebhookEventID               `json:"id"`
	Time                time.Time                    `json:"time"`
	ZendeskEventVersion string                       `json:"zendesk_event_version"`
	Subject             string                       `json:"subject"`
	Event               EventData                    `json:"event"`
	Detail              WebhookEventAgentStateDetail `json:"detail"`
}

type WebhookEventAgentStateDetail struct {
	AccountID string `json:"account_id"`
	AgentID   string `json:"agent_id"`
	Version   string `json:"version"`
}

type WebhookEventOmnichannelRoutingConfig[EventData WebhookOmnichannelRoutingConfigData] struct {
	Type                WebhookEventType                           `json:"type"`
	AccountID           AccountID                                  `json:"account_id"`
	ID                  WebhookEventID                             `json:"id"`
	Time                time.Time                                  `json:"time"`
	ZendeskEventVersion string                                     `json:"zendesk_event_version"`
	Subject             string                                     `json:"subject"`
	Event               EventData                                  `json:"event"`
	Detail              WebhookEventOmnichannelRoutingConfigDetail `json:"detail"`
}

type WebhookEventOmnichannelRoutingConfigDetail struct {
	AccountID string `json:"account_id"`
}

type WebhookArticleEventData interface {
	EventTypeArticleAuthorChangedEvent | any
}

type WebhookOrganizationEventData interface {
	WebhookEventDataEmpty |
		WebhookEventDataCustomFieldUpdate |
		WebhookEventDataSimpleStringUpdate |
		WebhookEventDataTagsChanged
}

type WebhookAgentStateEventData interface {
	any
}

type WebhookOmnichannelRoutingConfigData interface {
	WebhookEventDataSimpleBoolUpdateValue
}

type WebhookCommunityPostEventData interface {
	any
}

type EventTypeArticleAuthorChangedEvent struct{}

type WebhookEventUser[EventData WebhookUserEventData] struct {
	Type                WebhookEventType       `json:"type"`
	AccountID           AccountID              `json:"account_id"`
	ID                  WebhookEventID         `json:"id"`
	Time                time.Time              `json:"time"`
	ZendeskEventVersion string                 `json:"zendesk_event_version"`
	Subject             string                 `json:"subject"`
	Event               EventData              `json:"event"`
	Detail              WebhookEventUserDetail `json:"detail"`
}

type WebhookEventUserDetail struct {
	CreatedAt      time.Time `json:"created_at"`
	Email          string    `json:"email"`
	ExternalID     string    `json:"external_id"`
	DefaultGroupID string    `json:"default_group_id"`
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Role           string    `json:"role"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type WebhookUserEventData interface {
	WebhookEventDataSimpleStringUpdate |
		WebhookEventDataEmpty |
		WebhookEventDataCustomFieldUpdate |
		WebhookEventDataUserGroupMembershipChanged |
		WebhookEventDataUserIdentityChanged |
		WebhookEventDataUserIdentity |
		WebhookEventDataSimpleBoolUpdate |
		WebhookEventDataUserMerged |
		WebhookEventDataUserOrganizationMembershipChanged |
		WebhookEventDataTagsChanged
}

type WebhookEventUserActiveStatusChangedPayload struct {
	Current  bool `json:"current"`
	Previous bool `json:"previous"`
}

// https://developer.zendesk.com/api-reference/webhooks/event-types/community-events/#detail-object-properties
type WebhookEventDetailCommunityPost struct {
	BrandID BrandID `json:"brand_id"`
	PostID  PostID  `json:"post_id"`
	TopicID TopicID `json:"topic_id"`
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

// https://developer.zendesk.com/documentation/webhooks/verifying/
func (s WebhookService) verifyZendeskWebhookSignatureIsValid(
	r *http.Request,
	bodyBytes []byte,
	webhookSigningSecret string,
) bool {
	expectedZendeskSignature := r.Header.Get(WebhookHeaderSignature)
	zendeskSignatureTimestamp := r.Header.Get(WebhookHeaderSignatureTimestamp)

	if expectedZendeskSignature == "" || zendeskSignatureTimestamp == "" {
		return false
	}

	actualZendeskSignature := buildZendeskSignature(zendeskSignatureTimestamp, bodyBytes, webhookSigningSecret)

	return expectedZendeskSignature == actualZendeskSignature
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

type WebhookEventDataUserIdentitySchema struct {
	ID      string `json:"id"`
	Primary bool   `json:"primary"`
	Type    string `json:"type"`
	Value   string `json:"value"`
}

type WebhookEventDataSimpleStringUpdate struct {
	Current  string `json:"current"`
	Previous string `json:"previous"`
}

type WebhookEventDataSimpleBoolUpdate struct {
	Current  bool `json:"current"`
	Previous bool `json:"previous"`
}

type WebhookEventDataSimpleBoolUpdateValue struct {
	CurrentValue  bool `json:"current_value"`
	PreviousValue bool `json:"previous_value"`
}

type WebhookEventDataEmpty struct{}

type WebhookEventDataCustomFieldUpdate struct {
	Current  any `json:"current"`
	Previous any `json:"previous"`
	Field    struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"field"`
}

type WebhookEventDataUserGroupMembershipChanged struct {
	Group WebhookEventDataIDField `json:"group"`
}

type WebhookEventDataIDField struct {
	ID string `json:"id"`
}

type WebhookEventDataUserIdentityChanged struct {
	Current  WebhookEventDataUserIdentitySchema `json:"current"`
	Previous WebhookEventDataUserIdentitySchema `json:"previous"`
}

type WebhookEventDataUserIdentity struct {
	Identity WebhookEventDataUserIdentitySchema `json:"current"`
}

type WebhookEventDataUserMerged struct {
	User struct {
		ID string `json:"id"`
	} `json:"user"`
}

type WebhookEventDataUserOrganizationMembershipChanged struct {
	Organization struct {
		ID string `json:"id"`
	} `json:"organization"`
}

type WebhookEventDataTagsChanged struct {
	Added struct {
		Tags []string `json:"tags"`
	} `json:"added"`
	Removed struct {
		Tags []string `json:"tags"`
	} `json:"removed"`
}

type (
	WebhookEventUserAliasChangedPayload                  WebhookEventUser[WebhookEventDataSimpleStringUpdate]
	WebhookEventUserCreatedPayload                       WebhookEventUser[WebhookEventDataEmpty]
	WebhookEventUserCustomFieldChangedPayload            WebhookEventUser[WebhookEventDataCustomFieldUpdate]
	WebhookEventUserCustomRoleChangedPayload             WebhookEventUser[WebhookEventDataSimpleStringUpdate]
	WebhookEventUserDefaultGroupChangedPayload           WebhookEventUser[WebhookEventDataSimpleStringUpdate]
	WebhookEventUserDetailsChangedPayload                WebhookEventUser[WebhookEventDataSimpleStringUpdate]
	WebhookEventUserExternalIDChangedPayload             WebhookEventUser[WebhookEventDataSimpleStringUpdate]
	WebhookEventUserGroupMembershipCreatedPayload        WebhookEventUser[WebhookEventDataUserGroupMembershipChanged]
	WebhookEventUserGroupMembershipDeletedPayload        WebhookEventUser[WebhookEventDataUserGroupMembershipChanged]
	WebhookEventUserIdentityChangedPayload               WebhookEventUser[WebhookEventDataUserIdentityChanged]
	WebhookEventUserIdentityCreatedPayload               WebhookEventUser[WebhookEventDataUserIdentity]
	WebhookEventUserIdentityDeletedPayload               WebhookEventUser[WebhookEventDataUserIdentity]
	WebhookEventUserActiveChangedPayload                 WebhookEventUser[WebhookEventDataSimpleBoolUpdate]
	WebhookEventUserLastLoginChangedPayload              WebhookEventUser[WebhookEventDataSimpleStringUpdate]
	WebhookEventUserMergedPayload                        WebhookEventUser[WebhookEventDataUserMerged]
	WebhookEventUserNameChangedPayload                   WebhookEventUser[WebhookEventDataSimpleStringUpdate]
	WebhookEventUserNotesChangedPayload                  WebhookEventUser[WebhookEventDataSimpleStringUpdate]
	WebhookEventUserOnlyPrivateCommentsChangedPayload    WebhookEventUser[WebhookEventDataSimpleBoolUpdate]
	WebhookEventUserOrganizationMembershipCreatedPayload WebhookEventUser[WebhookEventDataUserOrganizationMembershipChanged]
	WebhookEventUserOrganizationMembershipDeletedPayload WebhookEventUser[WebhookEventDataUserOrganizationMembershipChanged]
	WebhookEventUserPasswordChangedPayload               WebhookEventUser[WebhookEventDataEmpty]
	WebhookEventUserPhotoChangedPayload                  WebhookEventUser[WebhookEventDataSimpleStringUpdate]
	WebhookEventUserRoleChangedPayload                   WebhookEventUser[WebhookEventDataSimpleStringUpdate]
	WebhookEventUserDeletedPayload                       WebhookEventUser[WebhookEventDataEmpty]
	WebhookEventUserSuspendedChangedPayload              WebhookEventUser[WebhookEventDataSimpleBoolUpdate]
	WebhookEventUserTagsChangedPayload                   WebhookEventUser[WebhookEventDataTagsChanged]
	WebhookEventUserTimeZoneChangedPayload               WebhookEventUser[WebhookEventDataSimpleStringUpdate]
)

type (
	WebhookEventOrganizationCreatedPayload            WebhookEventOrganization[WebhookEventDataEmpty]
	WebhookEventOrganizationCustomFieldChangedPayload WebhookEventOrganization[WebhookEventDataCustomFieldUpdate]
	WebhookEventOrganizationDeletedPayload            WebhookEventOrganization[WebhookEventDataEmpty]
	WebhookEventOrganizationExternalIDChangedPayload  WebhookEventOrganization[WebhookEventDataSimpleStringUpdate]
	WebhookEventOrganizationNameChangedPayload        WebhookEventOrganization[WebhookEventDataSimpleStringUpdate]
	WebhookEventOrganizationTagsChangedPayload        WebhookEventOrganization[WebhookEventDataTagsChanged]
)

type (
	WebhookEventArticlePublishedPayload           WebhookEventArticle[any]
	WebhookEventArticleSubscriptionCreatedPayload WebhookEventArticle[any]
	WebhookEventArticleUnpublishedPayload         WebhookEventArticle[any]
	WebhookEventArticleVoteCreatedPayload         WebhookEventArticle[any]
	WebhookEventArticleVoteChangedPayload         WebhookEventArticle[any]
	WebhookEventArticleVoteRemovedPayload         WebhookEventArticle[any]
	WebhookEventArticleCommentCreatedPayload      WebhookEventArticle[any]
	WebhookEventArticleCommentChangedPayload      WebhookEventArticle[any]
	WebhookEventArticleCommentPublishedPayload    WebhookEventArticle[any]
	WebhookEventArticleCommentUnpublishedPayload  WebhookEventArticle[any]
	WebhookEventArticleAuthorChangedPayload       WebhookEventArticle[any]
)

type (
	WebhookEventCommunityPostCreatedPayload             WebhookEventCommunityPost[any]
	WebhookEventCommunityPostChangedPayload             WebhookEventCommunityPost[any]
	WebhookEventCommunityPostPublishedPayload           WebhookEventCommunityPost[any]
	WebhookEventCommunityPostUnpublishedPayload         WebhookEventCommunityPost[any]
	WebhookEventCommunityPostSubscriptionCreatedPayload WebhookEventCommunityPost[any]
	WebhookEventCommunityPostVoteCreatedPayload         WebhookEventCommunityPost[any]
	WebhookEventCommunityPostVoteChangedPayload         WebhookEventCommunityPost[any]
	WebhookEventCommunityPostVoteRemovedPayload         WebhookEventCommunityPost[any]
	WebhookEventCommunityPostCommentCreatedPayload      WebhookEventCommunityPost[any]
	WebhookEventCommunityPostCommentChangedPayload      WebhookEventCommunityPost[any]
	WebhookEventCommunityPostCommentPublishedPayload    WebhookEventCommunityPost[any]
	WebhookEventCommunityPostCommentUnpublishedPayload  WebhookEventCommunityPost[any]
	WebhookEventCommunityPostCommentVoteCreatedPayload  WebhookEventCommunityPost[any]
	WebhookEventCommunityPostCommentVoteChangedPayload  WebhookEventCommunityPost[any]
)

type (
	WebhookEventAgentStateChangedPayload        WebhookEventAgentState[any]
	WebhookEventAgentWorkItemAddedPayload       WebhookEventAgentState[any]
	WebhookEventAgentWorkItemRemovedPayload     WebhookEventAgentState[any]
	WebhookEventAgentMaxCapacityChangedPayload  WebhookEventAgentState[any]
	WebhookEventAgentUnifiedStateChangedPayload WebhookEventAgentState[any]
	WebhookEventAgentChannelCreatedPayload      WebhookEventAgentState[any]
	WebhookEventAgentChannelDeletedPayload      WebhookEventAgentState[any]
)

type (
	WebhookEventOmnichannelRoutingConfigFeatureChangedPayload WebhookEventOmnichannelRoutingConfig[WebhookEventDataSimpleBoolUpdateValue]
)
