package main

import (
	"flag"
	"io/ioutil"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/ramss/handlers"
	"gopkg.in/yaml.v2"
)

type config struct {
	Services []string `yaml:"services"`
}

var (
	fileOpt = flag.String("f", "./config.yaml", "config file")
	portOpt = flag.String("p", "8080", "listen port")
)

func main() {
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}

	flag.Parse()

	file, err := ioutil.ReadFile(*fileOpt)
	if err != nil {
		panic(err)
	}

	c := &config{}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		panic(err)
	}

	// Error
	e.HTTPErrorHandler = handlers.DefaultErrorHandler

	// Routing
	e.GET("/api/v1/systemd/", handlers.Gets(c.Services))
	e.GET("/api/v1/systemd/:unit", handlers.Get(c.Services))
	e.POST("/api/v1/systemd/:unit", handlers.Post(c.Services))

	e.Logger.Fatal(e.Start(":" + *portOpt))
}
