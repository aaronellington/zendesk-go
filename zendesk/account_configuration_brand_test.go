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

func Test_AccountConfiguration_Brands_Show_200(t *testing.T) {
	ctx := context.Background()

	var expectedBrandID zendesk.BrandID = 987654321

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/brand/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/api/v2/brands/%d", expectedBrandID),
			},
		),
	})

	actual, err := z.AccountConfiguration().Brands().Show(ctx, expectedBrandID)
	if err != nil {
		t.Fatal(err)
	}

	if actual.ID != expectedBrandID {
		t.Fatalf("expected ID: %d - got ID: %d", expectedBrandID, actual.ID)
	}
}

func Test_AccountConfiguration_Brands_List_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/brand/list_page1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/brands",
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/brand/list_page2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/brands.json",
				Query: url.Values{
					"page[size]":  []string{"100"},
					"page[after]": []string{"aCursor="},
				},
			},
		),
	})

	actual := []zendesk.Brand{}

	if err := z.AccountConfiguration().Brands().List(ctx, func(response zendesk.BrandsResponse) error {
		actual = append(actual, response.Brands...)

		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if len(actual) != 3 {
		t.Fatalf("expected 3 brands, got %d", len(actual))
	}
}
