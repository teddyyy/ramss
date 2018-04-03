package handlers

import (
	"net/http"

	"../model"
	"../systemd"
	"github.com/labstack/echo"
)

// PostParam ...
type PostParam struct {
	Action string `json:"action"`
	Mode   string `json:"mode"`
}

// Get ...
func Get(systems []string) echo.HandlerFunc {
	return func(c echo.Context) error {
		for _, system := range systems {
			if system == c.Param("unit") {
				return c.JSON(http.StatusOK, systemd.Get(c.Param("unit")))
			}
		}

		return c.JSON(http.StatusBadRequest, "error")
	}
}

// Gets ...
func Gets(systems []string) echo.HandlerFunc {
	return func(c echo.Context) error {
		var us model.Units
		for _, system := range systems {
			u := systemd.Get(system)
			us.Units = append(us.Units, u)
		}

		return c.JSON(http.StatusOK, us)
	}
}

// Post ...
func Post(systems []string) echo.HandlerFunc {
	return func(c echo.Context) error {
		param := new(PostParam)
		if err := c.Bind(param); err != nil {
			return err
		}

		for _, system := range systems {
			if system == c.Param("unit") {
				systemd.Post(c.Param("unit"), param.Action, param.Mode)
				return c.JSON(http.StatusOK, "ok")
			}
		}

		return c.JSON(http.StatusBadRequest, "error")
	}
}
