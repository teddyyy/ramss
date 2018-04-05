package main

import (
	"flag"
	"io/ioutil"

	"./handlers"
	"github.com/labstack/echo"
	"gopkg.in/yaml.v2"
)

type config struct {
	Services []string `yaml:"services"`
}

var (
	fileOpt = flag.String("f", "./config.yaml", "help message for \"f\" option")
)

func main() {
	e := echo.New()

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

	e.Logger.Fatal(e.Start(":1323"))
}
