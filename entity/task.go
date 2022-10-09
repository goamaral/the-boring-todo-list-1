package entity

import (
	"time"

	"example.com/fiber-m3o-validator/errors"
)

type Task struct {
	AbstractEntity

	Title      string
	CompleteAt *time.Time
}

/* PUBLIC */
const TasksTableName = "tasks"

func TaskToMap(t Task) map[string]interface{} {
	return map[string]interface{}{
		"id":          t.Id,
		"created_at":  t.CreatedAt,
		"title":       t.Title,
		"complete_at": t.CompleteAt,
	}
}

func TaskFromMap(taskMap map[string]interface{}) (Task, error) {
	var task Task
	var ok bool

	task.Id, ok = taskMap["id"].(string)
	if !ok {
		return Task{}, errors.NewParseError("task id", taskMap["id"])
	}

	// TODO: created_at
	// TODO: title
	// TODO: completed_at

	return task, nil
}
