package zendesk

// https://developer.zendesk.com/api-reference/help_center/help-center-api/introduction/
type GuideService struct {
	articlesService *ArticlesService
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/
func (s *GuideService) Articles() *ArticlesService {
	return s.articlesService
}
