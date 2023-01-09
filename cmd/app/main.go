package main

import (
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/internal/server"
	gormprovider "example.com/the-boring-to-do-list-1/pkg/provider/gorm"
)

func main() {
	// Init gorm provider
	gormProvider, err := gormprovider.NewProvider()
	if err != nil {
		panic(err)
	}

	// Initialize repositories
	taskRepo := repository.NewTaskRepository(gormProvider)

	// Run http server
	err = server.NewServer(taskRepo).Run()
	if err != nil {
		panic(err)
	}
}
