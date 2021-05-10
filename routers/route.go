package routers

import (
	"net/http"

	"github.com/cwhuang29/article-sharing-website/config"
	"github.com/cwhuang29/article-sharing-website/handlers"
	"github.com/gin-gonic/gin"
)

var (
	router    = gin.Default() // Creates a gin router with default middleware: logger and recovery (crash-free) middleware
	htmlFiles = []string{
		"public/views/medication.html",
		"public/views/pharma.html",
		"public/views/about.html",
		"public/views/home.html",
		"public/views/browse.html",
		"public/views/login.html",
		"public/views/register.html",
		"public/views/overview.html",
		"public/views/editor.html",
		"public/views/auth/passwordResetRequest.html",
		"public/views/auth/passwordResetForm.html",
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
	router.LoadHTMLFiles(htmlFiles...)                                // router.LoadHTMLGlob("public/*")
}

func injectRoutes() {
	admin := router.Group("/admin") // /overview/... -> /admin/overview/...
	admin.Use(AdminRequired())
	{
		admin.GET("/overview", handlers.AdminOverview)
		admin.GET("/check-permisssion", handlers.CheckPermissionAndArticleExists)
		admin.GET("/create/article", handlers.CreateArticleView)
		admin.GET("/update/article", handlers.UpdateArticleView)

		admin.Use(CSRFProtection())
		{
			admin.POST("/create/article", handlers.CreateArticle)
			admin.PUT("/update/article", handlers.UpdateArticle)
			admin.DELETE("/delete/article", handlers.DeleteArticle)
		}
	}

	articles := router.Group("/articles")
	{
		articles.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/articles/weekly-update") })
		articles.GET("/weekly-update", handlers.Overview) // The main page
		articles.GET("/browse", handlers.Browse)
		articles.GET("/medication", handlers.Overview)
		articles.GET("/pharma", handlers.Overview)
		articles.GET("/fetch", handlers.FetchData)
		articles.GET("/tags", handlers.SearchTags)

		articles.GET("/bookmark", handlers.GetUserBookmarkedArticles)
		articles.GET("/bookmark/:articleId", handlers.Bookmark)
		articles.PUT("/bookmark/:articleId", handlers.UpdateBookmark)

		articles.GET("/like/:articleId", handlers.Like)
		articles.PUT("/like/:articleId", handlers.UpdateLike)
	}

	password := router.Group("/password")
	{
		password.GET("/reset", handlers.PasswordResetRequest)
		password.GET("/reset/:token", handlers.PasswordResetForm)
		password.POST("/email", handlers.PasswordResetEmail)
		password.Use(CSRFProtection())
		{
			password.PUT("/reset", handlers.PasswordUpdate)
		}
	}

	router.GET("/register", handlers.RegisterView)
	router.GET("/login", handlers.LoginView)
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)
	router.POST("/logout", handlers.Logout)

	router.GET("/home", handlers.Home)
	router.GET("/about", handlers.About)
	router.GET("/contact-us", handlers.ContactUs)
	router.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/articles/weekly-update") })
}

func serve() {
	cfg := config.GetCopy()

	http := cfg.App.HttpPort
	https := cfg.App.HttpsPort

	if http != "" && https != "" {
		go router.Run(":" + http)
		router.RunTLS(":"+https, "./certs/server.crt", "./certs/server.key")
	} else if http != "" {
		router.Run(":" + http)
	} else if https != "" {
		router.RunTLS(":"+https, "./certs/server.crt", "./certs/server.key")
	} else {
		panic("Either app.httpPort or app.HttpsPort should be set")
	}
}

func Router() {
	// gin.SetMode(gin.ReleaseMode)
	router.MaxMultipartMemory = 16 << 20 // Set a lower memory limit for multipart forms (default is 32 MiB)

	loadAssets()
	injectRoutes()
	serve()
}
