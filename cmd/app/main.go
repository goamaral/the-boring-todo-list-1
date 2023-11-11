package main

import (
	"os"
	"path/filepath"
	"runtime"

	"example.com/the-boring-to-do-list-1/internal/server"
	postgres_gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider/postgres"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	folderPath := filepath.Dir(b)

	/* PROVIDERS */
	// Gorm
	gorm_provider, err := postgres_gorm_provider.NewProvider(postgres_gorm_provider.DefaultDSN())
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
	jwt_provider, err := jwt_provider.NewProvider(privKeyFile, pubKeyFile)
	if err != nil {
		panic(err)
	}

	// HTTP server
	err = server.NewServer(jwt_provider, gorm_provider).Run()
	if err != nil {
		panic(err)
	}
}
