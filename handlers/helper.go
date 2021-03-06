package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cwhuang29/article-sharing-website/databases"
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
		"empty":            "The field can't be empty.",
		"long":             "This field can have no more than 255 characters.",
		"dateTooOld":       "The date chosen should be greater than 1960-01-01.",
		"dateFuture":       "The date chosen can't be in the future.",
		"tagsTooMany":      "You can target up to 5 tags at a time.",
		"tagsTooLong":      "Each tag can contain at most 20 charaters.",
		"emailInvalid":     "The email format is not correct.",
		"passwordTooShort": "Passwords must be at least 8 characters long.",
	}
	overviewContentLength = 800
	emailRegex            = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func isAdmin(c *gin.Context) bool {
	adminEmail, err := c.Cookie("is_admin")
	if err != nil {
		return false
	}

	if !databases.IsAdminUser(adminEmail) {
		return false
	}

	return isLogin(c)
}

func isLogin(c *gin.Context) bool {
	cookieEmail, err := c.Cookie("login_email")
	if err != nil {
		return false
	}

	cookieToken, err := c.Cookie("login_token")
	if err != nil {
		return false
	}

	loginCredentials := databases.GetLoginCredentials(cookieEmail)
	if cookieEmail != loginCredentials.Email || cookieToken != loginCredentials.Token || isExpired(loginCredentials.LastLogin, loginCredentials.MaxAge) {
		return false
	}

	return true
}

func isExpired(startTime time.Time, period int) bool {
	now := time.Now().UTC().Truncate(time.Second)
	if now.Sub(startTime).Seconds() > float64(period) {
		return true
	}
	return false
}

func checkArticleId(c *gin.Context, key string) int {
	if c.Query(key) == "" {
		return 0
	}

	id, err := strconv.Atoi(c.Query(key))
	if err != nil || id <= 0 {
		return 0
	}

	return id
}

func fetchData(category string, offset int, limit int) (articleList []OverviewArticle, err error) {
	if offset < 0 {
		err = fmt.Errorf("Invalid parameter: offset should not be negative.")
		return
	} else if limit <= 0 {
		err = fmt.Errorf("Invalid parameter: limit should not be negative.")
		return
	}

	dbFormatArticle := databases.GetArticlesList(category, offset, limit)
	for _, a := range dbFormatArticle {
		articleList = append(articleList, articleFormatDBToOverview(a))
	}
	return
}

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
	return string(byteS)
}

func getUUID() string {
	return uuid.NewString()
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("<div><p><strong>Some Severe Errors Occurred</strong></p><p>Please reload the page and try again.</p></div>")
	}
	return hashedPassword, err
}

func compareHashAndPassword(hashedPassword, password []byte) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		err = fmt.Errorf("<div><p><strong>Password Incorrect</strong></p><p>Please try again.</p></div>")
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

		keys := make([]interface{}, 0)
		for key := range unq {
			keys = append(keys, key)
		}
		return keys
	default:
		return nil
	}
}

func validateUserFormat(newUser models.User) (err map[string]interface{}) {
	err = make(map[string]interface{})

	if len(newUser.FirstName) == 0 {
		err["first_name"] = ErrInputMsg["empty"]
	}

	if len(newUser.LastName) == 0 {
		err["last_name"] = ErrInputMsg["empty"]
	}

	if len(newUser.Password) == 0 {
		err["password"] = ErrInputMsg["empty"]
	} else if len(newUser.Password) < 8 {
		err["password"] = ErrInputMsg["passwordTooShort"]
	}

	if len(newUser.Email) == 0 {
		err["email"] = ErrInputMsg["empty"]
	}

	if len(newUser.Gender) == 0 {
		err["gender"] = ErrInputMsg["empty"]
	}

	if len(newUser.Major) == 0 {
		err["major"] = ErrInputMsg["empty"]
	}

	return err
}

func validateArticleFormat(newArticle Article) (err map[string]interface{}) {
	err = make(map[string]interface{})

	// fmt.Println(newArticle.Date, time.Now().UTC().Format("2006-01-02"), OldestDate, OldestDate.String(), OldestDate.Local().String())
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
		// if time.Now().Truncate(time.Hour * 24).Sub(inpDate) < 0 {
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
