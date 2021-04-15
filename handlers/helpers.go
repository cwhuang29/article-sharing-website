package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/russross/blackfriday"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Each mandarin symbol takes 3 - 4 bytes.
// The following `limit` is a value to measured how many "characters" can be displayed,
// It measures not only number of characters, but also the width of each character.
// `ratio` is the width ratio of Chinese words versus English characters
func decodeRuneStringForFrontend(s string, limit float64, ratio float64) string {
	idx := 0

	for cnt := 0.; cnt < limit; {
		_, width := utf8.DecodeRuneInString(s[idx:])
		// fmt.Printf("%#U starts at byte position %d\n", runeValue, charWidth)

		idx += width
		if width == 1 {
			cnt += 1 // e.g. English alphabets
		} else {
			cnt += ratio
		}
	}
	return s[:idx]
}

func GetUserStatus(c *gin.Context) (status UserStatus, cookieEmail string) {
	cookieEmail, _ = c.Cookie("login_email") // If no such cookie, c.Cookie() returns empty string with error `named cookie not present`
	cookieToken, _ := c.Cookie("login_token")
	adminEmail, _ := c.Cookie("is_admin")

	memberOrAdmin := IsMember
	if adminEmail != "" && cookieEmail == adminEmail && databases.IsAdminUser(adminEmail) {
		memberOrAdmin = IsAdmin
	}

	user := databases.GetUser(cookieEmail)
	creds := databases.GetLoginCredentials(user.ID)
	for _, cred := range creds {
		isEpr := isExpired(cred.LastLogin, cred.MaxAge)
		if cookieEmail == cred.User.Email && cookieToken == cred.Token && !isEpr {
			status = memberOrAdmin
			return
		}
	}

	cookieEmail = ""
	return
}

func DetectIfUserIsAdmin(c *gin.Context) bool {
	status, _ := GetUserStatus(c)
	return status >= IsAdmin
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

func fetchData(types, query string, offset, limit int, isAdmin bool) (articleList []Article, err error) {
	var dbFormatArticle []models.Article

	switch types {
	case "time":
		// For first time, load the weekly articles (all articles in the latest 7 days)
		if offset == 0 {
			today := time.Now().UTC().Truncate(24 * time.Hour)
			sevenDaysAgo := today.AddDate(0, 0, -7)
			tomorrow := today.AddDate(0, 0, 1)
			dbFormatArticle = databases.GetArticlesInATimePeriod(sevenDaysAgo, tomorrow, isAdmin)
		} else {
			dbFormatArticle = databases.GetArticles(offset, limit, isAdmin)
		}
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

func parseMarkdownToHTML(s string) string {
	/*
		It is such a bad idea to self-implement markdown parser
		links := regexp.MustCompile(`\[([^\s]+)\]\(([^\s]+)\)`)
		code := regexp.MustCompile("`([^\r|\n]*)`")
		s = links.ReplaceAllString(s, `<a href="$2">$1</a>`)
		s = bold.ReplaceAllString(s, `<strong>$1</strong>`)
	*/
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
