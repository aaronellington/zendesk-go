package zendesk_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_SupportTickets_AddTags_200(t *testing.T) {
	ctx := context.Background()
	ticketID := zendesk.TicketID(2000)
	tagToAdd := zendesk.Tag("test_tag_zendesk_go")

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/tickets/add_tags_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPut,
				Path:   fmt.Sprintf("/api/v2/tickets/%d/tags", ticketID),
			},
		),
	})

	response, err := z.Support().Tickets().AddTags(ctx, ticketID, zendesk.Tags{tagToAdd})
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Contains(response, tagToAdd) {
		t.Fatal("response did not contain tag")
	}
}

func Test_SupportTickets_AddTags_422(t *testing.T) {
	ctx := context.Background()
	closedTicketID := zendesk.TicketID(2000)
	tagToAdd := zendesk.Tag("test_tag_zendesk_go")

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusUnprocessableEntity,
				FilePath:   "test_files/responses/support/tickets/add_tags_422.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPut,
				Path:   fmt.Sprintf("/api/v2/tickets/%d/tags", closedTicketID),
			},
		),
	})

	_, err := z.Support().Tickets().AddTags(ctx, closedTicketID, zendesk.Tags{tagToAdd})
	if err == nil {
		t.Fatal(err)
	}

	zendeskGoError := &zendesk.Error{}
	isZendeskGoError := errors.As(err, &zendeskGoError)

	if !isZendeskGoError {
		t.Fatalf("expected a custom zendesk-go error, got: %v", err)
	}

	// Check to confirm that we got a 422 error
	if !zendeskGoError.ImmutableRecord() {
		t.Fatal("did not receive an immutable error")
	}
}

func Test_SupportTickets_RemoveTags_200(t *testing.T) {
	ctx := context.Background()
	ticketID := zendesk.TicketID(2000)
	tagToRemove := zendesk.Tag("test_tag_zendesk_go")

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/tickets/remove_tags_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/api/v2/tickets/%d/tags", ticketID),
			},
		),
	})

	response, err := z.Support().Tickets().RemoveTags(ctx, ticketID, zendesk.Tags{tagToRemove})
	if err != nil {
		t.Fatal(err)
	}

	if slices.Contains(response, tagToRemove) {
		t.Fatal("response stil contained tag after update")
	}
}

func Test_SupportTickets_RemoveTags_422(t *testing.T) {
	ctx := context.Background()
	closedTicketID := zendesk.TicketID(2000)
	tagToAdd := zendesk.Tag("test_tag_zendesk_go")

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusUnprocessableEntity,
				FilePath:   "test_files/responses/support/tickets/remove_tags_422.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/api/v2/tickets/%d/tags", closedTicketID),
			},
		),
	})

	_, err := z.Support().Tickets().RemoveTags(ctx, closedTicketID, zendesk.Tags{tagToAdd})
	if err == nil {
		t.Fatal(err)
	}

	zendeskGoError := &zendesk.Error{}
	isZendeskGoError := errors.As(err, &zendeskGoError)

	if !isZendeskGoError {
		t.Fatalf("expected a custom zendesk-go error, got: %v", err)
	}

	// Check to confirm that we got a 422 error
	if !zendeskGoError.ImmutableRecord() {
		t.Fatal("did not receive an immutable error")
	}
}

func Test_SupportTickets_HasTags_OK(t *testing.T) {
	targetTag := zendesk.Tag("zendesk_go_is_awesome")
	tagThatDoesntExist := zendesk.Tag("this_tag_doesnt_exist_on_the_resource")

	tags := zendesk.Tags{
		"tag_one",
		"tag_two",
		targetTag,
		"AnotherTag",
		"A_new_random_tag_with_numbers1234",
	}

	// Assert that we can find the existing Tag in the Tags array
	if !tags.HasTag(targetTag) {
		t.Fatal("did not find targetTag in array")
	}

	// Assert that we do not return true when the Tag does not exist in the Tags array
	if tags.HasTag(tagThatDoesntExist) {
		t.Fatal("falsely reporting that tag exists in Tags when it does not")
	}
}

func Test_SupportTickets_SearchTags_200(t *testing.T) {
	searchTerm := "support"

	ctx := context.Background()
	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/tickets/search_tags_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/autocomplete/tags",
				Query: url.Values{
					"name": []string{searchTerm},
				},
			},
		),
	})

	actualResult, err := z.Support().TicketTags().Search(ctx, searchTerm)
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := zendesk.Tags{
		"abarcy_support",
		"abecedarian_support",
		"abequitate_support",
		"agelast_premium_support",
		"acoustic_guitar_free_support",
		"ameliorate_free_support",
		"altruistic_free_support",
		"acquaintances_paid_support",
		"astronaut_tier_support",
		"abomination_support",
		"agoraphobia_inside_only_support",
		"acquiesce_a_support_survey",
		"astringent_support",
		"acerbic_partner_support",
		"auspicious_ticket_support",
	}

	if err := study.Assert(expectedResult, actualResult); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTickets_SearchTags_TermShorterThanAllowed(t *testing.T) {
	searchTerm := "s"

	ctx := context.Background()
	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/tickets/search_tags_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/autocomplete/tags",
				Query: url.Values{
					"name": []string{searchTerm},
				},
			},
		),
	})

	_, err := z.Support().TicketTags().Search(ctx, searchTerm)
	if err == nil {
		t.Fatal("expected to get an error")
	}

	if err.Error() != "invalid request - 'searchTerm' must be at least 2 characters" {
		t.Fatalf("did not get expected error - got: '%s', expected: '%s'",
			err.Error(),
			"invalid searchterm - searchterm must be at least 2 characters",
		)
	}
}
