package handlers

import (
	"fmt"
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func About(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{"currPageCSS": "css/about.css"}) // Call the HTML method of the Context to render a template
}

func CheckPermission(c *gin.Context) {
	yes, _ := isLoginedAdmin(c)
	if !yes {
		errMsg := "<div><strong>You are not allowed to perform this action</strong><p>Login if you are administrator.</p></div>"
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": errMsg})
		return
	}

	id := checkArticleId(c, "articleId")
	if id == 0 {
		errMsg := "<div><p><strong>Article ID is an integer</strong></p><p>Please try again.</p></div>"
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": errMsg})
		return
	}

	if succeed := databases.IsArticleExists(id); succeed != true {
		errMsg := "<div><p><strong>Article ID Not Found</strong></p><p>Please try again.</p></div>"
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"err": errMsg})
		return
	}

	c.Status(http.StatusOK)
}

func WeeklyUpdate(c *gin.Context) {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	sevenDaysAgo := today.AddDate(0, 0, -7)

	dbFormatArticle := databases.GetArticlesInNDays(sevenDaysAgo)
	if len(dbFormatArticle) == 0 {
		c.HTML(http.StatusOK, "overview.html", gin.H{
			"currPageCSS": "css/overview.css",
			"title":       "Weekly News",
			"err":         "<strong>Oops ... </strong><br>No new articles in the past 7 days",
		})
	} else {
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
}

func FetchData(c *gin.Context) {
	category := c.DefaultQuery("category", "")
	if category == "" {
		errMsg := "Parameter category can not be empty."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": errMsg, "size": 0})
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		errMsg := "Parameter offset should be a positive integer."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": errMsg, "size": 0})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		errMsg := "Parameter limit should be a positive integer."
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": errMsg, "size": 0})
		return
	}

	data, err := fetchData(category, offset, limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error(), "size": 0})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articleList": data, "size": len(data)}) // Notice: if the data is an empty array [], frontend will get `null` instead of empty array
}

func Overview(c *gin.Context) {
	// fmt.Println(c.FullPath()) // If frontend trigger this route via window.location.href="/articles/browse?articleId=1", then c.FullPath() is /articles/browse
	title := "Pharma News"
	if c.FullPath() == "/articles/medication" {
		title = "Medication Related News"
	}

	c.HTML(http.StatusOK, "overview.html", gin.H{
		"currPageCSS": "css/overview.css",
		"title":       title,
	})
}

func Browse(c *gin.Context) {
	if c.Query("articleId") == "" {
		c.Redirect(http.StatusFound, "weekly-update")
		return
	}

	id, err := strconv.Atoi(c.Query("articleId"))
	if err != nil || id <= 0 {
		err := fmt.Errorf("<div><p><strong>Article Not Found</strong></p><p>Go back to previous page and try again.</p></div>")
		// c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		c.HTML(http.StatusBadRequest, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"err":         err.Error(),
		})
		return
	}

	if dbFormatArticle, succeed := databases.GetArticleFullContent(id); succeed != true {
		errMsg := "<div><p><strong>Article ID Not Found</strong></p><p>Please try again.</p></div>"
		c.HTML(http.StatusNotFound, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"err":         errMsg,
		})
	} else {
		article := articleFormatDBToDetailed(dbFormatArticle, true)
		c.HTML(http.StatusOK, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"success":     true,
			"title":       article.Title,
			"subtitle":    article.Subtitle,
			"date":        article.Date,
			"author":      article.Authors,
			"category":    article.Category,
			"tags":        article.Tags,
			"content":     article.Content,
		})
	}
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
