package handlers

import (
	"fmt"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

var (
	fileDir          = "public/upload/images/" // Do not start with "./" otherwise the images URL in articles content will be incorrect
	acceptedFileType = map[string][]string{"image": {"image/png", "image/jpeg", "image/gif", "image/webp", "image/apng"}}
	fileMaxSize      = 4 * 1000 * 1000 // 4MB
)

func writeFileLog(fileName, fileType, fileSize string) {
	fields := map[string]interface{}{
		"fileName": fileName,
		"fileType": fileType,
		"fileSize": fileSize,
	}
	logrus.WithFields(fields).Info("File upload")
}

func checkAndFilterFileType(mainType string, fileType string) string {
	for _, t := range acceptedFileType[mainType] {
		if fileType == t {
			return fileType[strings.LastIndex(fileType, "/")+1:]
		}
	}
	return ""
}

/*
 * Notice: Even though I have renamed the files (in the 3rd argument of JS formData.append API), Filenames can't be trust
 * fmt.Println("filename:", file.Filename, "size:", file.Size, "header:", file.Header)
 * filename: d5821d5a77.png size: 5387170 header: map[Content-Disposition:[form-data; name="uploadImages"; filename="d5821d5a77.png"] Content-Type:[image/png]]
 */
func getFilesFromForm(c *gin.Context, files []*multipart.FileHeader) (fileNamesMapping map[string]string, err error) {
	fileNamesMapping = make(map[string]string, len(files))
	for _, file := range files {
		if file.Size > int64(fileMaxSize) {
			err = fmt.Errorf("File size of %v is too large!", file.Filename)
			return
		}

		filteredType := checkAndFilterFileType("image", file.Header.Get("Content-Type"))
		if filteredType == "" {
			err = fmt.Errorf("File type of %v is not permitted!", file.Filename)
			return
		}

		fileID := time.Now().UTC().Format("20060102150405") + getUUID()
		fileName := fileDir + fileID + "." + filteredType
		fileNamesMapping[file.Filename] = fileName[7:] // Get rid of "public/" prefix

		err := c.SaveUploadedFile(file, fileName)
		if err != nil {
			logrus.Errorf("Create article error when saving images:", err)
		}

		writeFileLog(fileName, strconv.FormatInt(file.Size, 10), file.Header.Get("Content-Type"))
	}
	return
}

func getValuesFromForm(c *gin.Context, formVal map[string][]string) models.Article {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Create article error when retrieving values from form:", err)
		}
	}()

	date, err := time.Parse("2006-01-02", formVal["date"][0])
	if err != nil {
		date = OldestDate
	}

	auths := strings.Join(formVal["authors"], ",")

	tags := []models.Tag{}
	if formVal["tags"][0] != "" {
		formTags := strings.Split(formVal["tags"][0], ",") // JS's form.append() transformed array to a comma seperated string
		tags = make([]models.Tag, len(formTags))
		for i, t := range formTags {
			tags[i].Value = strings.TrimSpace(t)
		}
	}

	adminOnly, err := strconv.ParseBool(formVal["adminOnly"][0])
	if err != nil {
		adminOnly = true // Since there may be some unexpected errors, hide this article from non-admins
	}

	return models.Article{
		AdminOnly:   adminOnly,
		Title:       strings.TrimSpace(formVal["title"][0]), // If the form does not contain "title" field, the array's value extraction will panic
		Subtitle:    strings.TrimSpace(formVal["subtitle"][0]),
		ReleaseDate: date,
		Authors:     auths,
		Category:    formVal["category"][0],
		Tags:        tags,
		Content:     formVal["content"][0],
	}
}

func handleForm(c *gin.Context) (newArticle models.Article, invalids map[string]string, err error) {
	var form *multipart.Form
	var fileNamesMapping map[string]string

	form, err = c.MultipartForm() // form: &{map[authors:[Jasia] category:[Medication] ... title:[abcde]] map[uploadImages:[0xc0001f91d0 0xc0001f8000]]}
	if err != nil {
		return
	}

	newArticle = getValuesFromForm(c, form.Value)
	if newArticle.Title == "" {
		err = fmt.Errorf("Error occurred when extracting values from form.")
		return
	}
	invalids = validateArticleValues(newArticle)
	if len(invalids) != 0 {
		return
	}

	if fileNamesMapping, err = getFilesFromForm(c, form.File["uploadImages"]); err == nil {
		for key := range fileNamesMapping {
			newArticle.Content = strings.Replace(newArticle.Content, key, fileNamesMapping[key], -1)
		}
	}
	return
}
