package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"some-api/internal/chat"
	"some-api/internal/db"
	"sync"
	"time"
)

type Server struct {
	Router *echo.Echo
}

var connectionPool = struct {
	sync.RWMutex
	connections map[*websocket.Conn]struct{}
}{
	connections: make(map[*websocket.Conn]struct{}),
}

func NewServer() *Server {
	e := echo.New()
	a := &Server{
		Router: e,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ws", initChat)

	return a
}

func initChat(c echo.Context) error {
	ch := chat.NewChat(db.New())

	websocket.Handler(func(conn *websocket.Conn) {
		connectionPool.Lock()
		connectionPool.connections[conn] = struct{}{}

		defer func(conn2 *websocket.Conn) {
			connectionPool.Lock()
			err := conn2.Close()
			if err != nil {
				log.Warn(fmt.Sprintf("Unable to close websocket: %v", err))
			}
			delete(connectionPool.connections, conn2)
			connectionPool.Unlock()
		}(conn)

		connectionPool.Unlock()

		for {
			// Write
			msgs, _ := ch.LoadAllMessages(time.Now(), time.Now)
			if len(msgs) != 0 {
				connectionPool.RLock()

				for connection := range connectionPool.connections {
					for _, m := range msgs {
						err := websocket.Message.Send(connection, fmt.Sprintf("[%s] %s: %s", m.CreatedAt, m.UserId, m.Text))
						if err != nil {
							c.Logger().Error(err)
						}
					}
				}

				connectionPool.RUnlock()
			}

			// Read
			msg := ""
			err := websocket.Message.Receive(conn, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
			_ = ch.SaveMessage(uuid.New().String(), msg, time.Now)
		}

	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
