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

func storeLoginToken(email string, loginMaxAge int) (token string) {
	token = getUUID()
	databases.InsertLoginToken(email, token, loginMaxAge)
	return
}

func clearLoginToken(email string, token string) {
	// Notice: Users may have multiple tokens based on different user agents they have logged in from, and those tokens must be removed from DB when expired
	// It can be done at login, logout, or any other time. Currently I'll done this job when user logout
	databases.DeleteLoginToken(email, token)
}

func RegisterView(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{"currPageCSS": "css/register.css", "title": "Signup"}) // Call the HTML method of the Context to render a template
}

func Register(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": true, "errHead": err.Error()})
		return
	}

	invalids := validateUserFormat(newUser)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": false, "errTags": invalids})
		return
	}

	tmp := databases.GetUser(newUser.Email)
	if tmp.ID != 0 {
		errHead := "This email is already registered"
		errBody := "Please use another email."
		c.JSON(http.StatusConflict, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}

	if databases.IsAdminUser(newUser.Email) {
		newUser.Admin = true
	}

	hashedPwd, err := hashPassword(newUser.Password)
	if err != nil {
		errHead := "Some Severe Errors Occurred"
		errBody := "Please reload the page and try again."
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}

	newUser.Password = string(hashedPwd)
	_, res := databases.InsertUserToDB(newUser)
	if !res {
		errHead := "Some Severe Errors Occurred"
		errBody := "Please reload the page and try again."
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}

	token := storeLoginToken(newUser.Email, loginMaxAge)
	c.Header("Location", landingPage)
	c.SetCookie("login_token", token, loginMaxAge, "/", domain, false, false)
	c.SetCookie("login_email", newUser.Email, loginMaxAge, "/", domain, false, false)
	if newUser.Admin {
		c.SetCookie("is_admin", newUser.Email, loginMaxAge, "/", domain, false, false)
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func LoginView(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"currPageCSS": "css/login.css", "title": "Login"})
}

func LoginJSON(c *gin.Context) {
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"inputFormatInvalid": true, "errHead": err.Error()})
		return
	}

	invalids := validateLoginFormat(json.Email, json.Password)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"inputFormatInvalid": false, "errTags": invalids})
		return
	}

	var user models.User
	user = databases.GetUser(json.Email)
	if user.ID == 0 {
		errHead := "User not Found"
		errBody := "Please try again."
		c.JSON(http.StatusForbidden, gin.H{"inputFormatInvalid": false, "errHead": errHead, "errBody": errBody})
		return
	}

	err := compareHashAndPassword([]byte(user.Password), []byte(json.Password))
	if err != nil {
		errHead := "Password Incorrect"
		errBody := "Please try again."
		c.JSON(http.StatusForbidden, gin.H{"inputFormatInvalid": false, "errHead": errHead, "errBody": errBody})
		return
	}

	token := storeLoginToken(user.Email, loginMaxAge)
	c.Header("Location", landingPage)
	c.SetCookie("login_token", token, loginMaxAge, "/", domain, false, true)
	c.SetCookie("login_email", user.Email, loginMaxAge, "/", domain, false, false)
	if user.Admin {
		c.SetCookie("is_admin", user.Email, loginMaxAge, "/", domain, false, false)
	}
	c.JSON(http.StatusOK, gin.H{})
}

func Logout(c *gin.Context) {
	token, err := c.Cookie("login_token")
	if err != nil || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	yes, email := isLoginedAdmin(c)
	if yes {
		c.SetCookie("is_admin", "", loginMaxAge, "/", domain, false, true)
	}

	c.SetCookie("login_token", "", loginMaxAge, "/", domain, false, true)
	c.SetCookie("login_email", "", loginMaxAge, "/", domain, false, true)

	clearLoginToken(email, token)

	c.Header("Location", landingPage)
	c.JSON(http.StatusResetContent, gin.H{})
}

func PasswordResetRequest(c *gin.Context) {
}

func PasswordResetEmail(c *gin.Context) {
}

func PasswordResetView(c *gin.Context) {
}

func PasswordReset(c *gin.Context) {
}
