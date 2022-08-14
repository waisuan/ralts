package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func NewApi() *echo.Echo{
	e := echo.New()
	e.GET("/location/:id", getLocation)
	return e
}

func getLocation(c echo.Context) error {
	personId := c.Param("id")
	return c.String(http.StatusOK, personId)
}
