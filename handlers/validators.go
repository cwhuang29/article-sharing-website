package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"time"
)

var (
	OldestDate, _ = time.Parse("2006-01-02", "1960-01-01")
	TagsLmit      = 5
	TagsCharLmit  = 65 // TODO length of each emoji is about 13
	ErrInputMsg   = map[string]string{
		"empty":            "The field can't be empty.",
		"long":             "This field can have no more than 255 characters.",
		"dateTooOld":       "The date chosen should be greater than 1960-01-01.",
		"dateFuture":       "The date chosen can't be in the future.",
		"dateFormat":       "The format of date should be yyyy-mm-dd.",
		"tagsTooMany":      "You can target up to 5 tags at a time.",
		"tagsTooLong":      "Each tag can contain at most 20 charaters.",
		"emailInvalid":     "The email format is not correct.",
		"passwordTooShort": "Passwords must be at least 8 characters long.",
	}
)

func validateLoginFormat(email string, password string) (err map[string]interface{}) {
	err = make(map[string]interface{})

	if len(email) == 0 {
		err["email"] = ErrInputMsg["empty"]
	} else if !isEmailValid(email) {
		err["email"] = ErrInputMsg["emailInvalid"]
	}

	if len(password) == 0 {
		err["password"] = ErrInputMsg["empty"]
	}
	return
}

func validateUserFormat(newUser models.User) (err map[string]interface{}) {
	err = make(map[string]interface{})

	if len(newUser.FirstName) == 0 {
		err["first_name"] = ErrInputMsg["empty"]
	}

	if len(newUser.LastName) == 0 {
		err["last_name"] = ErrInputMsg["empty"]
	}

	if len(newUser.Password) == 0 {
		err["password"] = ErrInputMsg["empty"]
	} else if len(newUser.Password) < 8 {
		err["password"] = ErrInputMsg["passwordTooShort"]
	}

	if len(newUser.Email) == 0 {
		err["email"] = ErrInputMsg["empty"]
	}

	if len(newUser.Gender) == 0 {
		err["gender"] = ErrInputMsg["empty"]
	}

	if len(newUser.Major) == 0 {
		err["major"] = ErrInputMsg["empty"]
	}

	return err
}

func validateArticleValues(newArticle models.Article) (err map[string]interface{}) {
	err = make(map[string]interface{})

	if len(newArticle.Title) == 0 {
		err["title"] = ErrInputMsg["short"]
	} else if len(newArticle.Title) > 255 {
		err["title"] = ErrInputMsg["long"]
	}

	if len(newArticle.Subtitle) > 255 { // Subtitle can be empty
		err["subtitle"] = ErrInputMsg["long"]
	}

	if newArticle.ReleaseDate == OldestDate {
		err["date"] = ErrInputMsg["dateFormat"]
	} else {
		// if time.Now().Truncate(time.Hour * 24).Sub(inpDate) < 0 {
		//     err["date"] = ErrInputMsg["dateFuture"]
		// }
		if OldestDate.Sub(newArticle.ReleaseDate) > 0 {
			err["date"] = ErrInputMsg["dateTooOld"]
		}
	}

	if len(newArticle.Tags) > TagsLmit {
		err["tags"] = ErrInputMsg["tagsTooMany"]
	} else {
		for _, t := range newArticle.Tags {
			if len(t.Value) > TagsCharLmit {
				err["tags"] = ErrInputMsg["tagsTooLong"]
				break
			}
		}
	}

	if len(newArticle.Content) == 0 {
		err["content"] = ErrInputMsg["empty"]
	}
	return
}
