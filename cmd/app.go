package main

import (
	"github.com/ThEditor/clutter/internal/api"
	"github.com/ThEditor/clutter/internal/config"
)

func main() {
	cfg := config.Load()

	api.Start(cfg.BIND_ADDRESS, cfg.PORT)
}
