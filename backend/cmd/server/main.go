package main

import (
	"chat/db"
	_ "chat/docs"
	"chat/internal/cache"
	"chat/internal/chat"
	"chat/internal/config"
	"chat/internal/router"
	"chat/internal/topic"
	"chat/internal/user"
	"flag"
	"log"
	"net"
)

var (
	configPath = flag.String("configPath", "configs/dev.json", "configs path")
)

// @title Chat Swagger API Doc.
// @version 2.0
// @description This is a chat server.
// @termsOfService http://swagger.io/terms/

// @host localhost:8080/api/v1
// @BasePath /
// @schemes http
func main() {
	flag.Parse()
	conf := config.LoadConfiguration(*configPath)

	dbConn, err := db.New(conf.DB)
	if err != nil {
		log.Fatalf("db connection err: %s", err)
		return
	}

	rc := cache.New(conf.Redis)
	if err != nil {
		log.Fatalf("redis initalize err: %s", err)
		return
	}

	// user
	userRepository := user.NewRepository(dbConn.GetDB())
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	// topic
	topicRepository := topic.NewRedisTopicRepository(rc)
	topicService := topic.NewTopicService(topicRepository)
	topicHandler := topic.NewTopicHandler(topicService)

	// chat
	c := chat.NewChat(rc)
	chatHandler := chat.NewHandler(c)

	// init router
	router.New(userHandler, chatHandler, topicHandler)

	go c.Start()
	router.Listen(net.JoinHostPort(conf.Server.Host, conf.Server.Port))
}
