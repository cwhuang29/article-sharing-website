package handlers

import (
	"net/http"
	"strconv"

	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateArticleView(c *gin.Context) {
	uuid := utils.GetUUID()
	c.SetCookie("csrf_token", uuid, utils.CsrfTokenAge, "/", "", true, true)
	c.HTML(http.StatusOK, "editor.html", gin.H{
		"currPageCSS": "css/editor.css",
		"csrfToken":   uuid,
		"function":    "create",
		"title":       "Create New Post",
	})
}

func UpdateArticleView(c *gin.Context) {
	id := getParaId(c, "articleId")
	if id == 0 {
		errHead := "Article ID is An Positive Integer"
		errBody := "Please try again."
		c.HTML(http.StatusBadRequest, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"errHead":     errHead,
			"errBody":     errBody,
		})
		return
	}

	dbFormatArticle := databases.GetArticle(id, true)
	if dbFormatArticle.ID == 0 {
		errHead := "Article Not Found"
		errBody := "Please try again."
		c.HTML(http.StatusNotFound, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"errHead":     errHead,
			"errBody":     errBody,
		})
		return
	}

	article := articleFormatDBToDetailed(dbFormatArticle, false)
	uuid := utils.GetUUID()

	c.SetCookie("csrf_token", uuid, utils.CsrfTokenAge, "/", "", true, true)
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
	errHead := "Create Article Failed"
	errBody := "Please try again."

	newArticle, invalids, err := handleForm(c)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody, "errTags": invalids})
		return
	} else if err != nil {
		errBody = err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	id, res := databases.SubmitArticle(newArticle, "create")
	if !res {
		errBody = "An error occurred while writing to DB."
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}
	logrus.Infof("Create article with id %v", id)
	c.Header("Location", "/articles/browse?articleId="+strconv.Itoa(id)) // With Location header and status code 3XX (not 2XX), response.redirected becomes true
	c.JSON(http.StatusCreated, gin.H{"articleId": id})
}

func UpdateArticle(c *gin.Context) {
	id := getParaId(c, "articleId")
	if id == 0 {
		errHead := "Article ID is An Positive Integer"
		errBody := "Please try again."
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}

	if succeed := databases.IsArticleExists(id, true); succeed != true {
		errHead := "Article Not Found"
		errBody := "Please try again."
		c.JSON(http.StatusNotFound, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}

	errHead := "Update Article Failed"
	errBody := "Please try again."

	newArticle, invalids, err := handleForm(c)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody, "errTags": invalids})
		return
	} else if err != nil {
		errBody = err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}
	newArticle.ID = id

	id, res := databases.SubmitArticle(newArticle, "update")
	if !res {
		errBody = "An error occurred while writing to DB."
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}
	logrus.Infof("Update article with id %v", id)
	c.Header("Location", "/articles/browse?articleId="+strconv.Itoa(id))
	c.JSON(http.StatusCreated, gin.H{"articleId": id})
}

func DeleteArticle(c *gin.Context) {
	errHead := "Delete Article Failed"
	errBody := "Please try again."

	id := getParaId(c, "articleId")
	if id == 0 {
		errHead = "Article ID is An Positive Integer"
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	if res := databases.DeleteArticle(id, true); !res {
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}
	logrus.Infof("Delete article with id %v", id)
	c.Status(http.StatusNoContent)
}
