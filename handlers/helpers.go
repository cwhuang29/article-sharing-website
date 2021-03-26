package handlers

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/google/uuid"

	"github.com/russross/blackfriday"
	"golang.org/x/crypto/bcrypt"
)

var (
	overviewContentLength = 800
	emailRegex            = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func GetUserStatus(c *gin.Context) (status UserStatus, cookieEmail string) {
	cookieEmail, _ = c.Cookie("login_email") // If no such cookie, c.Cookie() returns empty string with error `named cookie not present`
	cookieToken, _ := c.Cookie("login_token")
	adminEmail, _ := c.Cookie("is_admin")

	memberOrAdmin := IsMember
	if adminEmail != "" && cookieEmail == adminEmail && databases.IsAdminUser(adminEmail) {
		memberOrAdmin = IsAdmin
	}

	creds := databases.GetLoginCredentials(cookieEmail)
	for i := 0; i < len(creds); i++ {
		isEpr := isExpired(creds[i].LastLogin, creds[i].MaxAge)
		if cookieEmail == creds[i].User.Email && cookieToken == creds[i].Token && !isEpr {
			status = memberOrAdmin
			return
		}
	}

	cookieEmail = ""
	return
}

// To reduce the number of queries to database, don't call GetUserStatus() here.
// If it is a critical (non-read) operation, there will be more rigorous validation in other functions.
func detectIfUserIsAdmin(c *gin.Context) bool {
	if admin, _ := c.Cookie("is_admin"); admin != "" {
		return true
	}
	return false
}

func isExpired(startTime time.Time, period int) bool {
	now := time.Now().UTC().Truncate(time.Second)
	if now.Sub(startTime).Seconds() > float64(period) {
		return true
	}
	return false
}

func getParaArticleId(c *gin.Context, key string) int {
	if c.Query(key) == "" {
		return 0
	}

	id, err := strconv.Atoi(c.Query(key))
	if err != nil || id <= 0 {
		return 0
	}

	return id
}

func getParaTagValue(c *gin.Context, key string) string {
	return c.Query(key)
}

func updateTagsStats(tag string) {
	databases.UpdateTagsStats(tag)
}

func fetchData(types, query string, offset, limit int, isAdmin bool) (articleList []Article, err error) {
	var dbFormatArticle []models.Article

	switch types {
	case "tag":
		dbFormatArticle = databases.GetSameTagArticles(query, offset, limit, isAdmin)
		for i := 0; i < len(dbFormatArticle); i++ {
			dbFormatArticle[i].Tags = databases.GetArticleTags(dbFormatArticle[i])
		}
	case "category":
		dbFormatArticle = databases.GetSameCategoryArticles(query, offset, limit, isAdmin)
	}

	for _, a := range dbFormatArticle {
		articleList = append(articleList, articleFormatDBToOverview(a))
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
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func compareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
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
