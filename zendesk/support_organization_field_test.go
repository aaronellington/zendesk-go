package zendesk_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_SupportOrganizationField_Show_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/organization_field/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/organization_fields/4321",
			},
		),
	})

	var exampleOrganizationFieldID zendesk.OrganizationFieldID = 4321

	actual, err := z.Support().OrganizationFields().Show(ctx, exampleOrganizationFieldID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleOrganizationFieldID, actual.ID); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportOrganizationField_Create_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusCreated,
				FilePath:   "test_files/responses/support/organization_field/create_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/organization_fields",
			},
		),
	})

	payloadData := struct {
		Active      bool   `json:"active"`
		Description string `json:"description"`
	}{
		Active:      true,
		Description: "A new field",
	}

	payload := zendesk.OrganizationFieldPayload{
		OrganizationField: payloadData,
	}

	var exampleOrganizationFieldID zendesk.OrganizationFieldID = 75

	actual, err := z.Support().OrganizationFields().Create(ctx, payload)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleOrganizationFieldID, actual.OrganizationField.ID); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportOrganizationField_Delete_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/organization_field/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/organization_fields/4321",
			},
		),
	})

	var exampleOrganizationFieldID zendesk.OrganizationFieldID = 4321

	actual, err := z.Support().OrganizationFields().Show(ctx, exampleOrganizationFieldID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleOrganizationFieldID, actual.ID); err != nil {
		t.Fatal(err)
	}
}
