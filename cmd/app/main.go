package main

import (
	"os"

	"example.com/the-boring-to-do-list-1/internal/initializer"
	"example.com/the-boring-to-do-list-1/internal/server"
	"example.com/the-boring-to-do-list-1/pkg/fs"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
)

func main() {
	// Gorm
	gorm_provider, err := initializer.NewProvider(initializer.DefaultDSN())
	if err != nil {
		panic(err)
	}

	// JWT
	privKeyFile, err := os.Open(fs.ResolveRelativePath("../../secrets/ecdsa"))
	if err != nil {
		panic(err)
	}
	pubKeyFile, err := os.Open(fs.ResolveRelativePath("../../secrets/ecdsa.pub"))
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
