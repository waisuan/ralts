package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"ralts/internal/chat"
	"ralts/internal/config"
	db "ralts/internal/db"
	"ralts/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.Info("I'm alive!")

	cfg := config.NewConfig(false)

	dbClient := db.NewRaltsDatabase(cfg)
	defer dbClient.Close()

	chatHandler := chat.NewChat(dbClient)
	h := server.Handlers{
		ChatHandler: chatHandler,
	}

	s := server.NewServer(&h, cfg)
	log.Fatalln(s.Router.Start(":8001"))
}
