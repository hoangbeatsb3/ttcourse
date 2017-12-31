package main

import (
	"github.com/hoangbeatsb3/ttcourse/config"
	"github.com/hoangbeatsb3/ttcourse/controller"
)

func main() {
	cfg := config.LoadEnvConfig()
	controller.Serve(cfg)
}
