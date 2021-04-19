package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func About(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{"currPageCSS": "css/about.css"}) // Call the HTML method of the Context to render a template
}

func Home(c *gin.Context) {
	userStatus, user := GetUserStatus(c)
	if userStatus < IsMember {
		errHead := "Login to view home page"
		errBody := ""
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	ttlBookmarks := databases.CountUserBookmarks(user.ID)

	c.HTML(http.StatusOK, "home.html", gin.H{
		"currPageCSS":    "css/overview.css", // About 95% of code is copied from overview.html and overview.js
		"title":          user.FirstName + " " + user.LastName,
		"totalBookmarks": ttlBookmarks,
	})
}

func Overview(c *gin.Context) {
	// If frontend trigger this route via window.location.href="/articles/browse?articleId=1", then c.FullPath() is /articles/browse
	title := ""
	switch c.FullPath() {
	case "/articles/weekly-update":
		title = "Weekly News"
	case "/articles/medication":
		title = "Medication Related News"
	case "/articles/pharma":
		title = "Pharma News"
	}

	c.HTML(http.StatusOK, "overview.html", gin.H{
		"currPageCSS": "css/overview.css",
		"title":       title,
	})
}

func SearchTags(c *gin.Context) {
	tag := c.Query("query") // Request via "/articles/tags?query=<value>"
	if tag == "" {
		c.Redirect(http.StatusFound, "/articles/weekly-update")
	}

	databases.UpdateTagsStats(tag)
	c.HTML(http.StatusOK, "overview.html", gin.H{
		"currPageCSS": "css/overview.css",
		"title":       "Results for: " + tag,
	})
}

func CheckPermissionAndArticleExists(c *gin.Context) {
	id := getParaId(c, "articleId")
	if id == 0 {
		errHead := "Article ID is An Positive Integer"
		errBody := "Please try again."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	if succeed := databases.IsArticleExists(id, true); succeed != true {
		errHead := "Article ID Not Found"
		errBody := "Please try again."
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	c.Status(http.StatusOK)
}

func FetchData(c *gin.Context) {
	errHead := "Invalid Parameter"

	types := c.DefaultQuery("type", "")
	if types == "" {
		errBody := "Parameter type can not be empty."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody, "size": 0})
		return
	}

	query := c.DefaultQuery("query", "")
	if query == "" {
		errBody := "Parameter query can not be empty."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody, "size": 0})
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 0 {
		errBody := "Parameter offset should be a non-negative integer."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody, "size": 0})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		errBody := "Parameter limit should be a positive integer."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody, "size": 0})
		return
	}

	isAdmin := isUserAdmin(c)
	data, err := fetchData(types, query, offset, limit, isAdmin)
	c.JSON(http.StatusOK, gin.H{"articleList": data, "size": len(data)}) // Notice: if the data is an empty array [], frontend will get `null` instead of an empty array
}

func Browse(c *gin.Context) {
	if c.Query("articleId") == "" {
		c.Redirect(http.StatusFound, "weekly-update")
		return
	}

	id, err := strconv.Atoi(c.Query("articleId"))
	if err != nil || id <= 0 {
		errHead := "Article ID is An Integer"
		errBody := "Go back to previous page and try again."
		c.HTML(http.StatusBadRequest, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"errHead":     errHead,
			"errBody":     errBody,
		})
		return
	}

	isAdmin := isUserAdmin(c)
	dbFormatArticle := databases.GetArticle(id, isAdmin)
	if dbFormatArticle.ID == 0 {
		errHead := "Article Not Found"
		errBody := "Go back to previous page and try again."
		c.HTML(http.StatusNotFound, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"errHead":     errHead,
			"errBody":     errBody,
		})
		return
	}

	article := articleFormatDBToDetailed(dbFormatArticle, true)

	uuid := ""
	cookieEmail, _ := c.Cookie("login_email")
	adminEmail, _ := c.Cookie("is_admin")
	if adminEmail != "" && cookieEmail == adminEmail && databases.IsAdminUser(adminEmail) {
		uuid = utils.GetUUID()
		c.SetCookie("csrf_token", uuid, utils.CsrfTokenAge, "/", "", true, true)
	}

	c.HTML(http.StatusOK, "browse.html", gin.H{
		"currPageCSS": "css/browse.css",
		"csrfToken":   uuid,
		"success":     true,
		"adminOnly":   article.AdminOnly,
		"title":       article.Title,
		"subtitle":    article.Subtitle,
		"date":        article.Date,
		"author":      article.Authors,
		"category":    article.Category,
		"tags":        article.Tags,
		"content":     article.Content,
	})
	/*
		c.HTML(http.StatusOK, "browse.html", gin.H{
			"currPageCSS":      "css/browse.css",
			"title":             "This is Title",
			"subtitle":          "-----  This is subtitle  -----",
			"content":           "123456789",
			"oneImageSrc":      "/assets/img_640_480.png",
			"twoImagesSrc01":   "/assets/img_360_270.png",
			"twoImagesSrc02":   "/assets/img_360_270.png",
			"threeImagesSrc01": "/assets/img_360_270.png",
			"threeImagesSrc02": "/assets/img_360_270.png",
			"threeImagesSrc03": "/assets/img_360_270.png",
		})
	*/
}

func ContactUs(c *gin.Context) {
}

func AdminOverview(c *gin.Context) {
}
