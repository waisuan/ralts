package main

import log "github.com/sirupsen/logrus"

func main() {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.Info("I'm alive!")
}
