package server

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"some-api/internal/location"
	"some-api/utils/db"
)

type Server struct {
	dataStore db.DataStore
	Router    *echo.Echo
}

func NewServer(dbClient db.DataStore) *Server {
	e := echo.New()
	a := &Server{
		dataStore: dbClient,
		Router:    e,
	}

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	e.GET("/location/:id", a.getLocation)
	return a
}

func (s *Server) getLocation(c echo.Context) error {
	personId := c.Param("id")
	out, err := location.Get(s.dataStore, personId)
	if err != nil {
		log.Error(err)
		c.Error(err)
	}

	if out == nil {
		return c.String(http.StatusNotFound, "Not found")
	}

	return c.JSON(http.StatusOK, out)
}
