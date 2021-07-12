package handlers

import (
	"net/http"

	"github.com/cwhuang29/article-sharing-website/config"
	"github.com/cwhuang29/article-sharing-website/constants"
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func PasswordResetRequest(c *gin.Context) {
	c.HTML(http.StatusOK, "passwordResetRequest.html", gin.H{"title": "Reset Password"})
}

func PasswordResetEmail(c *gin.Context) {
	var json = struct{ Email string }{} // Uppercase is required

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.ReloadAndRetry})
		return
	}

	user := databases.GetUser(json.Email)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"errHead": constants.EmailNotFound, "errBody": constants.EmailIsAddressCorrect})
		return
	}

	/*
	 * There should be a mechanism to remove expired tokens in Password table
	 * Currently I'll do it right here, and use go routine to clear expired tokens periodically in the future
	 */
	utils.ClearExpiredPasswordTokens(user.ID)

	if yes := utils.DoesUserHasEmailQuota(user.ID); !yes {
		c.JSON(http.StatusTooManyRequests, gin.H{"errHead": constants.TryTooOften, "errBody": constants.EmailTryLater})
		return
	}

	baseURL := config.GetCopy().App.Url
	token := utils.StorePasswordResetToken(user.ID, constants.ResetPasswordTokenMaxAge)
	link := baseURL + constants.URLResetPassword + token + "?email=" + user.Email
	expireMins := constants.ResetPasswordTokenMaxAge / 60
	name := user.FirstName + " " + user.LastName

	if ok := SendResetPasswordEmail(user.Email, name, link, expireMins); !ok {
		logrus.Error("[Failed] Password reset email. Sent to ", user.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.ReloadAndRetry})
		return
	}

	msgHead := "An email has been sent"
	msgBody := "Reset link has been sent to your email"
	logrus.Info("[Succeed] Password reset email. Sent to ", user.Email)
	c.JSON(http.StatusOK, gin.H{"msgHead": msgHead, "msgBody": msgBody})
}

func PasswordResetForm(c *gin.Context) {
	email := getQueryPara(c, constants.QueryEmail)
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.EmailRequestAgain})
		return
	}

	uuid := utils.GetUUID()
	c.SetCookie(constants.CookieCSRFToken, uuid, constants.CsrfTokenAge, "/", "", true, true)
	c.HTML(http.StatusOK, "passwordResetForm.html", gin.H{
		"title":       "Reset Password",
		"csrfToken":   uuid,
		"currPageCSS": "", // Fields can't be omitted even not using
		"email":       email,
	})
}

func PasswordUpdate(c *gin.Context) {
	var json = struct {
		Email    string
		Password string
		Token    string
	}{}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.EmailOpenAgain})
		return
	}

	email, password, token := json.Email, json.Password, json.Token
	if email == "" || password == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.QueryMissingErr, "errBody": constants.EmailOpenAgain})
		return
	}

	tokenObj := utils.GetPasswordResetTokenInstance(token)
	if email != tokenObj.User.Email {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.EmailOutdated})
		return
	}

	if utils.IsExpired(tokenObj.CreatedAt, tokenObj.MaxAge) {
		databases.DeletePasswordToken(token)
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.EmailLinkExpired, "errBody": constants.EmailRequestAgain})
		return
	}

	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.ReloadAndRetry})
		return
	}

	if ok := databases.UpdatePassword(tokenObj.User, string(hashedPwd), token); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errHead": constants.UnexpectedErr, "errBody": constants.ReloadAndRetry})
		return
	}

	utils.ClearPasswordResetToken(token)

	msgHead := "Reset Password Succeed"
	msgBody := "Now you can log in with the new password"
	c.Header("Location", constants.URLLoginPage)
	c.JSON(http.StatusCreated, gin.H{"msgHead": msgHead, "msgBody": msgBody})
}
