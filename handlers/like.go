package handlers

import (
	"net/http"
	"strconv"

	"github.com/cwhuang29/article-sharing-website/constants"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.ParameterErr, "errBody": constants.ParameterArticleIDErr})
		return
	}

	userStatus, user := GetUserStatus(c)
	if userStatus < IsMember {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errHead": constants.LoginFirst, "errBody": ""})
		return
	}

	isLiked := databases.GetLikeStatus(user.ID, articleId)
	c.JSON(http.StatusOK, gin.H{"isLiked": isLiked})
}

func UpdateLike(c *gin.Context) {
	articleId, err := strconv.Atoi(c.Param("articleId"))
	if err != nil || articleId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.ParameterErr, "errBody": constants.ParameterArticleIDErr})
		return
	}

	isLiked, err := strconv.Atoi(c.Query("liked"))
	if err != nil || isLiked != 0 && isLiked != 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.ReloadAndRetry})
		return
	}

	userStatus, user := GetUserStatus(c)
	if userStatus < IsMember {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errHead": constants.LoginFirst, "errBody": ""})
		return
	}

	if ok := databases.UpdateLikeStatus(user.ID, articleId, isLiked); !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.ReloadAndRetry})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isLiked": isLiked})
}
