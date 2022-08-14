package main

import (
	log "github.com/sirupsen/logrus"
	"some-api/services/api"
)

func main() {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.Info("I'm alive!")

	server := api.NewApi()
	log.Fatalln(server.Start(":8000"))
}
