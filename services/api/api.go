package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"some-api/services/location"
	"some-api/utils/db"
)

type Api struct {
	db db.DataStore
	echo *echo.Echo
}

func NewApi(dbClient db.DataStore) *Api {
	e := echo.New()
	a := &Api{
		db: dbClient,
		echo: e,
	}
	e.GET("/location/:id", a.getLocation)
	return a
}

func (a *Api) getLocation(c echo.Context) error {
	personId := c.Param("id")
	out, err := location.Get(a.db, personId)
	if err != nil {
		c.Error(err)
	}
	return c.JSON(http.StatusOK, out)
}
