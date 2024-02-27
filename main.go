package main

import (
	"fmt"
	"os"

	"github.com/JohnnyOhms/projectx/config"
	"github.com/JohnnyOhms/projectx/controller"
	"github.com/JohnnyOhms/projectx/middleware"
	"github.com/JohnnyOhms/projectx/services"
	"github.com/gin-gonic/gin"
)

var (
	AuthService    services.AuthService      = services.New()
	AuthController controller.AuthController = controller.New(AuthService)
)

func init() {
	if err := config.Loadenv(); err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Failed to load environment variables:", err)
		} else {
			fmt.Println(".env file not found, using default values for environment variables.")
		}
	}

	config.ConnectToDB()
	config.SyncDB()
}
func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to the server endpoint",
		})
	})

	// Authentication
	r.POST("/api/auth/register", AuthController.SignUpUser)
	r.POST("/api/auth/login", AuthController.LoginUser)

	// User Info
	r.GET("/api/getinfo", AuthController.ReteriveUserDetails)
	r.POST("/api/setinfo", AuthController.SetUserDetails)

	// User image
	r.POST("/api/upload", AuthController.UploadAvatar)
	r.GET("/api/upload", middleware.Authorization, AuthController.RetrieveAvatar)

	// Discord auth
	r.GET("/api/auth/discord/redirect", AuthController.DiscordAuth)

	// Create the "avatar" directory if it doesn't exist
	if err := os.MkdirAll("avatars", os.ModePerm); err != nil {
		fmt.Println("Error creating 'avatar' directory:", err)
		return
	}

	// Check if PORT environment variable is set
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("PORT environment variable not set. Defaulting to port 9000.")
		port = "9000"
	} else {
		fmt.Println("Using PORT:", port)
	}

	// Run the server
	if err := r.Run(":" + port); err != nil {
		fmt.Println("Error running server:", err)
	}
}
