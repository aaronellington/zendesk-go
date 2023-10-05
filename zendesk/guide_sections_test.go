package zendesk_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func TestSectionService_List(t *testing.T) {

	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/guide/sections.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/help_center/sections",
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
	})
	c := zendesk.SectionsResponse{}

	if err := z.Guide().Sections().List(ctx, func(response zendesk.SectionsResponse) error {

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
			name: "test1", actual: c.Sections[0].Name, expected: "Main Characters", result: true,
		},
		{
			name: "test2", actual: c.Sections[1].ParentSectionID, expected: nil, result: true,
		},
		{
			name: "test3", actual: c.Sections[1].CategoryID, expected: zendesk.CategoryID(2), result: false,
		},
		{
			name: "test4", actual: c.Sections[0].CategoryID, expected: zendesk.CategoryID(1), result: true,
		},
		{
			name: "test5", actual: c.Sections[0].CategoryID, expected: zendesk.SectionID(1), result: false,
		},
		{
			name: "test6", actual: len(c.Sections), expected: 5, result: false,
		},
		{
			name: "test7", actual: c.Sections[2].ThemeTemplate, expected: "default", result: true,
		},
	}

	for _, test := range tests {

		if (test.actual != test.expected) == test.result {

			t.Errorf("test %s failed", test.name)

		}

	}
}
