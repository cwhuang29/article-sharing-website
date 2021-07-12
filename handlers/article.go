package handlers

import (
	"net/http"
	"strconv"

	"github.com/cwhuang29/article-sharing-website/constants"
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateArticleView(c *gin.Context) {
	uuid := utils.GetUUID()
	c.SetCookie(constants.CookieCSRFToken, uuid, constants.CsrfTokenAge, "/", "", true, true)
	c.HTML(http.StatusOK, "editor.html", gin.H{
		"currPageCSS": "css/editor.css",
		"csrfToken":   uuid,
		"function":    "create",
		"title":       "Create New Post",
	})
}

func UpdateArticleView(c *gin.Context) {
	id, err := getQueryArticleID(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"errHead":     constants.QueryArticleIDErr,
			"errBody":     constants.TryAgain,
		})
		return
	}

	dbFormatArticle := databases.GetArticle(id, true)
	if dbFormatArticle.ID == 0 {
		c.HTML(http.StatusNotFound, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"errHead":     constants.ArticleNotFound,
			"errBody":     constants.TryAgain,
		})
		return
	}

	article := articleFormatDBToDetailed(dbFormatArticle, false)
	uuid := utils.GetUUID()

	c.SetCookie(constants.CookieCSRFToken, uuid, constants.CsrfTokenAge, "/", "", true, true)
	c.HTML(http.StatusOK, "editor.html", gin.H{
		"currPageCSS":  "css/editor.css",
		"csrfToken":    uuid,
		"function":     "update",
		"adminOnly":    article.AdminOnly,
		"title":        "Edit: " + article.Title,
		"articleTitle": article.Title,
		"subtitle":     article.Subtitle,
		"date":         article.Date,
		"author":       article.Authors,
		"category":     article.Category,
		"tags":         article.Tags,
		"outline":      article.Outline,
		"coverPhoto":   article.CoverPhoto,
		"content":      article.Content,
	})
}

func CreateArticle(c *gin.Context) {
	newArticle, invalids, err := handleForm(c)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.ArticleCreateErr, "errBody": constants.TryAgain, "errTags": invalids})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.ArticleCreateErr, "errBody": err.Error()})
		return
	}

	id, res := databases.SubmitArticle(*newArticle, "create")
	if !res {
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": constants.ArticleCreateErr, "errBody": constants.DatabaseErr})
		return
	}

	logrus.Infof("Create article with id %v", id)
	c.Header("Location", "/articles/browse?articleId="+strconv.Itoa(id)) // With Location header and status code 3XX (not 2XX), response.redirected becomes true
	c.JSON(http.StatusCreated, gin.H{"articleId": id})
}

func UpdateArticle(c *gin.Context) {
	id, err := getQueryArticleID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": false, "errHead": constants.QueryArticleIDErr, "errBody": constants.TryAgain})
		return
	}

	if succeed := databases.IsArticleExists(id, true); !succeed {
		c.JSON(http.StatusNotFound, gin.H{"bindingError": false, "errHead": constants.ArticleNotFound, "errBody": constants.TryAgain})
		return
	}

	newArticle, invalids, err := handleForm(c)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.ArticleUpdateErr, "errBody": constants.TryAgain, "errTags": invalids})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.ArticleUpdateErr, "errBody": err.Error()})
		return
	}
	newArticle.ID = id

	id, res := databases.SubmitArticle(*newArticle, "update")
	if !res {
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": constants.ArticleUpdateErr, "errBody": constants.DatabaseErr})
		return
	}

	logrus.Infof("Update article with id %v", id)
	c.Header("Location", constants.URLBrowse+"?articleId="+strconv.Itoa(id))
	c.JSON(http.StatusCreated, gin.H{"articleId": id})
}

func DeleteArticle(c *gin.Context) {
	id, err := getQueryArticleID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.QueryArticleIDErr, "errBody": constants.TryAgain})
		return
	}

	if res := databases.DeleteArticle(id, true); !res {
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": constants.ArticleDeleteErr, "errBody": constants.TryAgain})
		return
	}

	logrus.Infof("Delete article with id %v", id)
	c.Status(http.StatusNoContent)
}
