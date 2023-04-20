package server

import (
	log "github.com/sirupsen/logrus"
)

type Pool struct {
	Register   chan *Connection
	Unregister chan *Connection
	Clients    map[*Connection]bool
	Broadcast  chan Payload
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Connection),
		Unregister: make(chan *Connection),
		Clients:    make(map[*Connection]bool),
		Broadcast:  make(chan Payload),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			log.Infof("Size of Connection Pool: %d", len(pool.Clients))
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			log.Infof("Size of Connection Pool: %d", len(pool.Clients))
			break
		case message := <-pool.Broadcast:
			log.Info("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.C.WriteJSON(message); err != nil {
					log.Errorf("unable to broadcast message: %e", err)
					return
				}
			}
		}
	}
}
