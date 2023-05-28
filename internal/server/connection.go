package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"ralts/internal/chat"
	"ralts/internal/dependencies"
	"time"
)

type Connection struct {
	ID   string
	C    *websocket.Conn
	Pool *Pool
	Chat *chat.Chat
	Deps *dependencies.Dependencies
}

func (c *Connection) Read() {
	defer func() {
		c.Pool.Unregister <- c
		_ = c.C.Close()
	}()

	for {
		_, r, err := c.C.ReadMessage()
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
		var request Request
		err = json.Unmarshal(r, &request)
		if err != nil {
			log.Errorf("unable to unmarshal message: %e", err)
			c.respondWithError(err)
			continue
		}

		// Limit the no. of messages sent in a day per user.
		messageCount, err := c.Chat.GetMessageCount(request.UserId, time.Now)
		if err != nil {
			log.Errorf("unable to get message count: %e", err)
			c.respondWithError(err)
			continue
		}

		if messageCount >= c.Deps.Cfg.MaxSentMsgPerDay {
			log.Warnf("max message sent limit (%d) reached for %s", messageCount, request.UserId)
			_ = c.C.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "reached max no. of messages sent today"),
			)
			break
		}

		saved, err := c.Chat.SaveMessage(request.UserId, request.Message, time.Now)
		if err != nil {
			log.Errorf("unable to save message: %e", err)
			c.respondWithError(err)
			continue
		}

		resp := &Response{
			Payload: saved,
		}
		c.Pool.Broadcast <- resp
		log.Infof("Message Received: %+v", request)
	}
}

func (c *Connection) respondWithError(err error) {
	r := Response{
		Error: &Error{
			Code:    InternalServerError,
			Message: err.Error(),
		},
	}
	if err := c.C.WriteJSON(r); err != nil {
		log.Errorf("unable to send error response: %e", err)
	}
}
