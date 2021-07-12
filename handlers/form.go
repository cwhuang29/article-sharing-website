package handlers

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/cwhuang29/article-sharing-website/constants"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/cwhuang29/article-sharing-website/utils"
	"github.com/cwhuang29/article-sharing-website/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	acceptedFileType = map[string][]string{"image": {"image/png", "image/jpeg", "image/gif", "image/webp", "image/apng"}}
)

func writeFileLog(fileName, fileSize, fileType string) {
	fields := map[string]interface{}{
		"fileName": fileName,
		"fileType": fileType,
		"fileSize": fileSize,
	}
	logrus.WithFields(fields).Info("File upload")
}

func mapFilesName(content string, fileNamesMapping map[string]string) string {
	for key := range fileNamesMapping {
		content = strings.Replace(content, key, fileNamesMapping[key], -1)
	}

	return content
}

func saveFile(c *gin.Context, file *multipart.FileHeader, fileName string) (err error) {
	err = c.SaveUploadedFile(file, fileName)
	if err != nil {
		logrus.Errorf("Create article error when saving images:", err)
	}
	return
}

func checkFileSize(fileSize int64) bool {
	return fileSize <= int64(constants.FileMaxSize)
}

func checkFileType(fileType, mainType string) bool {
	for _, t := range acceptedFileType[mainType] {
		if fileType == t {
			return true
		}
	}
	return false
}

func generateFileName(fileType string) string {
	fileID := time.Now().UTC().Format("20060102150405") + utils.GetUUID()
	fileExt := fileType[strings.LastIndex(fileType, "/")+1:]
	return constants.UploadImageDir + fileID + "." + fileExt // Do not start with "./" otherwise the images URL in articles content will be incorrect
}

/*
 * Notice: Even though I have renamed the files (in the 3rd argument of JS formData.append API), Filenames can't be trusted
 * fmt.Println("filename:", file.Filename, "size:", file.Size, "header:", file.Header)
 * filename: d5821d5a77.png size: 5387170 header: map[Content-Disposition:[form-data; name="uploadImages"; filename="d5821d5a77.png"] Content-Type:[image/png]]
 */
func checkFileAndRename(file *multipart.FileHeader) (fileName string, err error) {
	if ok := checkFileSize(file.Size); !ok {
		err = fmt.Errorf("File size of %v is too large (max: 8MB per image)!", file.Filename)
		return
	}

	fileType := file.Header.Get("Content-Type")
	if ok := checkFileType(fileType, "image"); !ok {
		err = fmt.Errorf("File type of %v is not permitted!", file.Filename)
		return
	}

	fileName = generateFileName(fileType)
	return

}

func getImagesInContent(files []*multipart.FileHeader) (fileNames []string, fileNamesMapping map[string]string, err error) {
	var fileName string
	fileNames = make([]string, len(files))
	fileNamesMapping = make(map[string]string, len(files))

	for i, file := range files {
		if fileName, err = checkFileAndRename(file); err != nil {
			return
		}

		fileNames[i] = fileName
		fileNamesMapping[file.Filename] = fileName[len("public/"):] // Get rid of the prefix since we truncate it in router.Static()
		writeFileLog(fileName, strconv.FormatInt(file.Size, 10), file.Header.Get("Content-Type"))
	}
	return
}

func getCoverPhoto(file *multipart.FileHeader) (coverPhotoName string, err error) {
	if coverPhotoName, err = checkFileAndRename(file); err != nil {
		return
	}

	writeFileLog(coverPhotoName, strconv.FormatInt(file.Size, 10), file.Header.Get("Content-Type"))
	return
}

func saveFilesFromForm(c *gin.Context, files map[string][]*multipart.FileHeader) (fileNamesMapping map[string]string, coverPhotoURL string, err error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Error occurred when retrieving files from the input form:", err)
		}
	}()

	var coverPhotoName string
	var coverPhoto *multipart.FileHeader

	if len(files["coverPhoto"]) > 0 {
		coverPhoto = files["coverPhoto"][0] // There is only one cover photo
		if coverPhotoName, err = getCoverPhoto(coverPhoto); err != nil {
			return
		}
	}

	var fileNames []string
	contentImages := files["contentImages"]

	if fileNames, fileNamesMapping, err = getImagesInContent(contentImages); err != nil {
		return
	}

	// All files were retrieved. Start saving them

	if coverPhoto != nil { // User may not upload cover photo
		_ = saveFile(c, coverPhoto, coverPhotoName)
		coverPhotoURL = coverPhotoName[len("public/"):]
	}

	for i, name := range fileNames {
		if err = saveFile(c, contentImages[i], name); err != nil {
			return
		}
	}

	return
}

func getValuesFromForm(c *gin.Context, formVal map[string][]string) (*models.Article, error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Create article error when retrieving values from form:", err)
		}
	}()

	date, err := time.Parse("2006-01-02", formVal["date"][0])
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &models.Article{
		AdminOnly:   adminOnly,
		Title:       strings.TrimSpace(formVal["title"][0]), // If the form does not contain "title" field, the array's value extraction will panic
		Subtitle:    strings.TrimSpace(formVal["subtitle"][0]),
		ReleaseDate: date,
		Authors:     auths,
		Category:    formVal["category"][0],
		Tags:        tags,
		Outline:     formVal["outline"][0],
		Content:     formVal["content"][0],
	}, nil
}

func handleForm(c *gin.Context) (newArticle *models.Article, invalids map[string]string, err error) {
	var form *multipart.Form

	form, err = c.MultipartForm() // form: &{map[authors:[Jasia] category:[Medication] ... title:[abcde]] map[uploadImages:[0xc0001f91d0 0xc0001f8000]]}
	if err != nil {
		return
	}

	newArticle, err = getValuesFromForm(c, form.Value)
	if err != nil {
		return
	}
	if newArticle.Title == "" {
		err = fmt.Errorf("Error occurred when extracting values from form.")
		return
	}

	invalids = validator.ValidateArticleForm(newArticle)
	if len(invalids) != 0 {
		return
	}

	fileNamesMapping, coverPhotoURL, err := saveFilesFromForm(c, form.File)
	if err != nil {
		return
	}

	newArticle.Content = mapFilesName(newArticle.Content, fileNamesMapping)
	newArticle.CoverPhoto = coverPhotoURL

	return
}
