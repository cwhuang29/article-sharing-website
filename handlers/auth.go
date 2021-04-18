package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/cwhuang29/article-sharing-website/utils"
	"github.com/cwhuang29/article-sharing-website/utils/validator"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterView(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{"currPageCSS": "css/register.css", "title": "Signup"})
}

func LoginView(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"currPageCSS": "css/login.css", "title": "Login"})
}

func Register(c *gin.Context) {
	var newUser models.User
	errHead := "An Error Occurred"
	errBody := "Please reload the page and try again."

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": true, "errHead": err.Error()})
		return
	}

	invalids := validator.ValidateRegisterForm(newUser)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": false, "errTags": invalids})
		return
	}

	if tmp := databases.GetUser(newUser.Email); tmp.ID != 0 {
		errHead = "This email is already registered"
		errBody = "Please use another email."
		c.JSON(http.StatusConflict, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}

	if databases.IsAdminUser(newUser.Email) {
		newUser.Admin = true
	}

	hashedPwd, err := utils.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}

	newUser.Password = string(hashedPwd)
	id, res := databases.InsertUser(newUser)
	if !res {
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "errHead": errHead, "errBody": errBody})
		return
	}

	token := utils.StoreLoginToken(id, utils.LoginMaxAge)
	c.Header("Location", utils.LandingPage)
	c.SetCookie("login_token", token, utils.LoginMaxAge, "/", "", true, true)
	c.SetCookie("login_email", newUser.Email, utils.LoginMaxAge, "/", "", true, false) // Frontend relies on this cookie
	if newUser.Admin {
		c.SetCookie("is_admin", newUser.Email, utils.LoginMaxAge, "/", "", true, false) // Frontend relies on this cookie
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func Login(c *gin.Context) {
	json := struct {
		Email    string `form:"email" json:"email" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"inputFormatInvalid": true, "errHead": err.Error()})
		return
	}

	invalids := validator.ValidateLoginForm(json.Email, json.Password)
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

	err := utils.CompareHashAndPassword([]byte(user.Password), []byte(json.Password))
	if err != nil {
		errHead := "Password Incorrect"
		errBody := "Please try again."
		c.JSON(http.StatusForbidden, gin.H{"inputFormatInvalid": false, "errHead": errHead, "errBody": errBody})
		return
	}

	token := utils.StoreLoginToken(user.ID, utils.LoginMaxAge)
	c.Header("Location", utils.LandingPage)
	c.SetCookie("login_token", token, utils.LoginMaxAge, "/", "", true, true)
	c.SetCookie("login_email", user.Email, utils.LoginMaxAge, "/", "", true, false) // Frontend relies on this cookie
	if user.Admin {
		c.SetCookie("is_admin", user.Email, utils.LoginMaxAge, "/", "", true, false) // Frontend relies on this cookie
	}
	c.JSON(http.StatusOK, gin.H{})
}

func Logout(c *gin.Context) {
	email, _ := c.Cookie("login_email") // If no such cookie, first argument will be an empty string
	token, _ := c.Cookie("login_token")
	if email == "" {
		// We'll reach here if user logout in one tab and re-logout on the another tab subsequently
		// So don't regard this case as an error
		c.Header("Location", utils.LandingPage)
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	/*
	 * Notice: Users may have multiple tokens based on different user agents they have logged in from, and those
	 * tokens must be removed from DB when expired. For instance, the user has logged in from the cellphone and laptop.
	 * When the user logged out on the laptop, we'll check whether the login token for the cellphone expired
	 * It can be done at login, logout, or any other time. Currently, I'll do this task when the user logout
	 */
	utils.ClearLoginToken(token)
	user := databases.GetUser(email)
	utils.ClearExpiredLoginTokens(user.ID)

	c.SetCookie("login_token", "", 0, "/", "", true, true)
	c.SetCookie("login_email", "", 0, "/", "", true, true)
	c.SetCookie("is_admin", "", 0, "/", "", true, true)

	c.Header("Location", utils.LandingPage)
	c.JSON(http.StatusResetContent, gin.H{})
}
