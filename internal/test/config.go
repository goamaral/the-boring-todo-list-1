package test

import (
	"example.com/the-boring-to-do-list-1/pkg/env"
	"example.com/the-boring-to-do-list-1/pkg/fs"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	env.SetEnvIfNotDefined("ENV", "test")
	env.SetEnvIfNotDefined("DB_HOST", "localhost")
	env.SetEnvIfNotDefined("DB_PORT", "5432")
	env.SetEnvIfNotDefined("DB_NAME", "the_boring_todo_list_1")
	env.SetEnvIfNotDefined("DB_USER", "boring")
	env.SetEnvIfNotDefined("DB_PASS", "todo")
	godotenv.Overload(fs.ResolveRelativePath("../../secrets/.env.test"))
}
