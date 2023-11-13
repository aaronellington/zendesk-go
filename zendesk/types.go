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

type UserEndpointSideload string

const (
	UserEndpointSideloadIdentities      UserEndpointSideload = "identities"
	UserEndpointSideloadOrganizations   UserEndpointSideload = "organizations"
	UserEndpointSideloadRoles           UserEndpointSideload = "roles"
	UserEndpointSideloadAbilities       UserEndpointSideload = "abilities"
	UserEndpointSideloadGroups          UserEndpointSideload = "groups"
	UserEndpointSideloadOpenTicketCount UserEndpointSideload = "open_ticket_count"
)

type MalwareScanResult string

const (
	MalwareFound        MalwareScanResult = "malware_found"
	MalwareNotFound     MalwareScanResult = "malware_not_found"
	MalwareFailedToScan MalwareScanResult = "failed_to_scan"
	MalwareNotScanned   MalwareScanResult = "not_scanned"
)

type AuditLogAction string

const (
	Create   AuditLogAction = "create"
	Destroy  AuditLogAction = "destroy"
	Exported AuditLogAction = "exported"
	Login    AuditLogAction = "login"
	Update   AuditLogAction = "update"
)

type (
	ActorID             int64
	ArticleID           uint64
	AttachmentID        uint64
	AuditID             uint64
	AuditLogID          uint64
	CategoryID          uint64
	ChatAccountID       uint64
	ChatEngagementID    string
	ChatID              string
	CustomFieldOptionID uint64
	CustomRoleID        uint64
	CustomStatusID      uint64
	GroupID             uint64
	GroupMembershipID   uint64
	OrganizationID      uint64
	OrganizationFieldID uint64
	PermissionGroupID   uint64
	ScheduleID          uint64
	SectionID           uint64
	SourceID            int64
	SuspendedTicketID   uint64
	TicketAuditEventID  uint64
	TicketAuditID       uint64
	TicketCommentID     uint64
	TicketFieldID       uint64
	TicketFormID        uint64
	TicketID            uint64
	UploadToken         string
	UserFieldID         uint64
	UserID              uint64
	UserIdentityID      uint64
	UserSegmentID       uint64
)

func (userID *UserID) UnmarshalJSON(b []byte) error {
	// Try it as a uint64 first
	var targetUint64 uint64
	if err := json.Unmarshal(b, &targetUint64); err == nil {
		*userID = UserID(targetUint64)

		return nil
	}

	// Only try it as a string as a last resort
	var targetString string
	if err := json.Unmarshal(b, &targetString); err != nil {
		return err
	}

	typeUint64, _ := strconv.ParseUint(targetString, 0, 64)
	*userID = UserID(typeUint64)

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
