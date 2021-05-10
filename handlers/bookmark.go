package handlers

import (
	"net/http"
	"strconv"

	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/gin-gonic/gin"
)

/*
 * 0: User doesn't bookmark this article
 * 1: User has bookmarked this article
 */

func GetUserBookmarkedArticles(c *gin.Context) {
	errHead := "Invalid Parameter"
	errBody := "Parameter articleId should be a positive integer."

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody, "size": 0})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody, "size": 0})
		return
	}

	userStatus, user := GetUserStatus(c)
	if userStatus < IsMember {
		errHead = "Login to view home page"
		errBody = ""
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errHead": errHead, "errBody": errBody, "size": 0})
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
	articleId, err := strconv.Atoi(c.Param("articleId"))
	if err != nil || articleId <= 0 {
		errHead := "Invalid Parameter"
		errBody := "Parameter articleId should be a positive integer."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	userStatus, user := GetUserStatus(c)
	if userStatus < IsMember {
		errHead := "Login to save articles"
		errBody := ""
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	isBookmarked := databases.GetBookmarkStatus(user.ID, articleId)
	c.JSON(http.StatusOK, gin.H{"isBookmarked": isBookmarked})
}

func UpdateBookmark(c *gin.Context) {
	errHead := "Oops, this is unexpected"
	errBody := "Please reload the page and try again."

	articleId, err := strconv.Atoi(c.Param("articleId"))
	if err != nil || articleId <= 0 {
		errBody = "Parameter articleId should be a positive integer."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	isBookmarked, err := strconv.Atoi(c.Query("bookmarked"))
	if err != nil || isBookmarked != 0 && isBookmarked != 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	userStatus, user := GetUserStatus(c)
	if userStatus < IsMember {
		errHead = "You need to login first"
		errBody = ""
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	if ok := databases.UpdateBookmarkStatus(user.ID, articleId, isBookmarked); !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isBookmarked": isBookmarked})
}
