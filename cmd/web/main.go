package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"ralts/internal/config"
	"ralts/internal/dependencies"
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
	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	s := server.NewServer(deps)
	log.Fatalln(s.Router.Start(":8001"))
}
