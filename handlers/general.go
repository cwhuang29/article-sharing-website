package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func About(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{"currPageCSS": "css/about.css"}) // Call the HTML method of the Context to render a template
}

func CheckPermissionAndArticleExists(c *gin.Context) {
	id := getParaArticleId(c, "articleId")
	if id == 0 {
		errHead := "Article ID is An Positive Integer"
		errBody := "Please try again."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	if succeed := databases.IsArticleExists(id); succeed != true {
		errHead := "Article ID Not Found"
		errBody := "Please try again."
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	c.Status(http.StatusOK)
}

func WeeklyUpdate(c *gin.Context) {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	sevenDaysAgo := today.AddDate(0, 0, -7)

	dbFormatArticle := databases.GetArticlesInATimePeriod(sevenDaysAgo, today.AddDate(0, 0, 1))
	if len(dbFormatArticle) == 0 {
		c.HTML(http.StatusOK, "overview.html", gin.H{
			"currPageCSS": "css/overview.css",
			"title":       "Weekly News",
			"errHead":     "Oops ...",
			"errBody":     "No new articles in the past 7 days.",
		})
		return
	}

	var articleList []OverviewArticle

	for _, a := range dbFormatArticle {
		articleList = append(articleList, articleFormatDBToOverview(a))
	}
	c.HTML(http.StatusOK, "overview.html", gin.H{
		"currPageCSS":  "css/overview.css",
		"title":        "Weekly News",
		"articleList":  articleList,
		"initialCount": len(articleList),
	})
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

	data, err := fetchData(types, query, offset, limit)
	c.JSON(http.StatusOK, gin.H{"articleList": data, "size": len(data)}) // Notice: if the data is an empty array [], frontend will get `null` instead of empty array
}

func Overview(c *gin.Context) {
	// If frontend trigger this route via window.location.href="/articles/browse?articleId=1", then c.FullPath() is /articles/browse
	title := "Pharma News"
	if c.FullPath() == "/articles/medication" {
		title = "Medication Related News"
	}

	c.HTML(http.StatusOK, "overview.html", gin.H{
		"currPageCSS": "css/overview.css",
		"title":       title,
	})
}

func SearchTags(c *gin.Context) {
	tag := getParaTagValue(c, "query")
	if tag == "" {
		c.Redirect(http.StatusFound, "/articles/weekly-update")
	}

	updateTagsStats(tag)
	c.HTML(http.StatusOK, "overview.html", gin.H{
		"currPageCSS": "css/overview.css",
		"title":       "Results for: " + tag,
	})
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

	dbFormatArticle := databases.GetArticle(id)
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
	if cookieEmail == adminEmail && databases.IsAdminUser(adminEmail) {
		uuid = getUUID()
		c.SetCookie("csrf_token", uuid, csrfTokenAge, "/", "", true, true)
	}

	c.HTML(http.StatusOK, "browse.html", gin.H{
		"currPageCSS": "css/browse.css",
		"csrfToken":   uuid,
		"success":     true,
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
