package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"strings"
)

const (
	titleSizeLimit    = 34.
	subtitleSizeLimit = 54.
)

func articleFormatDBToOverview(article models.Article) (a Article) {
	a.ID = article.ID

	if len(article.Title) > titleSizeLimit {
		a.Title = decodeRunes(article.Title, titleSizeLimit) + " ..."
	} else {
		a.Title = article.Title
	}

	if len(article.Subtitle) > subtitleSizeLimit {
		a.Subtitle = decodeRunes(article.Subtitle, subtitleSizeLimit) + " ..."
	} else {
		a.Subtitle = article.Subtitle
	}

	a.Date = article.ReleaseDate.String()
	a.Authors = strings.Split(article.Authors, ",")
	a.Category = strings.ToLower(article.Category) // Because router only accepts lower case path

	a.Tags = []string{}
	for _, t := range article.Tags {
		a.Tags = append(a.Tags, t.Value)
	}

	truncate := false
	if len(article.Content) > overviewContentLength {
		truncate = true
	}
	a.Content = parseMarkdownToHTML(article.Content, truncate)
	a.AdminOnly = article.AdminOnly
	return
}

func articleFormatDBToDetailed(article models.Article, parseMarkdown bool) (a Article) {
	a.Title = article.Title
	a.Subtitle = article.Subtitle
	a.Date = article.ReleaseDate.Format("2006-01-02")
	a.Authors = strings.Split(article.Authors, ",")
	a.Category = strings.ToLower(article.Category)

	a.Tags = []string{} // Without initial, html template brokes (var tags = {{ .tags }};)
	for _, t := range article.Tags {
		a.Tags = append(a.Tags, t.Value)
	}

	if parseMarkdown {
		a.Content = parseMarkdownToHTML(article.Content, false)
	} else {
		a.Content = article.Content
	}

	a.AdminOnly = article.AdminOnly
	return
}
