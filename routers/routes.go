package routers

import (
	"github.com/cwhuang29/article-sharing-website/config"
	"github.com/cwhuang29/article-sharing-website/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	router    = gin.Default() // Creates a gin router with default middleware: logger and recovery (crash-free) middleware
	htmlFiles = []string{
		"public/views/medication.html",
		"public/views/pharma.html",
		"public/views/about.html",
		"public/views/browse.html",
		"public/views/login.html",
		"public/views/register.html",
		"public/views/overview.html",
		"public/views/editor.html",
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
		articles.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/articles/weekly-update")
		})
		articles.GET("/browse", handlers.Browse)
		articles.GET("/weekly-update", handlers.Overview) // The main page
		articles.GET("/medication", handlers.Overview)
		articles.GET("/pharma", handlers.Overview)
		articles.GET("/fetch", handlers.FetchData)
		articles.GET("/tags", handlers.SearchTags)
	}

	router.GET("/register", handlers.RegisterView)
	router.GET("/login", handlers.LoginView)
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.LoginJSON)
	router.POST("/logout", handlers.Logout)

	router.GET("/password/reset", handlers.PasswordResetRequest)
	router.GET("/password/reset/:token", handlers.PasswordResetView)
	router.POST("/password/email", handlers.PasswordResetEmail)
	router.POST("/password/reset", handlers.PasswordReset)

	router.GET("/about", handlers.About)
	router.GET("/contact-us", handlers.ContactUs)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/articles/weekly-update")
	})
}

func serve() {
	cfg := config.GetConfig()

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
