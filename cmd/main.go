package main

import (
	log "github.com/sirupsen/logrus"
	"some-api/internal/server"
	"some-api/utils/db"
)

func main() {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.Info("I'm alive!")

	server := server.NewServer(db.New())
	log.Fatalln(server.Router.Start(":8000"))
}
