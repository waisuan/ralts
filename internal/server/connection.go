package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"ralts/internal/chat"
	"time"
)

type Payload struct {
	UserId  string
	Message string
}

type Connection struct {
	ID   string
	C    *websocket.Conn
	Pool *Pool
	Chat *chat.Chat
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

		// Limit the no. of messages sent in a day per user.
		messageCount, err := c.Chat.GetMessageCount(payload.UserId, time.Now)
		if err != nil {
			log.Errorf("unable to get message count: %e", err)
			// TODO: Send an error response back to client-side.
			continue
		}

		// TODO: Save max message count as config env var
		if messageCount == 100 {
			// TODO: Code
		}

		saved, err := c.Chat.SaveMessage(payload.UserId, payload.Message, time.Now)
		if err != nil {
			log.Errorf("unable to save message: %e", err)
			// TODO: Send an error response back to client-side.
			continue
		}

		c.Pool.Broadcast <- saved
		log.Infof("Message Received: %+v", payload)
	}
}
