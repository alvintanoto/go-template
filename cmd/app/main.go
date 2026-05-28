package main

import (
	"fmt"

	"alvintanoto.id/go-template/internal/config"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	// import config
	c.Provide(config.NewConfig)

	return c
}

func main() {
	container := BuildContainer()
	err := container.Invoke(func(cfg *config.Config) error {
		fmt.Printf("%s App running in port :%s\n", cfg.AppEnv, cfg.Port)
		return nil
	})

	if err != nil {
		panic(err)
	}
}
