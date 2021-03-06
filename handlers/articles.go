package handlers

import (
	"fmt"
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CreateArticleView(c *gin.Context) {
	yes, _ := isLoginedAdmin(c)
	if !yes {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Unauthorized."})
		return
	}

	c.HTML(http.StatusOK, "editor.html", gin.H{
		"currPageCSS": "css/editor.css",
		"title":       "Create New Post",
		"function":    "create",
	})
}

func UpdateArticleView(c *gin.Context) {
	yes, _ := isLoginedAdmin(c)
	if !yes {
		errMsg := "<div><strong>You are not allowed to perform this action</strong><p>Login if you are administrator.</p></div>"
		c.HTML(http.StatusUnauthorized, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"err":         errMsg,
		})
		return
	}

	id := checkArticleId(c, "articleId")
	if id == 0 {
		errMsg := "<div><p><strong>Article ID is an integer</strong></p><p>Please try again.</p></div>"
		c.HTML(http.StatusBadRequest, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"err":         errMsg,
		})
		return
	}

	if succeed := databases.IsArticleExists(id); succeed != true {
		errMsg := "<div><p><strong>Article ID Not Found</strong></p><p>Please try again.</p></div>"
		c.HTML(http.StatusNotFound, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"err":         errMsg,
		})
		return
	}

	dbFormatArticle, succeed := databases.GetArticleFullContent(id)
	if succeed != true {
		errMsg := "<div><p><strong>Article ID Not Found</strong></p><p>Please try again.</p></div>"
		c.HTML(http.StatusNotFound, "browse.html", gin.H{
			"currPageCSS": "css/browse.css",
			"err":         errMsg,
		})
		return
	}

	article := articleFormatDBToDetailed(dbFormatArticle, false)
	c.HTML(http.StatusOK, "editor.html", gin.H{
		"function":     "update",
		"currPageCSS":  "css/editor.css",
		"title":        "Edit: " + article.Title,
		"articleTitle": article.Title,
		"subtitle":     article.Subtitle,
		"date":         article.Date,
		"author":       article.Authors,
		"category":     article.Category,
		"tags":         article.Tags,
		"content":      article.Content,
	})
}

func CreateArticle(c *gin.Context) {
	yes, _ := isLoginedAdmin(c)
	if !yes {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Unauthorized."})
		return
	}

	var newArticle Article
	if err := c.ShouldBindJSON(&newArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": true, "err": err})
		return
	}

	invalids := validateArticleFormat(newArticle)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": false, "err": invalids})
		return
	}

	newArticle.Tags = removeDuplicateTags(newArticle.Tags)
	dbFormatArticle := articleFormatDetailedToDB(newArticle)
	id, res := databases.InsertArticle(dbFormatArticle)
	if res {
		c.Header("Location", "/articles/browse?articleId="+strconv.Itoa(id)) // With Location header, response.redirected will become true (if status code is 3XX. e.g., 201 is always false)
		c.JSON(http.StatusCreated, gin.H{"articleId": id})                   // See public/js/editor.js to see the difference between 201 & 302
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "err": "An error occurred while writing to DB."})
	}
}

func UploadImages(c *gin.Context) {
	yes, _ := isLoginedAdmin(c)
	if !yes {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Unauthorized."})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println("Upload images error:", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	createdImages := map[string]string{}
	/*
		Single Files uploaded:
		form:  &{map[username:[abc123]] map[uploadImage:[0xc0001625f0]]}
		form.File: map[uploadImage:[0xc0001625f0]]
		form.File's key: uploadImage
		filename: 641.jpeg
		size: 102359
		header: map[Content-Disposition:[form-data; name="uploadImage"; filename="641.jpeg"] Content-Type:[image/jpeg]]

		Multipe Files uploaded:
		fmt.Println("form: ", form)          // form:  &{map[] map[uploadImage:[0xc0004234f0 0xc000423540 0xc000423590 0xc0004235e0]]}
		fmt.Println("form.File:", form.File) // form.File: map[uploadImage:[0xc0004234f0 0xc000423540 0xc000423590 0xc0004235e0]]
	*/
	for key := range form.File {
		fileDir := "./public/upload/images/"
		for _, file := range form.File[key] {
			fmt.Println("filename:", file.Filename, "size:", file.Size, "header:", file.Header) // filename can't be trusted !!!!!!!

			fileName := time.Now().UTC().Format("20060102150405") + getUUID()
			fileType := ""
			if ext := strings.Split(file.Filename, "."); len(ext) > 1 {
				fileType = "." + ext[len(ext)-1]
			}
			storedName := fileDir + fileName + fileType
			err := c.SaveUploadedFile(file, storedName)
			if err != nil {
				fmt.Println("Store file error:", err)
			} else {
				createdImages[file.Filename] = "/upload/images/" + fileName + fileType
			}
		}
	}
	c.JSON(http.StatusCreated, createdImages)
}

func UpdateArticle(c *gin.Context) {
	yes, _ := isLoginedAdmin(c)
	if !yes {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Unauthorized."})
		return
	}

	var newArticle Article
	if err := c.ShouldBindJSON(&newArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": true, "err": err})
		return
	}

	invalids := validateArticleFormat(newArticle)
	if len(invalids) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": false, "err": invalids})
		return
	}

	id := checkArticleId(c, "articleId")
	if id == 0 {
		errMsg := "<div><p><strong>Article ID is an integer</strong></p><p>Please try again.</p></div>"
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"bindingError": false, "err": errMsg})
		return
	}

	if succeed := databases.IsArticleExists(id); succeed != true {
		errMsg := "<div><p><strong>Article ID Not Found</strong></p><p>Please try again.</p></div>"
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"bindingError": false, "err": errMsg})
		return
	}

	newArticle.Tags = removeDuplicateTags(newArticle.Tags)
	dbFormatArticle := articleFormatDetailedToDB(newArticle)
	dbFormatArticle.ID = id

	id, res := databases.ReplaceArticle(dbFormatArticle)
	if res {
		c.Header("Location", "/articles/browse?articleId="+strconv.Itoa(id))
		c.JSON(http.StatusCreated, gin.H{"articleId": id})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "err": "An error occurred while writing to DB."})
	}
}

func DeleteArticle(c *gin.Context) {
	yes, _ := isLoginedAdmin(c)
	if !yes {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Unauthorized."})
		return
	}

	id := checkArticleId(c, "articleId")
	if id == 0 {
		errMsg := "<div><p><strong>Article ID is an integer</strong></p><p>Please try again.</p></div>"
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"bindingError": false, "err": errMsg})
		return
	}

	if succeed := databases.IsArticleExists(id); succeed != true {
		errMsg := "<div><p><strong>Article ID Not Found</strong></p><p>Please try again.</p></div>"
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"bindingError": false, "err": errMsg})
		return
	}

	fmt.Println("DeleteArticle with id:", id)
	if res := databases.DeleteArticle(id); res {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "err": "An error occurred while writing to DB."})
	}
}
