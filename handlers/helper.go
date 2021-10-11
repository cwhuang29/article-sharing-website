package handlers

import (
	"fmt"
	"time"

	"github.com/cwhuang29/article-sharing-website/constants"
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/cwhuang29/article-sharing-website/utils"
	"github.com/gin-gonic/gin"
)

func isUserAdmin(c *gin.Context) bool {
	status, _ := GetUserStatus(c)
	return status >= IsAdmin
}

func getURLPara(c *gin.Context, key string) string {
	return c.Param(key)
}

func getQueryPara(c *gin.Context, key string) string {
	return c.DefaultQuery(key, "")
}

func getParamArticleID(c *gin.Context) (int, error) {
	return utils.Str2PosInt(getURLPara(c, constants.ParamArticleID))
}

func getQueryArticleID(c *gin.Context) (int, error) {
	return utils.Str2PosInt(getQueryPara(c, constants.QueryArticleID))
}

func getQueryOffset(c *gin.Context) (int, error) {
	return utils.Str2Int(getQueryPara(c, constants.QueryOffset))
}

func getQueryLimit(c *gin.Context) (int, error) {
	return utils.Str2PosInt(getQueryPara(c, constants.QueryLimit))
}

func getQueryLiked(c *gin.Context) (int, error) {
	isLiked, err := utils.Str2Int(getQueryPara(c, constants.QueryLiked))
	if err != nil || (isLiked != 0 && isLiked != 1) {
		err = fmt.Errorf(constants.QueryLikedErr)
	}
	return isLiked, err
}

func getQueryBookmarked(c *gin.Context) (int, error) {
	isBookmarked, err := utils.Str2Int(getQueryPara(c, constants.QueryBookmarked))
	if err != nil || (isBookmarked != 0 && isBookmarked != 1) {
		err = fmt.Errorf(constants.QueryBookmarkedErr)
	}
	return isBookmarked, err
}

func GetUserStatus(c *gin.Context) (status UserStatus, user models.User) {
	cookieEmail, _ := c.Cookie(constants.CookieLoginEmail) // If no such cookie, c.Cookie() returns empty string with error `named cookie not present`
	cookieToken, _ := c.Cookie(constants.CookieLoginToken)
	adminEmail, _ := c.Cookie(constants.CookieIsAdmin)

	memberOrAdmin := IsMember
	if adminEmail != "" && cookieEmail == adminEmail && databases.IsAdminUser(adminEmail) {
		memberOrAdmin = IsAdmin
	}

	user = databases.GetUser(cookieEmail)
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

func fetchData(types, query string, offset, limit int, isAdmin bool) (articleList []Article, err error) {
	var dbFormatArticles []models.Article

	switch types {
	case "time":
		// For the first time, load the weekly articles (all articles in the latest 7 days)
		if offset == 0 {
			today := time.Now().UTC().Truncate(24 * time.Hour)
			sevenDaysAgo := today.AddDate(0, 0, -7)
			tomorrow := today.AddDate(0, 0, 1)
			dbFormatArticles = databases.GetArticlesInATimePeriod(sevenDaysAgo, tomorrow, isAdmin)
		} else {
			dbFormatArticles = databases.GetArticles(offset, limit, isAdmin)
		}
	case "tag":
		dbFormatArticles = databases.GetSameTagArticles(query, offset, limit, isAdmin)
		for i := 0; i < len(dbFormatArticles); i++ {
			dbFormatArticles[i].Tags = databases.GetArticleTags(dbFormatArticles[i])
		}
	case "category":
		dbFormatArticles = databases.GetSameCategoryArticles(query, offset, limit, isAdmin)
	}

	articleList = make([]Article, len(dbFormatArticles))
	for i, a := range dbFormatArticles {
		articleList[i] = articleFormatDBToOverview(a)
	}
	return
}
