package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/cwhuang29/article-sharing-website/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func GetUserStatus(c *gin.Context) (status UserStatus, cookieEmail string) {
	cookieEmail, _ = c.Cookie("login_email") // If no such cookie, c.Cookie() returns empty string with error `named cookie not present`
	cookieToken, _ := c.Cookie("login_token")
	adminEmail, _ := c.Cookie("is_admin")

	memberOrAdmin := IsMember
	if adminEmail != "" && cookieEmail == adminEmail && databases.IsAdminUser(adminEmail) {
		memberOrAdmin = IsAdmin
	}

	user := databases.GetUser(cookieEmail)
	creds := databases.GetLoginCredentials(user.ID)
	for _, cred := range creds {
		isEpr := utils.IsExpired(cred.LastLogin, cred.MaxAge)
		if cookieEmail == cred.User.Email && cookieToken == cred.Token && !isEpr {
			status = memberOrAdmin
			return
		}
	}

	cookieEmail = ""
	return
}

func isUserAdmin(c *gin.Context) bool {
	status, _ := GetUserStatus(c)
	return status >= IsAdmin
}

func getParaArticleId(c *gin.Context, key string) int {
	if c.Query(key) == "" {
		return 0
	}

	id, err := strconv.Atoi(c.Query(key))
	if err != nil || id <= 0 {
		return 0
	}

	return id
}

func fetchData(types, query string, offset, limit int, isAdmin bool) (articleList []Article, err error) {
	var dbFormatArticle []models.Article

	switch types {
	case "time":
		// For first time, load the weekly articles (all articles in the latest 7 days)
		if offset == 0 {
			today := time.Now().UTC().Truncate(24 * time.Hour)
			sevenDaysAgo := today.AddDate(0, 0, -7)
			tomorrow := today.AddDate(0, 0, 1)
			dbFormatArticle = databases.GetArticlesInATimePeriod(sevenDaysAgo, tomorrow, isAdmin)
		} else {
			dbFormatArticle = databases.GetArticles(offset, limit, isAdmin)
		}
	case "tag":
		dbFormatArticle = databases.GetSameTagArticles(query, offset, limit, isAdmin)
		for i := 0; i < len(dbFormatArticle); i++ {
			dbFormatArticle[i].Tags = databases.GetArticleTags(dbFormatArticle[i])
		}
	case "category":
		dbFormatArticle = databases.GetSameCategoryArticles(query, offset, limit, isAdmin)
	}

	for _, a := range dbFormatArticle {
		articleList = append(articleList, articleFormatDBToOverview(a))
	}
	return
}

func getPasswordResetTokenInstance(token string) models.Password {
	return databases.GetResetPasswordToken(token)
}

func doesUserHasEmailQuota(id int) bool {
	count := databases.CountUserResetPasswordTokens(id)
	return count < utils.ResetPasswordMaxRetry
}
