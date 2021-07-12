package constants

const (
	GeneralErr    = "Error"
	UnexpectedErr = "Oops, this is unexpected"

	PermissionDenied = "You are not allowed to perform this action"

	TryAgain       = "Please try again"
	ReloadAndRetry = "Please reload the page and try again"
	GobackAndRetry = "Go back to previous page and try again"
	TryTooOften    = "You are trying too often"

	UserNotFound = "User Not Found"

	DatabaseErr = "An error occurred while writing to DB"

	LoginFirst = "You need to login first"
	LoginTo    = "Login to %s"

	PasswordIncorrect = "Password incorrect"

	ArticleNotFound  = "Article Not Found"
	ArticleCreateErr = "Create Article Failed"
	ArticleUpdateErr = "Update Article Failed"
	ArticleDeleteErr = "Delete Article Failed"

	QueryErr           = "Invalid Parameter"
	QueryEmptyErr      = "Parameter %s can not be empty"
	QueryMissingErr    = "Some values are missing"
	QueryArticleIDErr  = "Parameter articleId should be a positive integer"
	QueryOffsetErr     = "Parameter offset should be a non-negative integer"
	QueryLimitErr      = "Parameter limit should be a positive integer"
	QueryBookmarkedErr = "Parameter bookmarked should be either 0 or 1"
	QueryLikedErr      = "Parameter liked should be either 0 or 1"

	EmailNotFound         = "Email Not Found"
	EmailOccupied         = "This email is already registered"
	EmailChangeAnother    = "Please use another email"
	EmailOpenAgain        = "Please reopen the link from email"
	EmailLinkExpired      = "The link has expired"
	EmailOutdated         = "Perhaps you didn't open the latest email"
	EmailRequestAgain     = "Please request a reset password email again"
	EmailIsAddressCorrect = "Did you fill in the correct email address?"
	EmailTryLater         = "Please try again in one hour"
)
