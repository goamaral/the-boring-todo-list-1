package main

import (
	"github.com/pkg/errors"
	m3o "go.m3o.com"

	"example.com/fiber-m3o-validator/pkg/env"
	"example.com/fiber-m3o-validator/server"
)

func main() {
	m3oClient := m3o.New(env.GetOrPanic("M3O_API_TOKEN"))

	// Run http server
	err := server.NewServer(m3oClient).Run()
	if err != nil {
		panic(errors.Wrap(err, "failed to run server"))
	}
}
