package validator

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/cwhuang29/article-sharing-website/utils"
)

var (
	errInputMsg = map[string]string{
		"empty":            "The field can't be empty.",
		"long":             "This field can have no more than 255 bytes (1 alphabet - 1 byte/1 Chinese word - 3 bytes/1 Emoji - 4 bytes).",
		"dateTooOld":       "The date chosen should be greater than 1960-01-01.",
		"dateFuture":       "The date chosen can't be in the future.",
		"dateFormat":       "The format of date should be yyyy-mm-dd.",
		"tagsTooMany":      "You can target up to 5 tags at a time.",
		"tagsTooLong":      "Each tag can contain at most 20 bytes (1 alphabet - 1 byte/1 Chinese word - 3 bytes/1 Emoji - 4 bytes)",
		"emailInvalid":     "The email format is not correct.",
		"passwordTooShort": "Passwords must be at least 8 characters long.",
	}
)

func ValidateLoginForm(email string, password string) (err map[string]string) {
	err = make(map[string]string)

	if len(email) == 0 {
		err["email"] = errInputMsg["empty"]
	} else if !utils.IsEmailValid(email) {
		err["email"] = errInputMsg["emailInvalid"]
	}

	if len(password) == 0 {
		err["password"] = errInputMsg["empty"]
	}
	return
}

func ValidateRegisterForm(newUser models.User) (err map[string]string) {
	err = make(map[string]string)

	if len(newUser.FirstName) == 0 {
		err["first_name"] = errInputMsg["empty"]
	}

	if len(newUser.LastName) == 0 {
		err["last_name"] = errInputMsg["empty"]
	}

	if len(newUser.Password) == 0 {
		err["password"] = errInputMsg["empty"]
	} else if len(newUser.Password) < 8 {
		err["password"] = errInputMsg["passwordTooShort"]
	}

	if len(newUser.Email) == 0 {
		err["email"] = errInputMsg["empty"]
	}

	if len(newUser.Gender) == 0 {
		err["gender"] = errInputMsg["empty"]
	}

	if len(newUser.Major) == 0 {
		err["major"] = errInputMsg["empty"]
	}

	return err
}
