package app

import (
	"bot/internal/config"
	"fmt"
)

type App struct {
	Cfg *config.Config
}

func Run() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
