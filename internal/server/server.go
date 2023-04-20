package server

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"ralts/internal/chat"
	"ralts/internal/config"
)

type Server struct {
	Router      *echo.Echo
	ChatHandler *chat.Chat
	Config      *config.Config
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewServer(chatHandler *chat.Chat, cfg *config.Config) *Server {
	e := echo.New()
	s := &Server{
		Router:      e,
		ChatHandler: chatHandler,
		Config:      cfg,
	}
	pool := NewPool()
	go pool.Start()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/ws", func(c echo.Context) error {
		return s.ServeChat(c, pool)
	})

	return s
}

func (s *Server) ServeChat(c echo.Context, pool *Pool) error {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := &Connection{
		ID:   uuid.NewString(),
		C:    conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()

	return nil
}
