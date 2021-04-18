package handlers

import (
	"github.com/cwhuang29/article-sharing-website/config"
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/cwhuang29/article-sharing-website/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func clearUserExpiredPasswordTokens(id int) {
	databases.DeleteExpiredPasswordTokens(id)
}

func doesUserHasEmailQuota(id int) bool {
	count := databases.CountUserResetPasswordTokens(id)
	return count < utils.ResetPasswordMaxRetry
}

func getPasswordResetToken(id, maxAge int) string {
	token := getUUID()

	if ok := databases.InsertResetPasswordToken(id, token, maxAge); !ok {
		return getPasswordResetToken(id, maxAge) // Try again if we got duplicate tokens
	}
	return token
}

func getPasswordResetTokenInstance(token string) models.Password {
	return databases.GetResetPasswordToken(token)
}

func PasswordResetRequest(c *gin.Context) {
	c.HTML(http.StatusOK, "passwordResetRequest.html", gin.H{"title": "Reset Password"})
}

func PasswordResetEmail(c *gin.Context) {
	var json = struct{ Email string }{} // Uppercase is required
	errHead := "Oops, this is unexpected"
	errBody := "Please reload the page and try again!"

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	user := databases.GetUser(json.Email)
	if user.ID == 0 {
		errHead = "Email Not Found"
		errBody = "Did you fill in the correct email address?"
		c.JSON(http.StatusUnauthorized, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	/*
	 * There should be a mechanism to remove expired tokens in Password table
	 * Currently I'll do it right here, and use go routine to clear expired tokens periodically in the future
	 */
	clearUserExpiredPasswordTokens(user.ID)

	if yes := doesUserHasEmailQuota(user.ID); !yes {
		errHead = "You are trying too often"
		errBody = "Please try again in one hour"
		c.JSON(http.StatusTooManyRequests, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	baseURL := config.GetConfig().App.Url
	token := getPasswordResetToken(user.ID, utils.ResetPasswordTokenMaxAge)
	link := baseURL + utils.ResetPasswordPath + token + "?email=" + user.Email
	expireMins := utils.ResetPasswordTokenMaxAge / 60
	name := user.FirstName + " " + user.LastName

	if ok := SendResetPasswordEmail(user.Email, name, link, expireMins); !ok {
		logrus.Error("[Failed] Password reset email. Sent to ", user.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	msgHead := "An email has been sent"
	msgBody := "Reset link has been sent to your email"
	logrus.Info("[Succeed] Password reset email. Sent to ", user.Email)
	c.JSON(http.StatusOK, gin.H{"msgHead": msgHead, "msgBody": msgBody})
}

func PasswordResetForm(c *gin.Context) {
	errHead := "Oops, something went wrong"
	errBody := "You may want to request a new reset password email."

	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	uuid := getUUID()
	c.SetCookie("csrf_token", uuid, utils.CsrfTokenAge, "/", "", true, true)
	c.HTML(http.StatusOK, "passwordResetForm.html", gin.H{
		"title":       "Reset Password",
		"csrfToken":   uuid,
		"currPageCSS": "", // Fields can't be omitted even not using
		"email":       email,
	})
}

func PasswordUpdate(c *gin.Context) {
	errHead := "Oops, something went wrong"
	errBody := "Please reopen the link from email."

	var json = struct {
		Email    string
		Password string
		Token    string
	}{} // Uppercase is required

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	email, password, token := json.Email, json.Password, json.Token
	if email == "" || password == "" || token == "" {
		errHead = "Some values are missing"
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	tokenObj := getPasswordResetTokenInstance(token)
	if email != tokenObj.User.Email {
		errBody = "Perhaps you didn't open the latest email."
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	if isExpired(tokenObj.CreatedAt, tokenObj.MaxAge) {
		databases.DeletePasswordToken(token)

		errHead = "The link has expired"
		errBody = "Please request a reset password email again."
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	hashedPwd, err := hashPassword(password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	if ok := databases.UpdatePassword(tokenObj.User, string(hashedPwd), token); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		return
	}

	databases.DeletePasswordToken(token)

	msgHead := "Reset Password Succeed"
	msgBody := "Now you can log in with the new password"
	c.Header("Location", utils.LoginPage)
	c.JSON(http.StatusCreated, gin.H{"msgHead": msgHead, "msgBody": msgBody})
}
