package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"net/http"
	"ralts/internal/chat"
	"ralts/internal/config"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	Router      *echo.Echo
	ChatHandler *chat.Chat
	Config      *config.Config
}

type Request struct {
	UserId  string
	Message string
}

var connectionPool = struct {
	sync.RWMutex
	connections map[*websocket.Conn]struct{}
}{
	connections: make(map[*websocket.Conn]struct{}),
}

func NewServer(chatHandler *chat.Chat, cfg *config.Config) *Server {
	e := echo.New()
	a := &Server{
		Router:      e,
		ChatHandler: chatHandler,
		Config:      cfg,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/ws", a.initChat)

	return a
}

func (s *Server) removeWsConn(conn *websocket.Conn) {
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

func (s *Server) initChat(c echo.Context) error {
	token := c.QueryParam("authorization")
	if token != s.Config.AuthToken {
		return c.JSON(http.StatusUnauthorized, "not authenticated")
	}

	websocket.Handler(func(conn *websocket.Conn) {
		connectionPool.Lock()
		connectionPool.connections[conn] = struct{}{}

		defer s.removeWsConn(conn)

		connectionPool.Unlock()

		var latestMsg *chat.Message
		forceDisconnect := false
		for {
			// Write
			if latestMsg == nil {
				msgs, _ := s.ChatHandler.LoadAllMessages()
				if len(msgs) != 0 {
					connectionPool.RLock()

					for connection := range connectionPool.connections {
						for _, m := range msgs {
							payload, err := json.Marshal(&m)
							if err != nil {
								c.Logger().Error(err)
							} else {
								err := websocket.Message.Send(connection, string(payload))
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
					}

					connectionPool.RUnlock()
				}
			} else {
				connectionPool.RLock()

				for connection := range connectionPool.connections {
					payload, err := json.Marshal(&latestMsg)
					if err != nil {
						c.Logger().Error(err)
					} else {
						err := websocket.Message.Send(connection, string(payload))
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
				var req Request
				err := json.Unmarshal([]byte(msg), &req)
				if err != nil {
					c.Logger().Error(err)
				} else {
					fmt.Printf("%s: %s\n", req.UserId, req.Message)
					latestMsg, _ = s.ChatHandler.SaveMessage(req.UserId, req.Message, time.Now)
				}
			}
		}

	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
