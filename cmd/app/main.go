package main

import (
	"example.com/fiber-m3o-validator/server"
	"github.com/pkg/errors"
)

func main() {
	// TODO: Load envs
	// TODO: Init M3O provider

	// Run http server
	err := server.NewServer().Run()
	if err != nil {
		panic(errors.Wrap(err, "failed to run server"))
	}
}
