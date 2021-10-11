package handlers

import (
	"fmt"
	"net/http"

	"github.com/cwhuang29/article-sharing-website/constants"
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/utils"
	"github.com/gin-gonic/gin"
)

func About(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{"currPageCSS": "css/about.css"}) // Call the HTML method of the Context to render a template
}

func Home(c *gin.Context) {
	userStatus, user := GetUserStatus(c)
	if userStatus < IsMember {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errHead": fmt.Sprintf(constants.LoginTo, "view home page"), "errBody": ""})
		return
	}

	ttlBookmarks := databases.CountUserBookmarks(user.ID)

	c.HTML(http.StatusOK, "home.html", gin.H{
		"currPageCSS":    "css/overview.css",
		"title":          user.FirstName + " " + user.LastName,
		"totalBookmarks": ttlBookmarks,
	})
}

func Overview(c *gin.Context) {
	// If frontend trigger this route via window.location.href="/articles/browse?articleId=1", then c.FullPath() is /articles/browse
	title := ""
	switch c.FullPath() {
	case constants.URLLandingPage:
		title = "Weekly News"
	case constants.URLTopicMed:
		title = "Medication Related News"
	case constants.URLTopicPharma:
		title = "Pharma News"
	}

	c.HTML(http.StatusOK, "overview.html", gin.H{
		"currPageCSS": "css/overview.css",
		"title":       title,
	})
}

func SearchTags(c *gin.Context) {
	tag := getQueryPara(c, constants.QueryTagSearch) // Request via "/articles/tags?q=<value>"
	if tag == "" {
		c.Redirect(http.StatusFound, constants.URLLandingPage)
	}

	databases.UpdateTagsStats(tag)
	c.HTML(http.StatusOK, "overview.html", gin.H{
		"currPageCSS": "css/overview.css",
		"title":       "Results for: " + tag,
	})
}

func CheckPermissionAndArticleExists(c *gin.Context) {
	id, err := getQueryArticleID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.QueryArticleIDErr, "errBody": constants.TryAgain})
		return
	}

	if succeed := databases.IsArticleExists(id, true); !succeed {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"errHead": constants.ArticleNotFound, "errBody": constants.TryAgain})
		return
	}

	c.Status(http.StatusOK)
}

func FetchData(c *gin.Context) {
	types := getQueryPara(c, "type")
	if types == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.QueryErr, "errBody": fmt.Sprintf(constants.QueryEmptyErr, "type"), "size": 0})
		return
	}

	query := getQueryPara(c, "query")
	if query == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.QueryErr, "errBody": fmt.Sprintf(constants.QueryEmptyErr, "query"), "size": 0})
		return
	}

	offset, err := getQueryOffset(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.QueryErr, "errBody": constants.QueryOffsetErr, "size": 0})
		return
	}

	limit, err := getQueryLimit(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errHead": constants.QueryErr, "errBody": constants.QueryLimitErr, "size": 0})
		return
	}

	isAdmin := isUserAdmin(c)
	data, err := fetchData(types, query, offset, limit, isAdmin)
	c.JSON(http.StatusOK, gin.H{"articleList": data, "size": len(data)}) // Notice: if the data is an empty array [], frontend will get `null` instead of an empty array
}

func Browse(c *gin.Context) {
	if getQueryPara(c, constants.QueryArticleID) == "" {
		c.Redirect(http.StatusFound, "weekly-update")
		return
	}

	id, err := getQueryArticleID(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"errHead":     constants.QueryArticleIDErr,
			"errBody":     constants.GobackAndRetry,
		})
		return
	}

	isAdmin := isUserAdmin(c)
	dbFormatArticle := databases.GetArticle(id, isAdmin)
	if dbFormatArticle.ID == 0 {
		c.HTML(http.StatusNotFound, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"errHead":     constants.ArticleNotFound,
			"errBody":     constants.GobackAndRetry,
		})
		return
	}

	article := articleFormatDBToDetailed(dbFormatArticle, true)

	uuid := ""
	cookieEmail, _ := c.Cookie(constants.CookieLoginEmail)
	adminEmail, _ := c.Cookie(constants.CookieIsAdmin)
	if adminEmail != "" && cookieEmail == adminEmail && databases.IsAdminUser(adminEmail) {
		uuid = utils.GetUUID()
		c.SetCookie(constants.CookieCSRFToken, uuid, constants.CsrfTokenAge, "/", "", true, true)
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
