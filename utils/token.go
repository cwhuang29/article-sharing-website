package utils

import (
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/databases/models"
)

func GetPasswordResetTokenInstance(token string) models.Password {
	return databases.GetResetPasswordToken(token)
}

func StoreLoginToken(id, loginMaxAge int) string {
	token := GetUUID()

	if ok := databases.InsertLoginToken(id, token, loginMaxAge); !ok {
		return StoreLoginToken(id, loginMaxAge) // Try again if we got duplicate tokens
	}
	return token
}

func StorePasswordResetToken(id, maxAge int) string {
	token := GetUUID()

	if ok := databases.InsertResetPasswordToken(id, token, maxAge); !ok {
		return StorePasswordResetToken(id, maxAge) // Try again if we got duplicate tokens
	}
	return token
}

func ClearLoginToken(userID int, token string) {
	databases.DeleteLoginToken(userID, token)
}

func ClearPasswordResetToken(token string) {
	databases.DeletePasswordToken(token)
}

func ClearExpiredLoginTokens(id int) {
	databases.DeleteExpiredLoginTokens(id)
}

func ClearExpiredPasswordTokens(id int) {
	databases.DeleteExpiredPasswordTokens(id)
}
