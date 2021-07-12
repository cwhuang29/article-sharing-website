package handlers

import (
	"fmt"
	"net/http"

	"github.com/cwhuang29/article-sharing-website/constants"
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/gin-gonic/gin"
)

/*
 * 0: User doesn't bookmark this article
 * 1: User has bookmarked this article
 */

func GetUserBookmarkedArticles(c *gin.Context) {
	offset, err := getQueryOffset(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.QueryErr, "errBody": constants.QueryOffsetErr, "size": 0})
		return
	}

	limit, err := getQueryLimit(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.QueryErr, "errBody": constants.QueryLimitErr, "size": 0})
		return
	}

	userStatus, user := GetUserStatus(c)
	if userStatus < IsMember {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errHead": fmt.Sprintf(constants.LoginTo, "view home page"), "errBody": "", "size": 0})
		return
	}

	isAdmin := false
	if userStatus >= IsAdmin {
		isAdmin = true
	}

	dbFormatArticles := databases.GetUserBookmarkArticles(user.ID, offset, limit, isAdmin)
	articleList := make([]Article, len(dbFormatArticles))
	for i, a := range dbFormatArticles {
		articleList[i] = articleFormatDBToOverview(a)
	}

	c.JSON(http.StatusOK, gin.H{"articleList": articleList, "size": len(articleList)})
}

func Bookmark(c *gin.Context) {
	id, err := getParamArticleID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.QueryErr, "errBody": constants.QueryArticleIDErr})
		return
	}

	userStatus, user := GetUserStatus(c)
	if userStatus < IsMember {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errHead": fmt.Sprintf(constants.LoginTo, "save articles"), "errBody": ""})
		return
	}

	isBookmarked := databases.GetBookmarkStatus(user.ID, id)
	c.JSON(http.StatusOK, gin.H{"isBookmarked": isBookmarked})
}

func UpdateBookmark(c *gin.Context) {
	id, err := getParamArticleID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.QueryArticleIDErr})
		return
	}

	isBookmarked, err := getQueryBookmarked(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.ReloadAndRetry})
		return
	}

	userStatus, user := GetUserStatus(c)
	if userStatus < IsMember {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errHead": constants.LoginFirst, "errBody": ""})
		return
	}

	if ok := databases.UpdateBookmarkStatus(user.ID, id, isBookmarked); !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.ReloadAndRetry})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isBookmarked": isBookmarked})
}
