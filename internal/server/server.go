package server

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"some-api/internal/chat"
	"some-api/internal/db"
	"sync"
	"syscall"
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
	e.Use(middleware.CORS())

	e.GET("/ws", initChat)

	return a
}

func removeWsConn(conn *websocket.Conn) {
	connectionPool.Lock()
	err := conn.Close()
	if err != nil {
		log.Warn(fmt.Sprintf("Unable to close websocket: %v", err))
	}
	log.Info(fmt.Sprintf(">>> Removing connection from connection pool (%d)...", len(connectionPool.connections)))
	delete(connectionPool.connections, conn)
	log.Info(fmt.Sprintf(">>> Connection pool size is now: %d", len(connectionPool.connections)))
	connectionPool.Unlock()
}

func initChat(c echo.Context) error {
	ch := chat.NewChat(db.New())

	websocket.Handler(func(conn *websocket.Conn) {
		connectionPool.Lock()
		connectionPool.connections[conn] = struct{}{}

		defer removeWsConn(conn)

		connectionPool.Unlock()

		forceDisconnect := false
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

							// Broken pipe, conn is probably dead.
							if errors.Is(err, syscall.EPIPE) {
								forceDisconnect = true
								break
							}
						}
					}
				}

				connectionPool.RUnlock()
			}

			if forceDisconnect {
				break
			}

			// Read
			msg := ""
			err := websocket.Message.Receive(conn, &msg)
			if err != nil {
				c.Logger().Error(err)

				// Disconnect initiated from caller.
				if err.Error() == "EOF" {
					break
				}
			} else {
				fmt.Printf("%s\n", msg)
				_ = ch.SaveMessage(uuid.New().String(), msg, time.Now)
			}
		}

	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
