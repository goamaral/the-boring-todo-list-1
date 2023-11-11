package config

import (
	"path/filepath"
	"runtime"

	"example.com/the-boring-to-do-list-1/pkg/env"
	"github.com/joho/godotenv"
)

func RelativePath(relativePath string) string {
	_, file, _, _ := runtime.Caller(1)
	folderPath := filepath.Dir(file)
	return folderPath + "/" + relativePath
}

func LoadTestEnv() {
	env.SetEnvIfNotDefined("ENV", "TEST")
	env.SetEnvIfNotDefined("DB_HOST", "localhost")
	env.SetEnvIfNotDefined("DB_PORT", "5432")
	env.SetEnvIfNotDefined("DB_NAME", "the_boring_todo_list_1")
	env.SetEnvIfNotDefined("DB_USER", "boring")
	env.SetEnvIfNotDefined("DB_PASS", "todo")
	godotenv.Overload(RelativePath("../../secrets/.env.test"))
}
