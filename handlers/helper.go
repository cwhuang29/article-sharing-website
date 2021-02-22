package handlers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	OldestDate, _ = time.Parse("2006-01-02", "1960-01-01")
	TagsLmit      = 5
	TagsCharLmit  = 65 // TODO length of each emoji is about 13
	ErrInputMsg   = map[string]string{
		"empty":        "The value can't be empty.",
		"long":         "Too many words.",
		"dateTooOld":   "The date chosen should be greater than 1960-01-01.",
		"dateFuture":   "The date chosen can't be in the future.",
		"tagsTooMany":  "You can target up to 5 tags at a time.",
		"tagsTooLong":  "Each tag can contain at most 20 charaters.",
		"emailInvalid": "The email format is not correct.",
	}
	overviewContentLength = 660
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
	if len(article.Content) > overviewContentLength { // This might cut on the falh of html tag
		//article.Content = article.Content[:overviewContentLength]
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
	dt.Content = article.Content
	return
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

func validateCreateArticle(newArticle Article) (err map[string]interface{}) {
	err = make(map[string]interface{})

	// fmt.Println(newArticle.Date, time.Now().Format("2006-01-02"), OldestDate, OldestDate.String(), OldestDate.Local().String())
	// 2020-01-01 2021-02-15 1960-01-01 00:00:00 +0000 UTC 1960-01-01 00:00:00 +0000 UTC 1960-01-01 08:00:00 +0800 CST

	if len(newArticle.Title) == 0 {
		err["title"] = ErrInputMsg["short"]
	} else if len(strings.Split(newArticle.Title, " ")) > 20 {
		err["title"] = ErrInputMsg["long"]
	}

	if len(strings.Split(newArticle.Subtitle, " ")) > 20 { // Subtitle can be empty
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
