package zendesk_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_Guide_Articles_Show_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/guide/articles/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/help_center/articles/123",
			},
		),
	})

	var exampleArticleID zendesk.ArticleID = 123

	actual, err := z.Guide().Articles().Show(ctx, exampleArticleID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleArticleID, actual.Article.ID); err != nil {
		t.Fatal(err)
	}
}
