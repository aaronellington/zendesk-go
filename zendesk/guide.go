package zendesk

// https://developer.zendesk.com/api-reference/help_center/help-center-api/introduction/
type GuideService struct {
	categoriesService *CategoriesService
	sectionsService   *SectionsService
	articlesService   *ArticlesService
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/categories/
func (s *GuideService) Categories() *CategoriesService {
	return s.categoriesService
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/sections/
func (s *GuideService) Sections() *SectionsService {
	return s.sectionsService
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/
func (s *GuideService) Articles() *ArticlesService {
	return s.articlesService
}
