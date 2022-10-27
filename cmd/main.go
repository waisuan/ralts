package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"some-api/internal/server"
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

	s := server.NewServer()
	log.Fatalln(s.Router.Start(":8001"))
}
