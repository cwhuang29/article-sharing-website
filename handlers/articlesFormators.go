package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"strings"
)

var (
	titleLength    = 55
	subtitleLength = 90
)

func articleFormatDBToOverview(article models.Article) (ov OverviewArticle) {
	ov.ID = article.ID

	if len(article.Title) > titleLength {
		ov.Title = article.Title[:titleLength] + " ..."
	} else {
		ov.Title = article.Title
	}
	if len(article.Subtitle) > subtitleLength {
		ov.Subtitle = article.Subtitle[:subtitleLength] + " ..."
	} else {
		ov.Subtitle = article.Subtitle
	}
	ov.Date = article.ReleaseDate.String()
	ov.Authors = strings.Split(article.Authors, ",")
	ov.Category = strings.ToLower(article.Category) // Because router only accepts lower case path

	ov.Tags = []string{}
	for _, t := range article.Tags {
		ov.Tags = append(ov.Tags, t.Value)
	}

	truncate := false
	if len(article.Content) > overviewContentLength {
		truncate = true
	}
	ov.Content = parseMarkdownToHTML(article.Content, truncate)
	return
}

func articleFormatDBToDetailed(article models.Article, parseMarkdown bool) (dt Article) {
	dt.Title = article.Title
	dt.Subtitle = article.Subtitle
	dt.Date = article.ReleaseDate.Format("2006-01-02")
	dt.Authors = strings.Split(article.Authors, ",")
	dt.Category = strings.ToLower(article.Category)

	dt.Tags = []string{} // Without initial, html template brokes (var tags = {{ .tags }};)
	for _, t := range article.Tags {
		dt.Tags = append(dt.Tags, t.Value)
	}

	if parseMarkdown {
		dt.Content = parseMarkdownToHTML(article.Content, false)
	} else {
		dt.Content = article.Content
	}
	return
}
