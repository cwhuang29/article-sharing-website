package handlers

import (
	"fmt"
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// router.GET("/welcome", func(c *gin.Context) { // The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
// 	firstname := c.DefaultQuery("firstname", "Guest")
// 	lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
// 	c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
// })

func About(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{"currPageCSS": "css/about.css"}) // Call the HTML method of the Context to render a template
}

func WeeklyUpdate(c *gin.Context) {
	today := time.Now().Truncate(24 * time.Hour)
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
			articleList = append(articleList, ArticleFormatDBToOverview(a))
		}
		c.HTML(http.StatusOK, "overview.html", gin.H{
			"currPageCSS": "css/overview.css",
			"title":       "Weekly News",
			"articleList": articleList,
		})
	}
}

func Overview(c *gin.Context) {
	offset := 0
	limit := 10
	category := "pharma"
	title := "Pharma News"
	if c.FullPath() == "/articles/medication" {
		category = "medication"
		title = "Medication Related News"
	}

	if dbFormatArticle, err := databases.GetArticlesList(category, offset, limit); err != nil {
		c.HTML(http.StatusBadRequest, "overview.html", gin.H{
			"currPageCSS": "css/overview.css",
			"title":       title,
			"err":         "<strong>Oops ... </strong><br>Please reload the page and try again",
		})
	} else {
		if len(dbFormatArticle) == 0 {
			c.HTML(http.StatusOK, "overview.html", gin.H{
				"currPageCSS": "css/overview.css",
				"title":       title,
				"err":         "<strong>Oops ... </strong><br>No new articles in the past 7 days",
			})
			return
		}
		var articleList []OverviewArticle
		for _, a := range dbFormatArticle {
			articleList = append(articleList, ArticleFormatDBToOverview(a))
		}
		c.HTML(http.StatusOK, "overview.html", gin.H{
			"currPageCSS": "css/overview.css",
			"title":       title,
			"articleList": articleList,
		})
	}
}

func Browse(c *gin.Context) {
	if c.Query("articleId") == "" {
		c.Redirect(http.StatusFound, "weekly-update")
		return
	}

	// fmt.Println(c.FullPath()) // Frontend trigger this route via window.location.href="/articles/browse?articleId=1" -> c.FullPath() is /articles/browse
	id, err := strconv.Atoi(c.Query("articleId"))
	if err != nil || id <= 0 {
		err := fmt.Errorf("<div><p><strong>Article Not Found</strong></p><p>Go back to previous page and try again</p></div>")
		// c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		c.HTML(http.StatusBadRequest, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"err":         err.Error(),
		})
		return
	}

	if dbFormatArticle, err := databases.GetArticleFullContent(id); err != nil {
		c.HTML(http.StatusNotFound, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"err":         err.Error(),
		})
	} else {
		article := ArticleFormatDBToDetailed(dbFormatArticle)
		// c.Header("Location", "/about") // With Location header, response.redirected will become true
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
