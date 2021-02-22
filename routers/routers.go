package routers

import (
	"github.com/cwhuang29/article-sharing-website/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var (
	router    = gin.Default() // Creates a gin router with default middleware: logger and recovery (crash-free) middleware
	htmlFiles = []string{
		"public/views/medication.html",
		"public/views/pharma.html",
		"public/views/about.html",
		"public/views/browse.html",
		"public/views/login.html",
		"public/views/overview.html",
		"public/views/articlesGenerator.html",
		"public/views/templates/header.tmpl",
		"public/views/templates/navbar.tmpl",
		"public/views/templates/footer.tmpl",
		"public/views/templates/error_message_notification.tmpl",
		"public/views/templates/notice_message_notification.tmpl",
		"public/views/templates/one_image.tmpl",
		"public/views/templates/two_images.tmpl",
		"public/views/templates/three_images.tmpl",
	}
)

func loadAssets() {
	router.Static("/upload/images", "public/upload/images")
	router.Static("/js", "public/js")
	router.Static("/css", "public/css") // Static serves files from the given file system root. Internally a http.FileServer is used
	router.Static("/assets", "public/assets")
	router.StaticFile("/favicon.ico", "public/assets/favicon-64.ico") // StaticFile registers a single route in order to serve a single file of the local filesystem
	router.LoadHTMLFiles(htmlFiles...)
}

func addRoutes() {
	admin := router.Group("/admin") // /overview/... -> /admin/overview/...
	{
		admin.GET("/overview", handlers.AdminOverview)
		admin.GET("/create/article", handlers.CreateArticleView)
		admin.GET("/delete/article", handlers.DeleteArticleView)
		admin.GET("/update/article", handlers.UpdateArticleView)

		admin.POST("/create/article", handlers.CreateArticle)
		admin.POST("/create/images", handlers.UploadImages)
	}

	articles := router.Group("/articles")
	{
		articles.GET("/", handlers.Browse)
		articles.GET("/browse", handlers.Browse)
		articles.GET("/weekly-update", handlers.WeeklyUpdate) // The main page
		articles.GET("/medication", handlers.Overview)
		articles.GET("/pharma", handlers.Overview)
	}

	router.GET("/register", handlers.RegisterView)
	router.POST("/register", handlers.Register)
	router.GET("/login", handlers.LoginView)
	router.POST("/login", handlers.LoginJSON)
	router.POST("/logout", handlers.Logout)
	router.GET("/password/reset", handlers.PasswordResetRequest)
	router.POST("/password/email", handlers.PasswordResetEmail)
	router.GET("/password/reset/:token", handlers.PasswordResetView)
	router.POST("/password/reset", handlers.PasswordReset)

	router.GET("/about", handlers.About)
	router.GET("/contact-us", handlers.ContactUs)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/articles/weekly-update")
	})
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "每週國際藥聞｜Irene報乎你知")
		// c.Fail(500, errors.New("something failed!"))
	})
}

func Router() {
	// gin.SetMode(gin.ReleaseMode)

	loadAssets()
	addRoutes()
	router.Run(":" + os.Getenv("APP_PORT"))
}
