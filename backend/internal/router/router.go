package router

import (
	"chat/internal/chat"
	"chat/internal/topic"
	"chat/internal/user"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Mode string `json:"mode"`
}

var jwtKey = []byte("secret")

func New(userHandler *user.Handler, chatHandler *chat.Handler, topicHandler *topic.TopicHandler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.POST("/logout", userHandler.Logout)

		v1.GET("/rooms/list", chatHandler.GetRooms)
		v1.GET("/rooms/join/:roomId", chatHandler.JoinRoom)

		v1.Use(authMiddleware(false))
		v1.POST("/rooms/create", chatHandler.CreateRoom)
		v1.GET("/rooms/clients/:roomId", chatHandler.GetClients)

		v1.Use(authMiddleware(true))
		v1.GET("/topics", topicHandler.GetTopics)
		v1.POST("/topics", topicHandler.CreateTopic)
		v1.POST("/topics/:id", topicHandler.UpdateTopic)
		v1.POST("/topics/delete/:id", topicHandler.DeleteTopic)
	}
}

func Listen(addr string) error {
	return r.Run(addr)
}
