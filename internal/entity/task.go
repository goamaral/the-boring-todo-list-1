package entity

import (
	"time"
)

type Task struct {
	AbstractEntity

	Title       string
	CompletedAt *time.Time
}

/* PUBLIC */
const TasksTableName = "tasks"

func TaskToMap(t Task) map[string]interface{} {
	return map[string]interface{}{
		"id":           t.Id,
		"created_at":   t.CreatedAt.Format(time.RFC3339),
		"title":        t.Title,
		"completed_at": t.CompletedAt,
	}
}

func TaskFromMap(taskMap map[string]interface{}) (Task, error) {
	var task Task
	var err error

	task.AbstractEntity, err = AbstractEntityFromMap(taskMap)
	if err != nil {
		return Task{}, err
	}

	task.Title = taskMap["title"].(string)

	completedAtStr, notNil := taskMap["completed_at"].(string)
	if notNil {
		completedAt, err := time.Parse(time.RFC3339, completedAtStr)
		if err != nil {
			return Task{}, err
		}
		task.CompletedAt = &completedAt
	}

	return task, nil
}
