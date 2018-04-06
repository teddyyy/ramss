package handlers

import (
	"net/http"

	"../model"
	"../systemd"
	"github.com/labstack/echo"
)

type (
	postParam struct {
		Action string `json:"action" validate:"required"`
		Mode   string `json:"mode" validate:"required"`
	}
)

func isIncludeConfig(unit string, systems []string) bool {
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
		if !isIncludeConfig(c.Param("unit"), systems) {
			return requestError(c, http.StatusBadRequest, 404, "Not found")
		}

		status, err := systemd.Get(c.Param("unit"))
		if err != nil {
			return httpError(c, http.StatusBadRequest, err, 400, "Invalid request")
		}

		return c.JSON(http.StatusOK, status)
	}
}

// Gets ...
func Gets(systems []string) echo.HandlerFunc {
	return func(c echo.Context) error {
		var us model.Units

		for _, system := range systems {
			status, err := systemd.Get(system)
			if err != nil {
				return httpError(c, http.StatusBadRequest, err, 400, "Invalid request")
			}
			us.Units = append(us.Units, status)
		}

		return c.JSON(http.StatusOK, us)
	}
}

// Post ...
func Post(systems []string) echo.HandlerFunc {
	return func(c echo.Context) error {
		param := new(postParam)
		if err := c.Bind(param); err != nil {
			return httpError(c, http.StatusBadRequest, err, 400, "Invalid json request")
		}

		if err := c.Validate(param); err != nil {
			return httpError(c, http.StatusBadRequest, err, 400, "Invalid json request")
		}

		if !isIncludeConfig(c.Param("unit"), systems) {
			return requestError(c, http.StatusBadRequest, 404, "Not found")
		}

		if err := systemd.Post(c.Param("unit"), param.Action, param.Mode); err != nil {
			return httpError(c, http.StatusBadRequest, err, 400, "Invalid request")
		}

		return postSuccess(c, http.StatusOK, 200, "success")
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

func httpError(c echo.Context, status int, err error, code int, msg string) error {
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

func postSuccess(c echo.Context, status int, code int, msg string) error {
	var apiscs model.APISuccess
	apiscs.Code = code
	apiscs.Message = msg

	return c.JSON(status, apiscs)
}
