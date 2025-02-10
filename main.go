package main

import (
	"log"
	"time"

	"github.com/LinuxSploit/Livekit-Mircroservice/handler"
	"github.com/LinuxSploit/Livekit-Mircroservice/internal"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	internal.Init()
}

func main() {
	r := gin.Default()

	// Apply the CORS middleware globally
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "OPTION"}, // Add methods as needed
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},      // Explicitly allow Content-Type
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/HostLiveStream", handler.HostLiveHandler)

	r.POST("/JoinLiveStream", handler.JoinLiveHandler)

	r.GET("/ListLiveStreams", handler.ListLiveHandler)

	r.POST("/EndLiveStream", handler.EndLiveHandler)

	r.POST("webhook", handler.WebhookHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
