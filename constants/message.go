package constants

const (
	UnexpectedErr = "Oops, this is unexpected"

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

	ParameterErr          = "Invalid Parameter"
	ParameterEmptyErr     = "Parameter %s can not be empty"
	ParameterArticleIDErr = "Parameter articleId is a positive integer"
	ParameterMissingErr   = "Some values are missing"

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
