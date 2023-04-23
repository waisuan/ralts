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
	Router   *echo.Echo
	Config   *config.Config
	Handlers *Handlers
}

type Handlers struct {
	ChatHandler *chat.Chat
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewServer(handlers *Handlers, cfg *config.Config) *Server {
	e := echo.New()
	s := &Server{
		Router:   e,
		Config:   cfg,
		Handlers: handlers,
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

	if len(pool.Clients) >= s.Config.MaxConnCount {
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "max no. of client connections reached"),
		)
	} else {
		client := &Connection{
			ID:   uuid.NewString(),
			C:    conn,
			Pool: pool,
			Chat: s.Handlers.ChatHandler,
		}

		pool.Register <- client
		client.Read()
	}

	return nil
}
