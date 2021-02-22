package handlers

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/google/uuid"

	"github.com/russross/blackfriday"
	"golang.org/x/crypto/bcrypt"
)

var (
	OldestDate, _ = time.Parse("2006-01-02", "1960-01-01")
	TagsLmit      = 5
	TagsCharLmit  = 65 // TODO length of each emoji is about 13
	ErrInputMsg   = map[string]string{
		"empty":        "The field can't be empty.",
		"long":         "This field can have no more than 255 characters.",
		"dateTooOld":   "The date chosen should be greater than 1960-01-01.",
		"dateFuture":   "The date chosen can't be in the future.",
		"tagsTooMany":  "You can target up to 5 tags at a time.",
		"tagsTooLong":  "Each tag can contain at most 20 charaters.",
		"emailInvalid": "The email format is not correct.",
	}
	overviewContentLength = 800
	emailRegex            = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

/*
 * There are 3 kinds of format:
 * 1. Detailed format (createing, updating ...)
 * 2. Overview format
 * 3. DB format
 */

func ArticleFormatDetailedToDB(article Article) (db models.Article) {
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

func ArticleFormatDBToOverview(article models.Article) (ov OverviewArticle) {
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
	if len(article.Content) > overviewContentLength {
		article.Content = parseMarkdownToHTML(article.Content, true)
	}
	// ov.Content = strings.ReplaceAll(article.Content, "\n", "<br>") // Done in frontend cause the respond will be escaped
	ov.Content = article.Content
	return
}

func ArticleFormatDBToDetailed(article models.Article) (dt Article) {
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
	dt.Content = parseMarkdownToHTML(article.Content, false)
	return
}

func parseMarkdownToHTML(s string, truncate bool) string {
	/*
		It is such a bad idea to self-implement markdown parser
		images := regexp.MustCompile(`!\[([^\s]+)\]\(([^\s]+)\)`)
		links := regexp.MustCompile(`\[([^\s]+)\]\(([^\s]+)\)`)
		strikes := regexp.MustCompile(`~~(\w.*\w)~~`)
		bold := regexp.MustCompile(`\*\*(\w.*\w)\*\*`)
		italic := regexp.MustCompile(`\*(\w.*\w)\*`)
		code := regexp.MustCompile("`([^\r|\n]*)`")
		s = images.ReplaceAllString(s, `<figure class="image is-16by9"><img alt="$1" href="$2"></figure>`)
		s = links.ReplaceAllString(s, `<a href="$2">$1</a>`)
		s = strikes.ReplaceAllString(s, `<del>$1</del>`)
		s = bold.ReplaceAllString(s, `<strong>$1</strong>`)
		s = italic.ReplaceAllString(s, `<em>$1</em>`)
		s = code.ReplaceAllString(s, `<code>$1</code>`)
	*/
	if truncate {
		trunc := overviewContentLength
		if len(s) < overviewContentLength {
			trunc = len(s)
		}
		s = s[:trunc]
	}
	byteS := blackfriday.MarkdownCommon([]byte(s))
	fmt.Println(string(byteS))
	return string(byteS)
}

func getUUID() string {
	return uuid.NewString()
}

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("<div><p><strong>Some Severe Errors Occurred</strong></p><p>Please reload the page and try again</p></div>")
	}
	return hashedPassword, err
}

func compareHashAndPassword(hashedPassword, password []byte) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		err = fmt.Errorf("<div><p><strong>Password Incorrect</strong></p><p>Please try again</p></div>")
	}
	return err
}

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func validateLoginFormat(email string, password string) (err map[string]interface{}) {
	err = make(map[string]interface{})

	if len(email) == 0 {
		err["email"] = ErrInputMsg["empty"]
	} else if !isEmailValid(email) {
		err["email"] = ErrInputMsg["emailInvalid"]
	}

	if len(password) == 0 {
		err["password"] = ErrInputMsg["empty"]
	}
	return
}

func removeDuplicateTags(t []string) []string {
	tmp := removeDuplicateValuesInSlice(t)
	tags := make([]string, len(tmp))
	for i, v := range tmp {
		tags[i] = v.(string)
	}
	return tags
}

func removeDuplicateValuesInSlice(t interface{}) []interface{} {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)
		unq := make(map[interface{}]bool)

		for i := 0; i < s.Len(); i++ {
			if _, ok := unq[s.Index(i).Interface()]; !ok {
				unq[s.Index(i).Interface()] = true
			}
		}

		switch reflect.TypeOf(t).String() {

		}
		keys := make([]interface{}, 0)
		for key := range unq {
			keys = append(keys, key)
		}
		return keys
	default:
		return nil
	}
}

func validateCreateArticle(newArticle Article) (err map[string]interface{}) {
	err = make(map[string]interface{})

	// fmt.Println(newArticle.Date, time.Now().Format("2006-01-02"), OldestDate, OldestDate.String(), OldestDate.Local().String())
	// 2020-01-01 2021-02-15 1960-01-01 00:00:00 +0000 UTC 1960-01-01 00:00:00 +0000 UTC 1960-01-01 08:00:00 +0800 CST

	if len(newArticle.Title) == 0 {
		err["title"] = ErrInputMsg["short"]
	} else if len(newArticle.Title) > 255 {
		err["title"] = ErrInputMsg["long"]
	}

	if len(newArticle.Subtitle) > 255 { // Subtitle can be empty
		err["subtitle"] = ErrInputMsg["long"]
	}

	if inpDate, dateErr := time.Parse("2006-01-02", newArticle.Date); dateErr != nil {
		err["date"] = dateErr.Error()
	} else {
		// if time.Now().Sub(inpDate) < 0 {
		//     err["date"] = ErrInputMsg["dateFuture"]
		// }
		if OldestDate.Sub(inpDate) > 0 {
			err["date"] = ErrInputMsg["dateTooOld"]
		}
	}

	if len(newArticle.Tags) > TagsLmit {
		err["tags"] = ErrInputMsg["tagsTooMany"]
	} else {
		for _, t := range newArticle.Tags {
			if len(t) > TagsCharLmit {
				err["tags"] = ErrInputMsg["tagsTooLong"]
				break
			}
		}
	}

	if len(newArticle.Content) == 0 {
		err["content"] = ErrInputMsg["empty"]
	}
	return
}
