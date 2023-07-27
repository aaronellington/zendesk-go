package zendesk

// https://developer.zendesk.com/api-reference/help_center/help-center-api/introduction/
type GuideService struct {
	categoriesService *CategoryService
	sectionsService   *SectionService
	articlesService   *ArticleService
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/categories/
func (s *GuideService) Categories() *CategoryService {
	return s.categoriesService
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/sections/
func (s *GuideService) Sections() *SectionService {
	return s.sectionsService
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/
func (s *GuideService) Articles() *ArticleService {
	return s.articlesService
}
