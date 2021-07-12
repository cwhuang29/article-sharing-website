package validator

import (
	"time"

	"github.com/cwhuang29/article-sharing-website/constants"
	"github.com/cwhuang29/article-sharing-website/databases/models"
)

var (
	OldestDate, _ = time.Parse("2006-01-02", "1960-01-01")
)

func ValidateArticleForm(newArticle *models.Article) (err map[string]string) {
	/*
	 * Note:
	 * outline can be empty
	 * In this function we are not checking the max length of outline and content fields due to the following reason:
	 *     The outline is for overview only (not an important field) and the limit word count of content is 20,000 which is super large
	 *     So if the input words count is really too large, just let the database truncate the content
	 */
	err = make(map[string]string)

	if len(newArticle.Title) == 0 {
		err["title"] = errInputMsg["short"]
	} else if len(newArticle.Title) > constants.TitleBytesLimit {
		err["title"] = errInputMsg["long"]
	}

	if len(newArticle.Subtitle) > constants.SubtitleBytesLimit { // Subtitle can be empty
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

	if len(newArticle.Tags) > constants.TagsNumLimit {
		err["tags"] = errInputMsg["tagsTooMany"]
	} else {
		for _, t := range newArticle.Tags {
			if len(t.Value) > constants.TagsBytesLimit {
				err["tags"] = errInputMsg["tagsTooLong"]
				break
			}
		}
	}

	if len(newArticle.Content) == 0 {
		err["content"] = errInputMsg["empty"]
	}

	return
}
