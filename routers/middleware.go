package routers

import (
	"net/http"
	"time"

	"github.com/cwhuang29/article-sharing-website/constants"
	"github.com/cwhuang29/article-sharing-website/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		csrfHeaders := c.Request.Header["X-Csrf-Token"]
		csrfToken, _ := c.Cookie(constants.CookieCSRFToken)

		if len(csrfHeaders) != 1 || csrfToken == "" || csrfHeaders[0] != csrfToken {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"errHead": constants.GeneralErr, "errBody": constants.PermissionDenied})
		}
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		userStatus, _ := handlers.GetUserStatus(c)
		if userStatus < handlers.IsAdmin {
			status := http.StatusUnauthorized
			if userStatus == handlers.IsMember {
				status = http.StatusForbidden
			}
			// If use JSON(), handler functions will be triggered subsequentlly
			c.AbortWithStatusJSON(status, gin.H{"errHead": constants.GeneralErr, "errBody": constants.PermissionDenied})
		}

		c.Next()

		cookieEmail, _ := c.Cookie(constants.CookieLoginEmail)
		fields := map[string]interface{}{
			"method":  c.Request.Method,
			"url":     c.Request.URL.String(),
			"status":  c.Writer.Status(),
			"latency": time.Since(t),
			"email":   cookieEmail,
		}
		logrus.WithFields(fields).Info("Admins routes")
	}
}
