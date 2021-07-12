package utils

import (
	"github.com/cwhuang29/article-sharing-website/constants"
	"github.com/cwhuang29/article-sharing-website/databases"
)

func DoesUserHasEmailQuota(id int) bool {
	count := databases.CountUserResetPasswordTokens(id)
	return count < constants.ResetPasswordMaxRetry
}
