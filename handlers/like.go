package handlers

import (
	"net/http"
	"strconv"

	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/gin-gonic/gin"
)

/*
 * 0: User doesn't like this article
 * 1: User has liked this article
 */

func Like(c *gin.Context) {
	articleId, err := strconv.Atoi(c.Param("articleId"))
	if err != nil || articleId <= 0 {
		errHead := "Invalid Parameter"
		errBody := "Parameter articleId should be a positive integer."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	userStatus, user := GetUserStatus(c)
	if userStatus < IsMember {
		errHead := "You need to login first"
		errBody := ""
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	isLiked := databases.GetLikeStatus(user.ID, articleId)
	c.JSON(http.StatusOK, gin.H{"isLiked": isLiked})
}

func UpdateLike(c *gin.Context) {
	errHead := "Oops, this is unexpected"
	errBody := "Please reload the page and try again."

	articleId, err := strconv.Atoi(c.Param("articleId"))
	if err != nil || articleId <= 0 {
		errHead = "Invalid Parameter"
		errBody = "Parameter articleId should be a positive integer."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	isLiked, err := strconv.Atoi(c.Query("liked"))
	if err != nil || isLiked != 0 && isLiked != 1 {
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

	if ok := databases.UpdateLikeStatus(user.ID, articleId, isLiked); !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isLiked": isLiked})
}
