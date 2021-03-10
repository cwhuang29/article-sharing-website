package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	fileDir          = "public/upload/images/" // Do not start with "./" otherwise the images URL in articles content will be incorrect
	acceptedFileType = map[string][]string{"image": {"image/png", "image/jpeg", "image/gif"}}
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

func getValuesFromForm(c *gin.Context, formVal map[string][]string) Article {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Create article error when retrieving values from form:", err)
			errHead := "An Error Occurred"
			errBody := "Please try again."
			c.JSON(http.StatusBadRequest, gin.H{"errHead": errHead, "errBody": errBody})
		}
	}()

	return Article{
		Title:    formVal["title"][0], // If the form does not contain "title" field, the array's value extraction will panic
		Subtitle: formVal["subtitle"][0],
		Date:     formVal["date"][0],
		Authors:  formVal["authors"],
		Category: formVal["category"][0],
		Tags:     formVal["tags"],
		Content:  formVal["content"][0],
	}
}

func handleForm(c *gin.Context) (newArticle Article, err error) {
	var form *multipart.Form
	form, err = c.MultipartForm() // form: &{map[authors:[Jasia] category:[Medication] ... title:[abcde]] map[uploadImages:[0xc0001f91d0 0xc0001f8000]]}
	if err != nil {
		return
	}

	newArticle = getValuesFromForm(c, form.Value)

	var fileNamesMapping map[string]string
	if fileNamesMapping, err = getFilesFromForm(c, form.File["uploadImages"]); err == nil {
		for key := range fileNamesMapping {
			newArticle.Content = strings.Replace(newArticle.Content, key, fileNamesMapping[key], -1)
		}
	}
	return
}