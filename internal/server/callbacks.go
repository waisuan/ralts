package server

import (
	log "github.com/sirupsen/logrus"
	"ralts/internal/dependencies"
)

const CONN_COUNT_KEY = "num_of_conn"

type Callbacks struct {
	Deps           *dependencies.Dependencies
	PostRegister   chan bool
	PostUnregister chan bool
}

func NewCallbacks(deps *dependencies.Dependencies) *Callbacks {
	return &Callbacks{
		Deps:           deps,
		PostRegister:   make(chan bool),
		PostUnregister: make(chan bool),
	}
}

func (c *Callbacks) Listen() {
	for {
		select {
		case <-c.PostRegister:
			err := c.Deps.Cache.Incr(CONN_COUNT_KEY)
			if err != nil {
				log.Errorf("unable to handle PostRegister callback: %e", err)
			}

			break
		case <-c.PostUnregister:
			err := c.Deps.Cache.Decr(CONN_COUNT_KEY)
			if err != nil {
				log.Errorf("unable to handle PostUnregister callback: %e", err)
			}

			break
		}
	}
}
