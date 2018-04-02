package handlers

import (
	"net/http"

	"../systemd"
	"github.com/labstack/echo"
)

// PostParam ...
type PostParam struct {
	Action string `json:"action"`
}

// Get ...
func Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, systemd.Get(c.Param("unit")))
	}
}

// Post ...
func Post() echo.HandlerFunc {
	return func(c echo.Context) error {
		param := new(PostParam)
		if err := c.Bind(param); err != nil {
			return err
		}

		systemd.Post(c.Param("unit"), param.Action)

		return c.JSON(http.StatusOK, "")
	}
}

// Delete ...
func Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		systemd.Delete(c.Param("unit"))

		return c.JSON(http.StatusOK, "")
	}
}

// Put ...
func Put() echo.HandlerFunc {
	return func(c echo.Context) error {
		systemd.Put(c.Param("unit"))

		return c.JSON(http.StatusOK, "")
	}
}
