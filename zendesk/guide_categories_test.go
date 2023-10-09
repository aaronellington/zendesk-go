package zendesk_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func TestCategoryService_List(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/guide/categories.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/help_center/categories",
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
	})
	c := zendesk.CategoriesResponse{}

	if err := z.Guide().Categories().List(ctx, func(response zendesk.CategoriesResponse) error {
		c = response

		return nil
	}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		actual   any
		expected any
		result   bool
	}{
		{
			name:     "test1",
			actual:   c.Categories[0].Name,
			expected: "South Park",
			result:   true,
		},
		{
			name:     "test2",
			actual:   c.Categories[1].URL,
			expected: "https://southpark.fandom.com/wiki/Portal:Episodes",
			result:   true,
		},
		{
			name:     "test3",
			actual:   c.Categories[1].ID,
			expected: zendesk.CategoryID(2),
			result:   true,
		},
		{
			name:     "test4",
			actual:   c.Categories[4].Description,
			expected: "South Park Trivia",
			result:   true,
		},
		{
			name:     "test5",
			actual:   c.Categories[3].Outdated,
			expected: false,
			result:   true,
		},
		{
			name:     "test6",
			actual:   c.Categories[3].Outdated,
			expected: "false",
			result:   false,
		},
	}

	for _, test := range tests {
		if (test.actual != test.expected) == test.result {
			t.Errorf("test %s failed", test.name)
		}
	}
}
