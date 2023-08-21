package zendesk

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
)

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

type (
	ArticleID           uint64
	CategoryID          uint64
	ChatAccountID       uint64
	ChatEngagementID    string
	ChatID              string
	CustomFieldOptionID uint64
	CustomRoleID        uint64
	GroupID             uint64
	GroupMembershipID   uint64
	OrganizationID      uint64
	PermissionGroupID   uint64
	ScheduleID          uint64
	SectionID           uint64
	SuspendedTicketID   uint64
	Tag                 string
	TicketAuditEventID  uint64
	TicketAuditID       uint64
	TicketFieldID       uint64
	TicketFormID        uint64
	TicketID            uint64
	UserFieldID         uint64
	UserID              uint64
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

func structToReader(x any) io.Reader {
	payloadBytes, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}

	return bytes.NewReader(payloadBytes)
}
