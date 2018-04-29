package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/teddyyy/ramss/model"
	"github.com/teddyyy/ramss/systemd"
)

type (
	postParam struct {
		Action string `json:"action" validate:"required"`
		Mode   string `json:"mode" validate:"required"`
	}
)

func isIncludeConfig(unit string, systems []model.Service) *model.Service {
	for _, system := range systems {
		if system.UnitName == unit {
			return &system
		}
	}
	return nil
}

// Get ...
func Get(systems []model.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		system := isIncludeConfig(c.Param("unit"), systems)
		if system == nil {
			return requestError(c, http.StatusBadRequest, 404, "Not found")
		}

		status, err := systemd.Get(system.UnitName, system.ServiceName)
		if err != nil {
			return httpError(c, http.StatusBadRequest, err, 400, "Invalid request")
		}

		return c.JSON(http.StatusOK, status)
	}
}

// Gets ...
func Gets(systems []model.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var us model.Units

		for _, system := range systems {
			status, err := systemd.Get(system.UnitName, system.ServiceName)
			if err != nil {
				return httpError(c, http.StatusBadRequest, err, 400, "Invalid request")
			}
			us.Units = append(us.Units, status)
		}

		return c.JSON(http.StatusOK, us)
	}
}

// Post ...
func Post(systems []model.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		param := new(postParam)
		if err := c.Bind(param); err != nil {
			return httpError(c, http.StatusBadRequest, err, 400, "Invalid json request")
		}

		if err := c.Validate(param); err != nil {
			return httpError(c, http.StatusBadRequest, err, 400, "Invalid json request")
		}

		system := isIncludeConfig(c.Param("unit"), systems)
		if system == nil {
			return requestError(c, http.StatusBadRequest, 404, "Not found")
		}

		if err := systemd.Post(system.ServiceName, param.Action, param.Mode); err != nil {
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
