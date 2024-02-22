package zendesk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
)

const timeFormat = "2006-01-02T15:04:05Z"

const (
	PriorityUrgent = "urgent"
	PriorityHigh   = "high"
	PriorityNormal = "normal"
	PriorityLow    = "low"
)

const (
	StatusNew     = "new"
	StatusOpen    = "open"
	StatusPending = "pending"
	StatusHold    = "hold"
	StatusSolved  = "solved"
	StatusClosed  = "closed"
	StatusDeleted = "deleted"
)

type UserSideload string

const (
	UserSideloadIdentities      UserSideload = "identities"
	UserSideloadOrganizations   UserSideload = "organizations"
	UserSideloadRoles           UserSideload = "roles"
	UserSideloadAbilities       UserSideload = "abilities"
	UserSideloadGroups          UserSideload = "groups"
	UserSideloadOpenTicketCount UserSideload = "open_ticket_count"
)

type TicketCommentSideload string

const (
	TicketCommentSideloadUsers TicketCommentSideload = "users"
)

type TicketSideload string

const (
	TicketSideloadDates TicketSideload = "dates"
)

type MalwareScanResult string

const (
	MalwareFound        MalwareScanResult = "malware_found"
	MalwareNotFound     MalwareScanResult = "malware_not_found"
	MalwareFailedToScan MalwareScanResult = "failed_to_scan"
	MalwareNotScanned   MalwareScanResult = "not_scanned"
)

type SatisfactionRatingScore string

const (
	SatisfactionRatingScoreOffered                SatisfactionRatingScore = "offered"
	SatisfactionRatingScoreUnOffered              SatisfactionRatingScore = "unoffered"
	SatisfactionRatingScoreReceived               SatisfactionRatingScore = "received"
	SatisfactionRatingScoreReceivedWithComment    SatisfactionRatingScore = "received_with_comment"
	SatisfactionRatingScoreReceivedWithoutComment SatisfactionRatingScore = "received_without_comment"
	SatisfactionRatingScoreGood                   SatisfactionRatingScore = "good"
	SatisfactionRatingScoreGoodWithComment        SatisfactionRatingScore = "good_with_comment"
	SatisfactionRatingScoreGoodWithoutComment     SatisfactionRatingScore = "good_without_comment"
	SatisfactionRatingScoreBad                    SatisfactionRatingScore = "bad"
	SatisfactionRatingScoreBadWithComment         SatisfactionRatingScore = "bad_with_comment"
	SatisfactionRatingScoreBadWithoutComment      SatisfactionRatingScore = "bad_without_comment"
)

type AuditLogAction string

const (
	Create   AuditLogAction = "create"
	Destroy  AuditLogAction = "destroy"
	Exported AuditLogAction = "exported"
	Login    AuditLogAction = "login"
	Update   AuditLogAction = "update"
)

type LiveChatTimeWindow uint64

const (
	LiveChatTimeWindow30Minutes LiveChatTimeWindow = 30
	LiveChatTimeWindow60Minutes LiveChatTimeWindow = 60
)

type LiveChatMetricKeyChat string

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/metrics/#chat-metrics
const (
	LiveChatMetricKeyIncomingChats    LiveChatMetricKeyChat = "incoming_chats"
	LiveChatMetricKeyAssignedChats    LiveChatMetricKeyChat = "assigned_chats"
	LiveChatMetricKeyMissedChats      LiveChatMetricKeyChat = "missed_chats"
	LiveChatMetricKeyActiveChats      LiveChatMetricKeyChat = "active_chats"
	LiveChatMetricKeyWaitingTimeAvg   LiveChatMetricKeyChat = "waiting_time_avg"
	LiveChatMetricKeyWaitingTimeMax   LiveChatMetricKeyChat = "waiting_time_max"
	LiveChatMetricKeyChatDurationAvg  LiveChatMetricKeyChat = "chat_duration_avg"
	LiveChatMetricKeyChatDurationMax  LiveChatMetricKeyChat = "chat_duration_max"
	LiveChatMetricKeyResponseTimeAvg  LiveChatMetricKeyChat = "response_time_avg"
	LiveChatMetricKeyResponseTimeMax  LiveChatMetricKeyChat = "response_time_max"
	LiveChatMetricKeySatisfactionGood LiveChatMetricKeyChat = "satisfaction_good"
	LiveChatMetricKeySatisfactionBad  LiveChatMetricKeyChat = "satisfaction_bad"
)

type LiveChatMetricKeyAgent string

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/metrics/#agent-metrics
const (
	LiveChatMetricKeyAgentsOnline    LiveChatMetricKeyAgent = "agents_online"
	LiveChatMetricKeyAgentsAway      LiveChatMetricKeyAgent = "agents_away"
	LiveChatMetricKeyAgentsInvisible LiveChatMetricKeyAgent = "agents_invisible"
)

type (
	AccountID                uint64
	ActorID                  int64
	ArticleID                uint64
	AttachmentID             uint64
	AuditID                  uint64
	AuditLogID               uint64
	AutomationID             uint64
	BrandID                  uint64
	CategoryID               uint64
	ChatAccountID            uint64
	ChatEngagementID         string
	ChatID                   string
	CustomFieldOptionID      uint64
	CustomRoleID             uint64
	CustomStatusID           uint64
	GroupID                  uint64
	GroupMembershipID        uint64
	LiveChatOAuthClientID    uint64
	OrganizationID           uint64
	OrganizationFieldID      uint64
	OrganizationMembershipID uint64
	PermissionGroupID        uint64
	PostID                   string
	TopicID                  string
	ReasonCode               uint64
	ReasonID                 uint64
	SatisfactionRatingID     uint64
	ScheduleID               uint64
	SectionID                uint64
	SourceID                 int64
	SuspendedTicketID        uint64
	TicketAuditEventID       uint64
	TicketAuditID            uint64
	TicketCommentID          uint64
	TicketFieldID            uint64
	TicketFormID             uint64
	TicketID                 uint64
	UploadToken              string
	UserFieldID              uint64
	UserID                   int64
	IdentityID               uint64
	UserSegmentID            uint64
	WebhookEventID           string
)

func (userID *UserID) UnmarshalJSON(b []byte) error {
	// Try it as a uint64 first
	var targetInt64 int64
	if err := json.Unmarshal(b, &targetInt64); err == nil {
		*userID = UserID(targetInt64)

		return nil
	}

	// Only try it as a string as a last resort
	var targetString string
	if err := json.Unmarshal(b, &targetString); err != nil {
		return err
	}

	typeInt64, _ := strconv.ParseInt(targetString, 10, 64)
	*userID = UserID(typeInt64)

	return nil
}

func (tag Tag) Validate() error {
	if strings.Contains(string(tag), " ") {
		return errors.New("zendesk tag cannot contain spaces")
	}

	return nil
}

func structToReader(x any) io.Reader {
	payloadBytes, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}

	return bytes.NewReader(payloadBytes)
}
