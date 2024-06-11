package zendesk_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_SupportOrganizationMembership_Show_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/organization_membership/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/organization_memberships/11224062128028",
			},
		),
	})

	var exampleOrganizationMembershipID zendesk.OrganizationMembershipID = 11224062128028

	actual, err := z.Support().OrganizationMemberships().Show(ctx, exampleOrganizationMembershipID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleOrganizationMembershipID, actual.ID); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportOrganizationMembership_List_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/organization_membership/list_page1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/organization_memberships",
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/organization_membership/list_page2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/organization_memberships.json",
				Query: url.Values{
					"page[size]":  []string{"2"},
					"page[after]": []string{"aCursor="},
				},
			},
		),
	})

	expectedMembershipsLen := 4
	actualMembershipsLen := 0

	if err := z.Support().OrganizationMemberships().List(ctx,
		func(response zendesk.OrganizationMembershipsResponse) error {
			for range response.OrganizationMemberships {
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

func Test_SupportOrganizationMembership_ListByOrganizationID_200(t *testing.T) {
	ctx := context.Background()
	organizationID := zendesk.OrganizationID(12345)

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/organization_membership/list_by_organization_page1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/api/v2/organizations/%d/organization_memberships", organizationID),
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/organization_membership/list_by_organization_page2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/api/v2/organizations/%d/organization_memberships.json", organizationID),
				Query: url.Values{
					"page[size]":  []string{"2"},
					"page[after]": []string{"aCursor="},
				},
			},
		),
	})

	expectedMembershipsLen := 4
	actualMembershipsLen := 0

	if err := z.Support().OrganizationMemberships().ListByOrganizationID(
		ctx,
		organizationID,
		func(response zendesk.OrganizationMembershipsResponse) error {
			for range response.OrganizationMemberships {
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

func Test_SupportOrganizationMembership_Create_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusCreated,
				FilePath:   "test_files/responses/support/organization_membership/create_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/organization_memberships",
			},
		),
	})

	orgID := zendesk.OrganizationID(11224050897308)
	userID := zendesk.UserID(11174914247836)

	payload := zendesk.OrganizationMembershipPayload{
		OrganizationMembership: zendesk.OrganizationMembershipPayloadData{
			OrganizationID: orgID,
			UserID:         userID,
		},
	}

	_, err := z.Support().OrganizationMemberships().Create(ctx, payload)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SupportOrganizationMembership_Create_422(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusUnprocessableEntity,
				FilePath:   "test_files/responses/support/organization_membership/create_422.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/organization_memberships",
			},
		),
	})

	orgID := zendesk.OrganizationID(11224050897308)
	userID := zendesk.UserID(11174914247836)

	payload := zendesk.OrganizationMembershipPayload{
		OrganizationMembership: zendesk.OrganizationMembershipPayloadData{
			OrganizationID: orgID,
			UserID:         userID,
		},
	}

	_, err := z.Support().OrganizationMemberships().Create(ctx, payload)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
