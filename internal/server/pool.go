package server

import (
	log "github.com/sirupsen/logrus"
	"ralts/internal/chat"
)

type Pool struct {
	Register   chan *Connection
	Unregister chan *Connection
	Clients    map[*Connection]bool
	Broadcast  chan *chat.Message
	Callbacks  *Callbacks
}

func NewPool(callbacks *Callbacks) *Pool {
	return &Pool{
		Register:   make(chan *Connection),
		Unregister: make(chan *Connection),
		Clients:    make(map[*Connection]bool),
		Broadcast:  make(chan *chat.Message),
		Callbacks:  callbacks,
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			log.Infof("Size of Connection Pool: %d", len(pool.Clients))

			msgs, err := client.Chat.LoadAllMessages()
			if err != nil {
				log.Errorf("unable to send recent messages to new client: %e", err)
			} else {
				for _, msg := range msgs {
					if err := client.C.WriteJSON(msg); err != nil {
						log.Errorf("unable to send message to new client: %e", err)
					}
				}
			}

			pool.Callbacks.PostRegister <- true

			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			log.Infof("Size of Connection Pool: %d", len(pool.Clients))

			pool.Callbacks.PostUnregister <- true

			break
		case message := <-pool.Broadcast:
			log.Info("Sending message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.C.WriteJSON(message); err != nil {
					log.Errorf("unable to broadcast message: %e", err)
					return
				}
			}
		}
	}
}
