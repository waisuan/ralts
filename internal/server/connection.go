package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Payload struct {
	UserId  string
	Message string
}

type Connection struct {
	ID   string
	C    *websocket.Conn
	Pool *Pool
}

func (c *Connection) Read() {
	defer func() {
		c.Pool.Unregister <- c
		_ = c.C.Close()
	}()

	for {
		_, p, err := c.C.ReadMessage()
		if err != nil {
			log.Errorf("unable to read message: %e", err)
			if websocket.IsCloseError(err) {
				break
			}
			if websocket.IsUnexpectedCloseError(err) {
				break
			}
			continue
		}
		var payload Payload
		err = json.Unmarshal(p, &payload)
		if err != nil {
			log.Errorf("unable to unmarshal message: %e", err)
			continue
		}

		c.Pool.Broadcast <- payload
		log.Infof("Message Received: %+v", payload)
	}
}
