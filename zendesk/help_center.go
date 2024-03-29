package zendesk

type HelpCenterService struct {
	accountCustomClaims  *HelpCenterAccountCustomClaimsService
	articleAttachments   *HelpCenterArticleAttachmentsService
	articleComments      *HelpCenterArticleCommentsService
	articleLabels        *HelpCenterArticleLabelsService
	articles             *HelpCenterArticlesService
	badgeAssignments     *HelpCenterBadgeAssignmentsService
	badgeCategories      *HelpCenterBadgeCategoriesService
	badges               *HelpCenterBadgesService
	categories           *HelpCenterCategoriesService
	contentSubscriptions *HelpCenterContentSubscriptionsService
	contentTags          *HelpCenterContentTagsService
	jwts                 *HelpCenterJWTsService
	permissionGroups     *HelpCenterPermissionGroupsService
	postComments         *HelpCenterPostCommentsService
	posts                *HelpCenterPostsService
	search               *HelpCenterSearchService
	sections             *HelpCenterSectionsService
	theming              *HelpCenterThemingService
	topics               *HelpCenterTopicsService
	translations         *HelpCenterTranslationsService
	userImages           *HelpCenterUserImagesService
	userSegments         *HelpCenterUserSegmentsService
	userSubscriptions    *HelpCenterUserSubscriptionsService
	votes                *HelpCenterVotesService
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/account_custom_claims/
func (s *HelpCenterService) AccountCustomClaims() *HelpCenterAccountCustomClaimsService {
	return s.accountCustomClaims
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/article_attachments/
func (s *HelpCenterService) ArticleAttachments() *HelpCenterArticleAttachmentsService {
	return s.articleAttachments
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/article_comments/
func (s *HelpCenterService) ArticleComments() *HelpCenterArticleCommentsService {
	return s.articleComments
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/article_labels/
func (s *HelpCenterService) ArticleLabels() *HelpCenterArticleLabelsService {
	return s.articleLabels
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/
func (s *HelpCenterService) Articles() *HelpCenterArticlesService {
	return s.articles
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/badge_assignments/
func (s *HelpCenterService) BadgeAssignments() *HelpCenterBadgeAssignmentsService {
	return s.badgeAssignments
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/badge_categories/
func (s *HelpCenterService) BadgeCategories() *HelpCenterBadgeCategoriesService {
	return s.badgeCategories
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/badges/
func (s *HelpCenterService) Badges() *HelpCenterBadgesService {
	return s.badges
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/categories/
func (s *HelpCenterService) Categories() *HelpCenterCategoriesService {
	return s.categories
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/content_subscriptions/
func (s *HelpCenterService) ContentSubscriptions() *HelpCenterContentSubscriptionsService {
	return s.contentSubscriptions
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/content_tags/
func (s *HelpCenterService) ContentTags() *HelpCenterContentTagsService {
	return s.contentTags
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/help_center_jwts/
func (s *HelpCenterService) JWTs() *HelpCenterJWTsService {
	return s.jwts
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/permission_groups/
func (s *HelpCenterService) PermissionGroups() *HelpCenterPermissionGroupsService {
	return s.permissionGroups
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/post_comments/
func (s *HelpCenterService) PostComments() *HelpCenterPostCommentsService {
	return s.postComments
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/posts/
func (s *HelpCenterService) Posts() *HelpCenterPostsService {
	return s.posts
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/search/
func (s *HelpCenterService) Search() *HelpCenterSearchService {
	return s.search
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/sections/
func (s *HelpCenterService) Sections() *HelpCenterSectionsService {
	return s.sections
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/theming/
func (s *HelpCenterService) Theming() *HelpCenterThemingService {
	return s.theming
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/topics/
func (s *HelpCenterService) Topics() *HelpCenterTopicsService {
	return s.topics
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/translations/
func (s *HelpCenterService) Translations() *HelpCenterTranslationsService {
	return s.translations
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/user_images/
func (s *HelpCenterService) UserImages() *HelpCenterUserImagesService {
	return s.userImages
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/user_segments/
func (s *HelpCenterService) UserSegments() *HelpCenterUserSegmentsService {
	return s.userSegments
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/user_subscriptions/
func (s *HelpCenterService) UserSubscriptions() *HelpCenterUserSubscriptionsService {
	return s.userSubscriptions
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/votes/
func (s *HelpCenterService) Votes() *HelpCenterVotesService {
	return s.votes
}
