package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func CreateArticleView(c *gin.Context) {
	c.HTML(http.StatusOK, "editor.html", gin.H{
		"currPageCSS": "css/editor.css",
		"title":       "Create New Post",
		"function":    "create",
	})
}

func UpdateArticleView(c *gin.Context) {
	id := checkArticleId(c, "articleId")
	if id == 0 {
		errHead := "Article ID is An Integer"
		errBody := "Please try again."
		c.HTML(http.StatusBadRequest, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"errHead":     errHead,
			"errBody":     errBody,
		})
		return
	}

	if succeed := databases.IsArticleExists(id); succeed != true {
		errHead := "Article Not Found"
		errBody := "Please try again."
		c.HTML(http.StatusNotFound, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"errHead":     errHead,
			"errBody":     errBody,
		})
		return
	}

	dbFormatArticle, succeed := databases.GetArticleFullContent(id)
	if succeed != true {
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
	c.HTML(http.StatusOK, "editor.html", gin.H{
		"function":     "update",
		"currPageCSS":  "css/editor.css",
		"title":        "Edit: " + article.Title,
		"articleTitle": article.Title,
		"subtitle":     article.Subtitle,
		"date":         article.Date,
		"author":       article.Authors,
		"category":     article.Category,
		"tags":         article.Tags,
		"content":      article.Content,
	})
}

func CreateArticle(c *gin.Context) {
	errHead := "An Error Occurred"
	errBody := "Please try again."

	newArticle, err := handleForm(c)
	if err != nil {
		errHead := "An Error Occurred"
		errBody := err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	invalids := validateArticleFormat(newArticle)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody, "errTags": invalids})
		return
	}

	newArticle.Tags = removeDuplicateTags(newArticle.Tags)
	dbFormatArticle := articleFormatDetailedToDB(newArticle)
	id, res := databases.InsertArticle(dbFormatArticle)
	if !res {
		errBody = "An error occurred while writing to DB."
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}
	logrus.Infof("Create article with id %v\n", id)
	c.Header("Location", "/articles/browse?articleId="+strconv.Itoa(id)) // With Location header and status code 3XX (but not 2XX), response.redirected becomes true
	c.JSON(http.StatusCreated, gin.H{"articleId": id})
}

func UpdateArticle(c *gin.Context) {
	id := checkArticleId(c, "articleId")
	if id == 0 {
		errHead := "Article ID is An Integer"
		errBody := "Please try again."
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}

	if succeed := databases.IsArticleExists(id); succeed != true {
		errHead := "Article Not Found"
		errBody := "Please try again."
		c.JSON(http.StatusNotFound, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}

	errHead := "An Error Occurred"
	errBody := "Please try again."

	newArticle, err := handleForm(c)
	if err != nil {
		errHead := "An Error Occurred"
		errBody := err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	invalids := validateArticleFormat(newArticle)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody, "errTags": invalids})
		return
	}

	newArticle.Tags = removeDuplicateTags(newArticle.Tags)
	dbFormatArticle := articleFormatDetailedToDB(newArticle)
	dbFormatArticle.ID = id

	id, res := databases.ReplaceArticle(dbFormatArticle)
	if !res {
		errHead := "An Error Occurred"
		errBody := "An error occurred while writing to DB."
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}
	logrus.Infof("Update article with id %v\n", id)
	c.Header("Location", "/articles/browse?articleId="+strconv.Itoa(id))
	c.JSON(http.StatusCreated, gin.H{"articleId": id})
}

func DeleteArticle(c *gin.Context) {
	id := checkArticleId(c, "articleId")
	if id == 0 {
		errHead := "Article ID is An Integer"
		errBody := "Please try again."
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	if succeed := databases.IsArticleExists(id); succeed != true {
		errHead := "Article Not Found"
		errBody := "Please try again."
		c.JSON(http.StatusNotFound, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	if res := databases.DeleteArticle(id); !res {
		errHead := "An Error Occurred"
		errBody := "An error occurred while writing to DB."
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}
	logrus.Infof("Delete article with id %v\n", id)
	c.Status(http.StatusNoContent)
}
