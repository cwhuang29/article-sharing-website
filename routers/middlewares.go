package routers

import (
	"github.com/cwhuang29/article-sharing-website/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		csrfHeaders := c.Request.Header["X-Csrf-Token"]
		csrfToken, _ := c.Cookie("csrf_token")

		if len(csrfHeaders) != 1 || csrfHeaders[0] != csrfToken {
			errHead := "Unauthorized"
			errBody := "You are not allowed to perform this action."
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"errHead": errHead, "errBody": errBody})
			return
		}
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		cookieEmail, _ := c.Cookie("login_email")

		userStatus, _ := handlers.GetUserStatus(c)
		if userStatus < handlers.IsAdmin {
			status := http.StatusUnauthorized
			if userStatus == handlers.IsMember {
				status = http.StatusForbidden
			}

			errHead := "Unauthorized"
			errBody := "You are not allowed to perform this action."
			c.AbortWithStatusJSON(status, gin.H{"errHead": errHead, "errBody": errBody}) // If use JSON(), handlers functions will be triggered subsequentlly
		}

		c.Next()

		fields := map[string]interface{}{
			"method":  c.Request.Method,
			"url":     c.Request.URL.String(),
			"status":  c.Writer.Status(),
			"latency": time.Since(t),
			"email":   cookieEmail,
		}
		logrus.WithFields(fields).Info("Admins related routes")
	}
}
