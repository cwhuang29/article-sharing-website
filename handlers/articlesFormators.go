package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"strings"
)

const (
	// Give slightly more word counts since the CSS style (e.g. text-justify: inter-word) affects layout in varying degrees
	titleSizeLimit                 = 35.
	subtitleSizeLimit              = 54.
	outlineSizeLimit               = 510.
	outlineSizeLimitWithCoverPhoto = 340.
)

func articleFormatDBToOverview(article models.Article) (a Article) {
	a.ID = article.ID

	if len(article.Title) > titleSizeLimit {
		// Chinese words are about 1.8 times wider than English alphabets in title and subtitle
		a.Title = decodeRuneStringForFrontend(article.Title, titleSizeLimit, 1.78) + "&nbsp;..."
	} else {
		a.Title = article.Title
	}

	if len(article.Subtitle) > subtitleSizeLimit {
		a.Subtitle = decodeRuneStringForFrontend(article.Subtitle, subtitleSizeLimit, 1.78) + "&nbsp;..."
	} else {
		a.Subtitle = article.Subtitle
	}

	a.Date = article.ReleaseDate.String()
	a.Authors = strings.Split(article.Authors, ",")
	a.Category = strings.ToLower(article.Category) // Because router only accepts lower case path
	a.CoverPhoto = article.CoverPhoto              // The url of cover photo
	a.AdminOnly = article.AdminOnly

	a.Tags = []string{}
	for _, t := range article.Tags {
		a.Tags = append(a.Tags, t.Value)
	}

	if article.CoverPhoto != "" && len(article.Outline) > outlineSizeLimitWithCoverPhoto {
		a.Outline = decodeRuneStringForFrontend(article.Outline, outlineSizeLimitWithCoverPhoto, 2.15)
	} else if len(article.Outline) > outlineSizeLimit {
		a.Outline = decodeRuneStringForFrontend(article.Outline, outlineSizeLimit, 2.15)
	} else {
		a.Outline = article.Outline
	}

	return
}

func articleFormatDBToDetailed(article models.Article, parseMarkdown bool) (a Article) {
	a.Title = article.Title
	a.Subtitle = article.Subtitle
	a.Date = article.ReleaseDate.Format("2006-01-02")
	a.Authors = strings.Split(article.Authors, ",")
	a.Category = strings.ToLower(article.Category)
	a.Outline = article.Outline
	a.CoverPhoto = article.CoverPhoto
	a.AdminOnly = article.AdminOnly

	a.Tags = []string{} // Without initial, html template brokes (var tags = {{ .tags }};)
	for _, t := range article.Tags {
		a.Tags = append(a.Tags, t.Value)
	}

	if parseMarkdown {
		a.Content = parseMarkdownToHTML(article.Content)
	} else {
		a.Content = article.Content
	}

	return
}
