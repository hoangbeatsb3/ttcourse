package main

import (
	"gitlab.com/hoangbeatsb3/ttcourse/config"
	"gitlab.com/hoangbeatsb3/ttcourse/controller"
)

func main() {
	cfg := config.LoadEnvConfig()
	controller.Serve(cfg)
}
