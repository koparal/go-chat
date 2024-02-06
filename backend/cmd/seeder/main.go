package main

import (
	"chat/db"
	"chat/internal/config"
	"chat/internal/user"
	"chat/internal/utils"
	"flag"
	"fmt"
	"log"
)

var (
	configPath = flag.String("configPath", "configs/dev.json", "configs path")
)

func main() {
	flag.Parse()
	conf := config.LoadConfiguration(*configPath)

	dbConn, err := db.New(conf.DB)
	if err != nil {
		log.Fatalf("db connection err: %s", err)
		return
	}

	admin := user.User{
		Username: "admin",
		Password: "123",
		IsAdmin:  true,
	}

	hashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		panic(err)
	}

	_, err = dbConn.DB.Exec("INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3)",
		admin.Username, hashedPassword, admin.IsAdmin)
	if err != nil {
		panic(err)
	}

	fmt.Println("Admin user created.")
}
