package main

import (
	"os"
	"path/filepath"
	"runtime"

	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/internal/server"
	"example.com/the-boring-to-do-list-1/pkg/gormprovider"
	"example.com/the-boring-to-do-list-1/pkg/jwtprovider"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	folderPath := filepath.Dir(b)

	/* PROVIDERS */
	// Gorm
	gormProvider, err := gormprovider.NewProvider()
	if err != nil {
		panic(err)
	}

	// JWT
	privKeyFile, err := os.Open(folderPath + "/../../secrets/ecdsa")
	if err != nil {
		panic(err)
	}
	pubKeyFile, err := os.Open(folderPath + "/../../secrets/ecdsa.pub")
	if err != nil {
		panic(err)
	}
	jwtProvider, err := jwtprovider.NewProvider(privKeyFile, pubKeyFile)
	if err != nil {
		panic(err)
	}

	// Repositories
	taskRepo := repository.NewTaskRepository(gormProvider)
	userRepo := repository.NewUserRepository(gormProvider)

	// HTTP server
	err = server.NewServer(jwtProvider, taskRepo, userRepo).Run()
	if err != nil {
		panic(err)
	}
}
