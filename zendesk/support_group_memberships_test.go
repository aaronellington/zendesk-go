package zendesk_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_SupportGroupMembership_Show_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/group_membership/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/group_memberships/11292733779740",
			},
		),
	})

	var exampleGroupMembershipID zendesk.GroupMembershipID = 11292733779740

	actual, err := z.Support().GroupMemberships().Show(ctx, exampleGroupMembershipID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleGroupMembershipID, actual.GroupMembership.ID); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportGroupMembership_List_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/group_membership/list_page1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/group_memberships",
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/group_membership/list_page2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/group_memberships.json",
				Query: url.Values{
					"page[size]":  []string{"1"},
					"page[after]": []string{"aCursor="},
				},
			},
		),
	})

	expectedMembershipsLen := 2
	actualMembershipsLen := 0

	if err := z.Support().GroupMemberships().List(ctx,
		func(response zendesk.GroupMembershipsResponse) error {
			for range response.GroupMemberships {
				actualMembershipsLen++
			}

			return nil
		},
	); err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(expectedMembershipsLen, actualMembershipsLen); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportGroupMembership_Create_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusCreated,
				FilePath:   "test_files/responses/support/group_membership/create_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/group_memberships",
			},
		),
	})

	groupID := zendesk.GroupID(11328996577564)
	userID := zendesk.UserID(11174914247836)

	_, err := z.Support().GroupMemberships().Create(ctx, userID, groupID)
	if err != nil {
		t.Fatal(err)
	}
}
