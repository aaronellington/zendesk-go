package zendesk

import (
	"context"
	"fmt"
	"net/http"
)

// https://developer.zendesk.com/api-reference/ticketing/side_conversation/side_conversation/
type SideConversationService struct {
	client *client
}

type SideConversationCreatePayload struct {
	Message SideConversationMessage `json:"message"`
}

type SideConversationMessage struct {
	Subject string                   `json:"subject"`
	Body    string                   `json:"body"`
	To      []SideConversationTarget `json:"to"`
}

type SideConversationTarget interface {
	SideConversationTarget()
}

type SideConversationTargetChildTicket struct {
	SupportGroupID GroupID `json:"support_group_id"`
	SupportAgentID UserID  `json:"support_agent_id,omitempty"`
}

func (s SideConversationTargetChildTicket) SideConversationTarget() {}

// https://developer.zendesk.com/api-reference/ticketing/side_conversation/side_conversation/#create-side-conversation
func (s *SideConversationService) Create(
	ctx context.Context,
	ticketID TicketID,
	payload SideConversationCreatePayload,
) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/tickets/%d/side_conversations", ticketID),
		structToReader(payload),
	)
	if err != nil {
		return err
	}

	return s.client.ZendeskRequest(request, nil)
}
