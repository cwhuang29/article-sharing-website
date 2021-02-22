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
	cookieEmail, err := c.Cookie("login_email")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Page not found"})
		return
	}

	cookieToken, err := c.Cookie("login_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Page not found"})
		return
	}

	if cookieToken != databases.GetLoginToken(cookieEmail) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Page not found"})
		return
	}

	c.HTML(http.StatusOK, "articlesGenerator.html", gin.H{
		"currPageCSS": "css/articlesGenerator.css",
		"title":       "Create New Post",
		"function":    "create",
	})
}

func UpdateArticleView(c *gin.Context) {
	// postName := "Old news title"
	// c.HTML(http.StatusOK, "articlesGenerator.html", gin.H{
	//     "currPageCSS": "css/articlesGenerator.css",
	//     "title":       "Update Post: " + postName,
	//     "function":    "update",
	// })
}

func DeleteArticleView(c *gin.Context) {
}

func CreateArticle(c *gin.Context) {
	cookieEmail, err := c.Cookie("login_email")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Page not found"})
		return
	}

	cookieToken, err := c.Cookie("login_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Page not found"})
		return
	}

	if cookieToken != databases.GetLoginToken(cookieEmail) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Page not found"})
		return
	}

	var newArticle Article
	if err := c.ShouldBindJSON(&newArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": true, "err": err})
		return
	}

	invalids := validateCreateArticle(newArticle)
	if len(invalids) == 0 {
		newArticle.Tags = removeDuplicateTags(newArticle.Tags)
		dbFormatArticle := ArticleFormatDetailedToDB(newArticle)
		id, res := databases.InsertArticleToDB(dbFormatArticle)
		if res {
			c.Header("Location", "/articles/browse?articleId="+strconv.Itoa(id)) // With Location header, response.redirected will become true (if status code is 3XX. e.g., 201 is always false)
			c.JSON(http.StatusCreated, gin.H{"articleId": id})                   // See public/js/articlesGenerator.js to see the difference between 201 & 302
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"bindingError": false, "err": "An error occurred while writing to DB."})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"bindingError": false, "err": err})
	}
}

func UploadImages(c *gin.Context) {
	cookieEmail, err := c.Cookie("login_email")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Page not found"})
		return
	}

	cookieToken, err := c.Cookie("login_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Page not found"})
		return
	}

	if cookieToken != databases.GetLoginToken(cookieEmail) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Page not found"})
		return
	}

	createdImages := map[string]string{}
	if form, err := c.MultipartForm(); err != nil {
		fmt.Println("-------- ERROR  --------")
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, createdImages)
	} else {
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

				fileName := time.Now().Format("20060102150405") + getUUID()
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
}
