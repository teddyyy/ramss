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

func isIncludeUnit(unit string, systems []string) bool {
	for _, system := range systems {
		if system == unit {
			return true
		}
	}
	return false
}

// Get ...
func Get(systems []string) echo.HandlerFunc {
	return func(c echo.Context) error {
		if isIncludeUnit(c.Param("unit"), systems) {
			return c.JSON(http.StatusOK, systemd.Get(c.Param("unit")))
		}

		return requestError(c, http.StatusBadRequest, 400, "Invalid request")
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
			return jsonError(c, http.StatusBadRequest, err, 400, "Invalid json request")
		}

		for _, system := range systems {
			if system == c.Param("unit") {
				systemd.Post(c.Param("unit"), param.Action, param.Mode)
				return requestSuccess(c, http.StatusOK, 200, "success")
			}
		}

		return requestError(c, http.StatusBadRequest, 400, "Invalid request")
	}
}

// DefaultErrorHandler ...
func DefaultErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}

	var apierr model.APIError
	apierr.Code = code
	apierr.Message = msg

	if !c.Response().Committed {
		c.JSON(code, apierr)
	}
}

func jsonError(c echo.Context, status int, err error, code int, msg string) error {
	var apierr model.APIError
	apierr.Code = code
	apierr.Message = msg

	c.JSON(status, apierr)
	return err
}

func requestError(c echo.Context, status int, code int, msg string) error {
	var apierr model.APIError
	apierr.Code = code
	apierr.Message = msg

	return c.JSON(status, apierr)
}

func requestSuccess(c echo.Context, status int, code int, msg string) error {
	var apiscs model.APISuccess
	apiscs.Code = code
	apiscs.Message = msg

	return c.JSON(status, apiscs)
}
