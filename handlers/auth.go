package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	landingPage = "/articles/weekly-update"
	loginMaxAge = 30 * 86400 // 1 month
	domain      = ""
)

func RegisterView(c *gin.Context) {
}

func Register(c *gin.Context) {
}

func LoginView(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"currPageCSS": "css/login.css", "title": "Login"}) // Call the HTML method of the Context to render a template
}

func LoginJSON(c *gin.Context) {
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"inputFormatInvalid": true, "err": err.Error()})
		return
	}

	invalids := validateLoginFormat(json.Email, json.Password)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"inputFormatInvalid": false, "errTags": invalids})
		return
	}

	var user models.User
	user, err := databases.GetUser(json.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"inputFormatInvalid": false, "err": err.Error()})
		return
	}

	err = compareHashAndPassword([]byte(user.Password), []byte(json.Password))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"inputFormatInvalid": false, "err": err.Error()})
		return
	}

	token := getUUID()
	databases.InsertLoginToken(user.Email, token, loginMaxAge)
	c.Header("Location", landingPage)
	c.SetCookie("login_token", token, loginMaxAge, "/", domain, false, true)
	c.SetCookie("login_email", user.Email, loginMaxAge, "/", domain, false, false)
	c.JSON(http.StatusOK, gin.H{})
}

func Logout(c *gin.Context) {
}

func PasswordResetRequest(c *gin.Context) {
}

func PasswordResetEmail(c *gin.Context) {
}

func PasswordResetView(c *gin.Context) {
}

func PasswordReset(c *gin.Context) {
}
