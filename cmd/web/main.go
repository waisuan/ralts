package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"ralts/internal/chat"
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

	dbClient := db.NewRaltsDatabase(os.Getenv("DATABASE_URL"))
	defer dbClient.Close()

	chatHandler := chat.NewChat(dbClient)

	s := server.NewServer(chatHandler)
	log.Fatalln(s.Router.Start(":8001"))
}
