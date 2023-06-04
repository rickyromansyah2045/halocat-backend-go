package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rickyromansyah2045/halocat-backend-go/auth"
	haloCatConfig "github.com/rickyromansyah2045/halocat-backend-go/config"
	"github.com/rickyromansyah2045/halocat-backend-go/constant"
	"github.com/rickyromansyah2045/halocat-backend-go/content"
	"github.com/rickyromansyah2045/halocat-backend-go/handler"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"github.com/rickyromansyah2045/halocat-backend-go/middleware"
	"github.com/rickyromansyah2045/halocat-backend-go/user"
)

func main() {
	isProduction := flag.Bool("production", true, "production mode?")
	flag.Parse()

	_, b, _, _ := runtime.Caller(0)
	projectRootPath := filepath.Join(filepath.Dir(b), "")
	envLocation := projectRootPath + "/.env"

	if *isProduction {
		envLocation = "C:/Users/ricky/OneDrive/Dokumen/Halodoc Cat/halocat-backend-go/.env"

	}

	if err := godotenv.Load(envLocation); err != nil {
		log.Fatal("error while loading or open .env file, err: ", err.Error())
	}

	// initial constants
	constant.InitGeneralConstant()
	constant.InitDBConstant()
	constant.InitAuthConstant()
	constant.InitRedisConstant()

	// initial database
	db := haloCatConfig.InitDB(*isProduction)

	// repositories
	userRepository := user.NewRepository(db)
	contentRepository := content.NewRepository(db)

	// services
	userSvc := user.NewService(userRepository)
	authSvc := auth.NewService()
	contentSvc := content.NewService(contentRepository, userRepository)

	// handlers
	userHandler := handler.NewUserHandler(userSvc, authSvc)
	contentHandler := handler.NewContentHandler(contentSvc, userSvc)

	// for activate release mode
	if *isProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	// gin app configuration
	app := gin.Default()
	app.SetTrustedProxies(nil)
	app.Static("/images", "./images")
	app.Use(gzip.Gzip(gzip.DefaultCompression))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "POST"},
		AllowHeaders:     []string{"Host", "Origin", "Content-Length", "Content-Type", "Authorization", "User-Agent", "X-Forwarded-For", "Accept-Encoding", "Connection"},
		ExposeHeaders:    []string{"Content-Length", "Content-Encoding"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// middleware
	mAuth := middleware.Auth(authSvc, userSvc)
	mAdminAuth := middleware.AdminAuth(authSvc, userSvc)

	// routing
	api := app.Group("/api/v1")
	{
		// >>>>>>>>>>>>>>> begin strict endpoint <<<<<<<<<<<<<<<

		// account settings
		api.GET("/users/data", mAuth, userHandler.GetUserData)
		api.PUT("/users/data/change", mAuth, userHandler.ChangeUserData)

		// users (for admin only)
		api.GET("/users", mAdminAuth, userHandler.GetAllUser)
		api.GET("/users/:id", mAdminAuth, userHandler.GetUserByID)
		api.PUT("/users/:id", mAdminAuth, userHandler.UpdateUser)
		api.POST("/users", mAdminAuth, userHandler.CreateUser)
		api.DELETE("/users/:id", mAdminAuth, userHandler.DeleteUser)

		// contents
		api.PUT("/contents/:id", mAuth, contentHandler.UpdateContent)
		api.POST("/contents", mAuth, contentHandler.CreateContent)
		api.DELETE("/contents/:id", mAuth, contentHandler.DeleteContent)

		// contents -> images
		api.POST("/contents/images", mAuth, contentHandler.UploadImage)
		api.DELETE("/contents/images/:id", mAuth, contentHandler.DeleteContentImage)

		// contents -> categories (for admin only)
		api.PUT("/contents/categories/:id", mAdminAuth, contentHandler.UpdateContentCategory)
		api.POST("/contents/categories", mAdminAuth, contentHandler.CreateContentCategory)
		api.DELETE("/contents/categories/:id", mAdminAuth, contentHandler.DeleteContentCategory)

		// admin datatables
		api.GET("admin/datatables/users", mAdminAuth, userHandler.AdminDataTablesUsers)
		api.GET("admin/datatables/categories", mAdminAuth, contentHandler.AdminDataTablesCategories)
		api.GET("admin/datatables/contents", mAdminAuth, contentHandler.AdminDataTablesContents)

		// datatables for user
		api.GET("datatables/contents", mAuth, contentHandler.UserDataTablesContents)

		// >>>>>>>>>>>>>>> end strict endpoint <<<<<<<<<<<<<<<

		// >>>>>>>>>>>>>>> begin non-strict endpoint <<<<<<<<<<<<<<<

		// authentication
		api.POST("/users/register", userHandler.Register)
		api.POST("/users/login", userHandler.Login)

		// forgot password
		api.GET("/users/forgot-password/:token", userHandler.ProcessForgotPasswordToken)
		api.POST("/users/forgot-password", userHandler.CreateForgotPasswordToken)

		// users
		api.GET("/users/name/:id", userHandler.GetNameByID)

		// contents
		api.GET("/contents", contentHandler.GetAllContent)
		api.GET("/contents/:id", contentHandler.GetContentByID)

		// contents -> images
		api.GET("/contents/images", contentHandler.GetAllContentImage)
		api.GET("/contents/images/:id", contentHandler.GetContentImageByID)

		// contents -> categories
		api.GET("/contents/categories", contentHandler.GetAllContentCategory)
		api.GET("/contents/categories/:id", contentHandler.GetContentCategoryByID)

		// >>>>>>>>>>>>>>> end non-strict endpoint <<<<<<<<<<<<<<<
	}

	// handle invalid method
	app.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, helper.BasicAPIResponseError(http.StatusMethodNotAllowed, "Request invalid, invalid method!"))
	})

	// handle invalid path or invalid endpoint
	app.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, helper.BasicAPIResponseError(http.StatusNotFound, "Request invalid, path not found!"))
	})

	// run http server
	app.Run(os.Getenv("APP_RUN_ON"))
}
