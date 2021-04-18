package handlers

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/cwhuang29/article-sharing-website/utils"
	"time"
)

var (
	OldestDate, _ = time.Parse("2006-01-02", "1960-01-01")
	errInputMsg   = map[string]string{
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

func validateLoginFormat(email string, password string) (err map[string]string) {
	err = make(map[string]string)

	if len(email) == 0 {
		err["email"] = errInputMsg["empty"]
	} else if !isEmailValid(email) {
		err["email"] = errInputMsg["emailInvalid"]
	}

	if len(password) == 0 {
		err["password"] = errInputMsg["empty"]
	}
	return
}

func validateUserFormat(newUser models.User) (err map[string]string) {
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

func validateArticleValues(newArticle models.Article) (err map[string]string) {
	err = make(map[string]string)

	if len(newArticle.Title) == 0 {
		err["title"] = errInputMsg["short"]
	} else if len(newArticle.Title) > utils.TitleBytesLimit {
		err["title"] = errInputMsg["long"]
	}

	if len(newArticle.Subtitle) > utils.SubtitleBytesLimit { // Subtitle can be empty
		err["subtitle"] = errInputMsg["long"]
	}

	if newArticle.ReleaseDate == OldestDate {
		err["date"] = errInputMsg["dateFormat"]
	} else {
		// if time.Now().Truncate(time.Hour * 24).Sub(inpDate) < 0 {
		//     err["date"] = ErrInputMsg["dateFuture"]
		// }
		if OldestDate.Sub(newArticle.ReleaseDate) > 0 {
			err["date"] = errInputMsg["dateTooOld"]
		}
	}

	if len(newArticle.Tags) > utils.TagsNumLimit {
		err["tags"] = errInputMsg["tagsTooMany"]
	} else {
		for _, t := range newArticle.Tags {
			if len(t.Value) > utils.TagsBytesLimit {
				err["tags"] = errInputMsg["tagsTooLong"]
				break
			}
		}
	}

	if len(newArticle.Content) == 0 {
		err["content"] = errInputMsg["empty"]
	}
	/*
	 * Note:
	 * outline can be empty
	 * In this function we are not checking the max length of outline and content fields due to the following reason:
	 * The outline is for overview only (not an important field) and the limit word count of content is 20,000 which is super large
	 * So if the input word count are really too large, just let the database truncate them
	 */
	return
}
