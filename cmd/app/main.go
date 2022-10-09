package main

import (
	m3o "go.m3o.com"

	"example.com/fiber-m3o-validator/internal/server"
	"example.com/fiber-m3o-validator/internal/service"
	"example.com/fiber-m3o-validator/pkg/env"
)

func main() {
	m3oClient := m3o.New(env.GetOrPanic("M3O_API_TOKEN"))
	taskService := service.NewTaskService(m3oClient.Db)

	// Run http server
	err := server.NewServer(taskService).Run()
	if err != nil {
		panic(err)
	}
}
