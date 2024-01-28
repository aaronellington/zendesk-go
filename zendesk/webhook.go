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
	client            *client
	eventHandlers     *WebhookEventHandlers
	eventHandlerCache map[WebhookEventType]int
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
	WebhookEventUserPasswordChanged               WebhookEventType = "zen:event-type:user.password_changed"
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

type WebhookEventHandlers struct {
	WebhookEventArticlePublished           func(eventData WebhookEventArticlePublishedPayload) error           ``
	WebhookEventArticleSubscriptionCreated func(eventData WebhookEventArticleSubscriptionCreatedPayload) error ``
	WebhookEventArticleUnpublished         func(eventData WebhookEventArticleUnpublishedPayload) error         ``
	WebhookEventArticleVoteCreated         func(eventData WebhookEventArticleVoteCreatedPayload) error         ``
	WebhookEventArticleVoteChanged         func(eventData WebhookEventArticleVoteChangedPayload) error         ``
	WebhookEventArticleVoteRemoved         func(eventData WebhookEventArticleVoteRemovedPayload) error         ``
	WebhookEventArticleCommentCreated      func(eventData WebhookEventArticleCommentCreatedPayload) error      ``
	WebhookEventArticleCommentChanged      func(eventData WebhookEventArticleCommentChangedPayload) error      ``
	WebhookEventArticleCommentPublished    func(eventData WebhookEventArticleCommentPublishedPayload) error    ``
	WebhookEventArticleCommentUnpublished  func(eventData WebhookEventArticleCommentUnpublishedPayload) error  ``

	WebhookEventOrganizationCreated            func(eventData WebhookEventOrganizationCreatedPayload) error            ``
	WebhookEventOrganizationCustomFieldChanged func(eventData WebhookEventOrganizationCustomFieldChangedPayload) error ``
	WebhookEventOrganizationDeleted            func(eventData WebhookEventOrganizationDeletedPayload) error            ``
	WebhookEventOrganizationExternalIDChanged  func(eventData WebhookEventOrganizationExternalIDChangedPayload) error  ``
	WebhookEventOrganizationNameChanged        func(eventData WebhookEventOrganizationNameChangedPayload) error        ``
	WebhookEventOrganizationTagsChanged        func(eventData WebhookEventOrganizationTagsChangedPayload) error        ``

	WebhookEventUserCreated                       func(eventData WebhookEventUserCreatedPayload) error                       `zenevent:"zen:event-type:user.created"`
	WebhookEventUserCustomFieldChanged            func(eventData WebhookEventUserCustomFieldChangedPayload) error            `zenevent:"zen:event-type:user.custom_field_changed"`
	WebhookEventUserCustomRoleChanged             func(eventData WebhookEventUserCustomRoleChangedPayload) error             ``
	WebhookEventUserDefaultGroupChanged           func(eventData WebhookEventUserDefaultGroupChangedPayload) error           ``
	WebhookEventUserExternalIDChanged             func(eventData WebhookEventUserExternalIDChangedPayload) error             ``
	WebhookEventUserGroupMembershipCreated        func(eventData WebhookEventUserGroupMembershipCreatedPayload) error        ``
	WebhookEventUserGroupMembershipDeleted        func(eventData WebhookEventUserGroupMembershipDeletedPayload) error        ``
	WebhookEventUserIdentityChanged               func(eventData WebhookEventUserIdentityChangedPayload) error               ``
	WebhookEventUserIdentityCreated               func(eventData WebhookEventUserIdentityCreatedPayload) error               ``
	WebhookEventUserIdentityDeleted               func(eventData WebhookEventUserIdentityDeletedPayload) error               ``
	WebhookEventUserActiveChanged                 func(eventData WebhookEventUserActiveChangedPayload) error                 ``
	WebhookEventUserLastLoginChanged              func(eventData WebhookEventUserLastLoginChangedPayload) error              ``
	WebhookEventUserMerged                        func(eventData WebhookEventUserMergedPayload) error                        ``
	WebhookEventUserNameChanged                   func(eventData WebhookEventUserNameChangedPayload) error                   ``
	WebhookEventUserNotesChanged                  func(eventData WebhookEventUserNotesChangedPayload) error                  ``
	WebhookEventUserOnlyPrivateCommentsChanged    func(eventData WebhookEventUserOnlyPrivateCommentsChangedPayload) error    ``
	WebhookEventUserOrganizationMembershipCreated func(eventData WebhookEventUserOrganizationMembershipCreatedPayload) error ``
	WebhookEventUserOrganizationMembershipDeleted func(eventData WebhookEventUserOrganizationMembershipDeletedPayload) error ``
	WebhookEventUserPasswordChanged               func(eventData WebhookEventUserPasswordChangedPayload) error               ``
	WebhookEventUserPhotoChanged                  func(eventData WebhookEventUserPhotoChangedPayload) error                  ``
	WebhookEventUserRoleChanged                   func(eventData WebhookEventUserRoleChangedPayload) error                   ``
	WebhookEventUserDeleted                       func(eventData WebhookEventUserDeletedPayload) error                       ``
	WebhookEventUserSuspendedChanged              func(eventData WebhookEventUserSuspendedChangedPayload) error              ``
	WebhookEventUserTagsChanged                   func(eventData WebhookEventUserTagsChangedPayload) error                   ``
	WebhookEventUserTimeZoneChanged               func(eventData WebhookEventUserTimeZoneChangedPayload) error               ``

	WebhookEventCommunityPostCreated             func(eventData WebhookEventCommunityPostCreatedPayload) error             ``
	WebhookEventCommunityPostChanged             func(eventData WebhookEventCommunityPostChangedPayload) error             ``
	WebhookEventCommunityPostPublished           func(eventData WebhookEventCommunityPostPublishedPayload) error           ``
	WebhookEventCommunityPostUnpublished         func(eventData WebhookEventCommunityPostUnpublishedPayload) error         ``
	WebhookEventCommunityPostSubscriptionCreated func(eventData WebhookEventCommunityPostSubscriptionCreatedPayload) error ``
	WebhookEventCommunityPostVoteCreated         func(eventData WebhookEventCommunityPostVoteCreatedPayload) error         ``
	WebhookEventCommunityPostVoteChanged         func(eventData WebhookEventCommunityPostVoteChangedPayload) error         ``
	WebhookEventCommunityPostVoteRemoved         func(eventData WebhookEventCommunityPostVoteRemovedPayload) error         ``
	WebhookEventCommunityPostCommentCreated      func(eventData WebhookEventCommunityPostCommentCreatedPayload) error      ``
	WebhookEventCommunityPostCommentChanged      func(eventData WebhookEventCommunityPostCommentChangedPayload) error      ``
	WebhookEventCommunityPostCommentPublished    func(eventData WebhookEventCommunityPostCommentPublishedPayload) error    ``
	WebhookEventCommunityPostCommentUnpublished  func(eventData WebhookEventCommunityPostCommentUnpublishedPayload) error  ``
	WebhookEventCommunityPostCommentVoteCreated  func(eventData WebhookEventCommunityPostCommentVoteCreatedPayload) error  ``
	WebhookEventCommunityPostCommentVoteChanged  func(eventData WebhookEventCommunityPostCommentVoteChangedPayload) error  ``

	WebhookEventAgentStateChanged        func(eventData WebhookEventAgentStateChangedPayload) error        ``
	WebhookEventAgentWorkItemAdded       func(eventData WebhookEventAgentWorkItemAddedPayload) error       ``
	WebhookEventAgentWorkItemRemoved     func(eventData WebhookEventAgentWorkItemRemovedPayload) error     ``
	WebhookEventAgentMaxCapacityChanged  func(eventData WebhookEventAgentMaxCapacityChangedPayload) error  ``
	WebhookEventAgentUnifiedStateChanged func(eventData WebhookEventAgentUnifiedStateChangedPayload) error ``
	WebhookEventAgentChannelCreated      func(eventData WebhookEventAgentChannelCreatedPayload) error      ``
	WebhookEventAgentChannelDeleted      func(eventData WebhookEventAgentChannelDeletedPayload) error      ``

	WebhookEventOmnichannelRoutingConfigFeatureChanged func(eventData WebhookEventOmnichannelRoutingConfigFeatureChangedPayload) error ``
}

// Base webhookevent for Event Based webhooks
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

func (s *WebhookService) HandleWebhookEvent(webhookSigningSecret string) http.Handler {
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

		if webhookSigningSecret != "" {
			if !s.verifyZendeskWebhookSignatureIsValid(r, webhookBody, webhookSigningSecret) {
				respondToWebhookRequest(w, http.StatusBadRequest, "Bad Request")

				return
			}
		}

		switch baseTarget.Type {
		case WebhookEventArticlePublished:
			if s.eventHandlers.WebhookEventArticlePublished != nil {
				target := WebhookEventArticlePublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventArticlePublished(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventArticleSubscriptionCreated:
			if s.eventHandlers.WebhookEventArticleSubscriptionCreated != nil {
				target := WebhookEventArticleSubscriptionCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventArticleSubscriptionCreated(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventArticleUnpublished:
			if s.eventHandlers.WebhookEventArticleUnpublished != nil {
				target := WebhookEventArticleUnpublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventArticleUnpublished(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventArticleVoteCreated:
			if s.eventHandlers.WebhookEventArticleVoteCreated != nil {
				target := WebhookEventArticleVoteCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventArticleVoteCreated(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventArticleVoteChanged:
			if s.eventHandlers.WebhookEventArticleVoteChanged != nil {
				target := WebhookEventArticleVoteChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventArticleVoteChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventArticleCommentCreated:
			if s.eventHandlers.WebhookEventArticleCommentCreated != nil {
				target := WebhookEventArticleCommentCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventArticleCommentCreated(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventArticleCommentChanged:
			if s.eventHandlers.WebhookEventArticleCommentChanged != nil {
				target := WebhookEventArticleCommentChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventArticleCommentChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventArticleCommentPublished:
			if s.eventHandlers.WebhookEventArticleCommentPublished != nil {
				target := WebhookEventArticleCommentPublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventArticleCommentPublished(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventArticleCommentUnpublished:
			if s.eventHandlers.WebhookEventArticleCommentUnpublished != nil {
				target := WebhookEventArticleCommentUnpublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventArticleCommentUnpublished(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventOrganizationCreated:
			if s.eventHandlers.WebhookEventOrganizationCreated != nil {
				target := WebhookEventOrganizationCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventOrganizationCreated(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventOrganizationCustomFieldChanged:
			if s.eventHandlers.WebhookEventOrganizationCustomFieldChanged != nil {
				target := WebhookEventOrganizationCustomFieldChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventOrganizationCustomFieldChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventOrganizationDeleted:
			if s.eventHandlers.WebhookEventOrganizationDeleted != nil {
				target := WebhookEventOrganizationDeletedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventOrganizationDeleted(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserLastLoginChanged:
			if s.eventHandlers.WebhookEventUserLastLoginChanged != nil {
				target := WebhookEventUserLastLoginChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserLastLoginChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserMerged:
			if s.eventHandlers.WebhookEventUserMerged != nil {
				target := WebhookEventUserMergedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserMerged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserNameChanged:
			if s.eventHandlers.WebhookEventUserNameChanged != nil {
				target := WebhookEventUserNameChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserNameChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserNotesChanged:
			if s.eventHandlers.WebhookEventUserNotesChanged != nil {
				target := WebhookEventUserNotesChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserNotesChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserOnlyPrivateCommentsChanged:
			if s.eventHandlers.WebhookEventUserOnlyPrivateCommentsChanged != nil {
				target := WebhookEventUserOnlyPrivateCommentsChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserOnlyPrivateCommentsChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserOrganizationMembershipCreated:
			if s.eventHandlers.WebhookEventUserOrganizationMembershipCreated != nil {
				target := WebhookEventUserOrganizationMembershipCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserOrganizationMembershipCreated(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserOrganizationMembershipDeleted:
			if s.eventHandlers.WebhookEventUserOrganizationMembershipDeleted != nil {
				target := WebhookEventUserOrganizationMembershipDeletedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserOrganizationMembershipDeleted(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserPasswordChanged:
			if s.eventHandlers.WebhookEventUserPasswordChanged != nil {
				target := WebhookEventUserPasswordChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserPasswordChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserPhotoChanged:
			if s.eventHandlers.WebhookEventUserPhotoChanged != nil {
				target := WebhookEventUserPhotoChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserPhotoChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserRoleChanged:
			if s.eventHandlers.WebhookEventUserRoleChanged != nil {
				target := WebhookEventUserRoleChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserRoleChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserDeleted:
			if s.eventHandlers.WebhookEventUserDeleted != nil {
				target := WebhookEventUserDeletedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserDeleted(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserActiveChanged:
			if s.eventHandlers.WebhookEventUserActiveChanged != nil {
				target := WebhookEventUserActiveChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserActiveChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserSuspendedChanged:
			if s.eventHandlers.WebhookEventUserSuspendedChanged != nil {
				target := WebhookEventUserSuspendedChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserSuspendedChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserTagsChanged:
			if s.eventHandlers.WebhookEventUserTagsChanged != nil {
				target := WebhookEventUserTagsChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserTagsChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventUserTimeZoneChanged:
			if s.eventHandlers.WebhookEventUserTimeZoneChanged != nil {
				target := WebhookEventUserTimeZoneChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventUserTimeZoneChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostCreated:
			if s.eventHandlers.WebhookEventCommunityPostCreated != nil {
				target := WebhookEventCommunityPostCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostCreated(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostChanged:
			if s.eventHandlers.WebhookEventCommunityPostChanged != nil {
				target := WebhookEventCommunityPostChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostPublished:
			if s.eventHandlers.WebhookEventCommunityPostPublished != nil {
				target := WebhookEventCommunityPostPublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostPublished(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostUnpublished:
			if s.eventHandlers.WebhookEventCommunityPostUnpublished != nil {
				target := WebhookEventCommunityPostUnpublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostUnpublished(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostSubscriptionCreated:
			if s.eventHandlers.WebhookEventCommunityPostSubscriptionCreated != nil {
				target := WebhookEventCommunityPostSubscriptionCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostSubscriptionCreated(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostVoteCreated:
			if s.eventHandlers.WebhookEventCommunityPostVoteCreated != nil {
				target := WebhookEventCommunityPostVoteCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostVoteCreated(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostVoteChanged:
			if s.eventHandlers.WebhookEventCommunityPostVoteChanged != nil {
				target := WebhookEventCommunityPostVoteChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostVoteChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostVoteRemoved:
			if s.eventHandlers.WebhookEventCommunityPostVoteRemoved != nil {
				target := WebhookEventCommunityPostVoteRemovedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostVoteRemoved(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentCreated:
			if s.eventHandlers.WebhookEventCommunityPostCommentCreated != nil {
				target := WebhookEventCommunityPostCommentCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostCommentCreated(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentChanged:
			if s.eventHandlers.WebhookEventCommunityPostCommentChanged != nil {
				target := WebhookEventCommunityPostCommentChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostCommentChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentPublished:
			if s.eventHandlers.WebhookEventCommunityPostCommentPublished != nil {
				target := WebhookEventCommunityPostCommentPublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostCommentPublished(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentUnpublished:
			if s.eventHandlers.WebhookEventCommunityPostCommentUnpublished != nil {
				target := WebhookEventCommunityPostCommentUnpublishedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostCommentUnpublished(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentVoteCreated:
			if s.eventHandlers.WebhookEventCommunityPostCommentVoteCreated != nil {
				target := WebhookEventCommunityPostCommentVoteCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostCommentVoteCreated(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventCommunityPostCommentVoteChanged:
			if s.eventHandlers.WebhookEventCommunityPostCommentVoteChanged != nil {
				target := WebhookEventCommunityPostCommentVoteChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventCommunityPostCommentVoteChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventAgentStateChanged:
			if s.eventHandlers.WebhookEventAgentStateChanged != nil {
				target := WebhookEventAgentStateChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventAgentStateChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventAgentWorkItemAdded:
			if s.eventHandlers.WebhookEventAgentWorkItemAdded != nil {
				target := WebhookEventAgentWorkItemAddedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventAgentWorkItemAdded(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventAgentWorkItemRemoved:
			if s.eventHandlers.WebhookEventAgentWorkItemRemoved != nil {
				target := WebhookEventAgentWorkItemRemovedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventAgentWorkItemRemoved(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventAgentMaxCapacityChanged:
			if s.eventHandlers.WebhookEventAgentMaxCapacityChanged != nil {
				target := WebhookEventAgentMaxCapacityChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventAgentMaxCapacityChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventAgentUnifiedStateChanged:
			if s.eventHandlers.WebhookEventAgentUnifiedStateChanged != nil {
				target := WebhookEventAgentUnifiedStateChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventAgentUnifiedStateChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventAgentChannelCreated:
			if s.eventHandlers.WebhookEventAgentChannelCreated != nil {
				target := WebhookEventAgentChannelCreatedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventAgentChannelCreated(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventAgentChannelDeleted:
			if s.eventHandlers.WebhookEventAgentChannelDeleted != nil {
				target := WebhookEventAgentChannelDeletedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventAgentChannelDeleted(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)

					return
				}
			}
		case WebhookEventOmnichannelRoutingConfigFeatureChanged:
			if s.eventHandlers.WebhookEventOmnichannelRoutingConfigFeatureChanged != nil {
				target := WebhookEventOmnichannelRoutingConfigFeatureChangedPayload{}
				if err := json.Unmarshal(webhookBody, &target); err != nil {
					respondToWebhookRequest(w, http.StatusBadRequest, err.Error())

					return
				}

				if err := s.eventHandlers.WebhookEventOmnichannelRoutingConfigFeatureChanged(target); err != nil {
					respondToWebhookRequest(
						w,
						http.StatusInternalServerError,
						err.Error(),
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

func (s *WebhookService) WebhookTriggerHandler(handler func(b []byte), secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

type WebhookEventOrganization[EventData WebhookOrganizationEventData] struct {
	Type                WebhookEventType          `json:"type"`
	AccountID           AccountID                 `json:"account_id"`
	ID                  WebhookEventID            `json:"id"`
	Time                time.Time                 `json:"time"`
	ZendeskEventVersion string                    `json:"zendesk_event_version"`
	Subject             string                    `json:"subject"`
	Event               EventData                 `json:"event"`
	Detail              WebhookEventArticleDetail `json:"detail"`
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

type WebhookEventCommunityPost[EventData WebhookCommunityPostEventData] struct {
	Type                WebhookEventType          `json:"type"`
	AccountID           AccountID                 `json:"account_id"`
	ID                  WebhookEventID            `json:"id"`
	Time                time.Time                 `json:"time"`
	ZendeskEventVersion string                    `json:"zendesk_event_version"`
	Subject             string                    `json:"subject"`
	Event               EventData                 `json:"event"`
	Detail              WebhookEventArticleDetail `json:"detail"`
}

type WebhookEventAgentState[EventData WebhookAgentStateEventData] struct {
	Type                WebhookEventType          `json:"type"`
	AccountID           AccountID                 `json:"account_id"`
	ID                  WebhookEventID            `json:"id"`
	Time                time.Time                 `json:"time"`
	ZendeskEventVersion string                    `json:"zendesk_event_version"`
	Subject             string                    `json:"subject"`
	Event               EventData                 `json:"event"`
	Detail              WebhookEventArticleDetail `json:"detail"`
}

type WebhookEventOmnichannelRoutingConfig[EventData WebhookOmnichannelRoutingConfigData] struct {
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

type WebhookArticleEventData interface {
	EventTypeArticleAuthorChangedEvent | any
}

type WebhookOrganizationEventData interface {
	any
}

type WebhookAgentStateEventData interface {
	any
}

type WebhookOmnichannelRoutingConfigData interface {
	any
}

type WebhookCommunityPostEventData interface {
	any
}

type WebhookEventArticleAuthorChangedPayload WebhookEventArticle[EventTypeArticleAuthorChangedEvent]

type EventTypeArticleAuthorChangedEvent struct {
}

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
	CreatedAt      time.Time    `json:"created_at"`
	Email          string       `json:"email"`
	ExternalID     string       `json:"external_id"`
	DefaultGroupID string       `json:"default_group_id"`
	ID             string       `json:"id"`
	OrganizationID string       `json:"organization_id"`
	Role           CustomRoleID `json:"role"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type WebhookUserEventData interface {
	WebhookEventUserAliasChangedPayload |
		WebhookEventUserActiveStatusChangedPayload |
		WebhookEventUserDetailsChangedPayload |
		any
}

type WebhookEventUserActiveStatusChangedPayload struct {
	Current  bool `json:"current"`
	Previous bool `json:"previous"`
}

type WebhookEventUserDetailsChangedPayload struct {
	Current  string `json:"current"`
	Previous string `json:"previous"`
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
	return expectedZendeskSignature != actualZendeskSignature
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

type WebhookEventUserAliasChangedPayload WebhookEventUser[any]
type WebhookEventUserCreatedPayload WebhookEventUser[any]
type WebhookEventUserCustomFieldChangedPayload WebhookEventUser[any]
type WebhookEventUserCustomRoleChangedPayload WebhookEventUser[any]
type WebhookEventUserDefaultGroupChangedPayload WebhookEventUser[any]
type WebhookEventUserExternalIDChangedPayload WebhookEventUser[any]
type WebhookEventUserGroupMembershipCreatedPayload WebhookEventUser[any]
type WebhookEventUserGroupMembershipDeletedPayload WebhookEventUser[any]
type WebhookEventUserIdentityChangedPayload WebhookEventUser[any]
type WebhookEventUserIdentityCreatedPayload WebhookEventUser[any]
type WebhookEventUserIdentityDeletedPayload WebhookEventUser[any]
type WebhookEventUserActiveChangedPayload WebhookEventUser[any]
type WebhookEventUserLastLoginChangedPayload WebhookEventUser[any]
type WebhookEventUserMergedPayload WebhookEventUser[any]
type WebhookEventUserNameChangedPayload WebhookEventUser[any]
type WebhookEventUserNotesChangedPayload WebhookEventUser[any]
type WebhookEventUserOnlyPrivateCommentsChangedPayload WebhookEventUser[any]
type WebhookEventUserOrganizationMembershipCreatedPayload WebhookEventUser[any]
type WebhookEventUserOrganizationMembershipDeletedPayload WebhookEventUser[any]
type WebhookEventUserPasswordChangedPayload WebhookEventUser[any]
type WebhookEventUserPhotoChangedPayload WebhookEventUser[any]
type WebhookEventUserRoleChangedPayload WebhookEventUser[any]
type WebhookEventUserDeletedPayload WebhookEventUser[any]
type WebhookEventUserSuspendedChangedPayload WebhookEventUser[any]
type WebhookEventUserTagsChangedPayload WebhookEventUser[any]
type WebhookEventUserTimeZoneChangedPayload WebhookEventUser[any]

type WebhookEventOrganizationCreatedPayload WebhookEventOrganization[any]
type WebhookEventOrganizationCustomFieldChangedPayload WebhookEventOrganization[any]
type WebhookEventOrganizationDeletedPayload WebhookEventOrganization[any]
type WebhookEventOrganizationExternalIDChangedPayload WebhookEventOrganization[any]
type WebhookEventOrganizationNameChangedPayload WebhookEventOrganization[any]
type WebhookEventOrganizationTagsChangedPayload WebhookEventOrganization[any]

type WebhookEventArticlePublishedPayload WebhookEventArticle[any]
type WebhookEventArticleSubscriptionCreatedPayload WebhookEventArticle[any]
type WebhookEventArticleUnpublishedPayload WebhookEventArticle[any]
type WebhookEventArticleVoteCreatedPayload WebhookEventArticle[any]
type WebhookEventArticleVoteChangedPayload WebhookEventArticle[any]
type WebhookEventArticleVoteRemovedPayload WebhookEventArticle[any]
type WebhookEventArticleCommentCreatedPayload WebhookEventArticle[any]
type WebhookEventArticleCommentChangedPayload WebhookEventArticle[any]
type WebhookEventArticleCommentPublishedPayload WebhookEventArticle[any]
type WebhookEventArticleCommentUnpublishedPayload WebhookEventArticle[any]

type WebhookEventCommunityPostCreatedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostChangedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostPublishedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostUnpublishedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostSubscriptionCreatedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostVoteCreatedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostVoteChangedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostVoteRemovedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostCommentCreatedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostCommentChangedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostCommentPublishedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostCommentUnpublishedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostCommentVoteCreatedPayload WebhookEventCommunityPost[any]
type WebhookEventCommunityPostCommentVoteChangedPayload WebhookEventCommunityPost[any]

type WebhookEventAgentStateChangedPayload WebhookEventAgentState[any]
type WebhookEventAgentWorkItemAddedPayload WebhookEventAgentState[any]
type WebhookEventAgentWorkItemRemovedPayload WebhookEventAgentState[any]
type WebhookEventAgentMaxCapacityChangedPayload WebhookEventAgentState[any]
type WebhookEventAgentUnifiedStateChangedPayload WebhookEventAgentState[any]
type WebhookEventAgentChannelCreatedPayload WebhookEventAgentState[any]
type WebhookEventAgentChannelDeletedPayload WebhookEventAgentState[any]

type WebhookEventOmnichannelRoutingConfigFeatureChangedPayload WebhookEventOmnichannelRoutingConfig[any]
