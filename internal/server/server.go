package server

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ralts/internal/chat"
	"ralts/internal/dependencies"
	"ralts/internal/newsfeed"
	"strconv"
)

type Server struct {
	Router *echo.Echo
	Deps   *dependencies.Dependencies
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewServer(deps *dependencies.Dependencies) *Server {
	e := echo.New()
	s := &Server{
		Router: e,
		Deps:   deps,
	}

	callbacks := NewCallbacks(deps)
	go callbacks.Listen()

	pool := NewPool(callbacks)
	go pool.Start()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	//middleware.BasicAuth()

	e.GET("/ws", func(c echo.Context) error {
		return s.ServeChat(c, pool)
	})
	e.GET("/conn_count", s.GetConnCount)
	e.GET("/news-feed", func(c echo.Context) error {
		return s.GetNewsFeed(c, newsfeed.NewNewsFeed(deps))
	})

	return s
}

func (s *Server) ServeChat(c echo.Context, pool *Pool) error {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Errorf("unable to initiate websocket request: %e", err)
		return err
	}
	defer conn.Close()

	if len(pool.Clients) >= s.Deps.Cfg.MaxConnCount {
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "max no. of client connections reached"),
		)
	} else {
		client := &Connection{
			ID:   uuid.NewString(),
			C:    conn,
			Pool: pool,
			Chat: chat.NewChat(s.Deps),
			Deps: s.Deps,
		}

		pool.Register <- client
		client.Read()
	}

	return nil
}

func (s *Server) GetConnCount(c echo.Context) error {
	type response struct {
		Count int `json:"count"`
	}

	v, err := s.Deps.Cache.Get(CONN_COUNT_KEY)
	if err != nil {
		log.Errorf("unable to fetch conn count from cache: %e", err)
		return c.JSON(http.StatusInternalServerError, "unable to handle request")
	}

	vi, _ := strconv.Atoi(v)

	res := &response{
		Count: vi,
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) GetNewsFeed(c echo.Context, newsFeedHandler newsfeed.NewsFeedHandler) error {
	a, err := newsFeedHandler.LoadAllArticles()
	if err != nil {
		log.Errorf("unable to load news feed articles: %e", err)
		return c.JSON(http.StatusInternalServerError, "unable to handle request")
	}

	return c.JSON(http.StatusOK, a)
}
