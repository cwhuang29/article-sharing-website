package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"strings"
	"time"
)

/*
 * There are 3 kinds of format:
 * 1. Detailed format (createing, updating ...)
 * 2. Overview format
 * 3. DB format
 */

func articleFormatDetailedToDB(article Article) (db models.Article) {
	t, _ := time.Parse("2006-01-02", article.Date)

	db.Title = article.Title
	db.Subtitle = article.Subtitle
	db.ReleaseDate = t
	db.Author = strings.Join(article.Authors, ",")
	db.Category = strings.ToLower(article.Category)
	db.Tag = strings.Join(article.Tags, ",")
	db.Content = article.Content
	return
}

func articleFormatDBToOverview(article models.Article) (ov OverviewArticle) {
	ov.ID = article.ID
	ov.Title = article.Title
	ov.Subtitle = article.Subtitle
	ov.Date = article.ReleaseDate.String()
	ov.Authors = strings.Split(article.Author, ",")
	ov.Category = article.Category

	if article.Tag == "" {
		ov.Tags = []string{} // length is 0
	} else {
		ov.Tags = strings.Split(article.Tag, ",") // len(strings.Split("", ",")) is 1 (so html/template will show an empty element) !!!
	}

	toParse := false
	if len(article.Content) > overviewContentLength {
		toParse = true
	}
	article.Content = parseMarkdownToHTML(article.Content, toParse)
	ov.Content = article.Content
	return
}

func articleFormatDBToDetailed(article models.Article, parseMarkdown bool) (dt Article) {
	dt.Title = article.Title
	dt.Subtitle = article.Subtitle
	dt.Date = article.ReleaseDate.Format("2006-01-02")
	dt.Authors = strings.Split(article.Author, ",")
	dt.Category = article.Category
	if article.Tag == "" {
		dt.Tags = []string{} // length is 0
	} else {
		dt.Tags = strings.Split(article.Tag, ",") // len(strings.Split("", ",")) is 1 (so html/template will show an empty element) !!!
	}
	if parseMarkdown {
		dt.Content = parseMarkdownToHTML(article.Content, false)
	} else {
		dt.Content = article.Content
	}
	return
}
