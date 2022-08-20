package entity

import "time"

type Task struct {
	abstractEntity

	Title      string
	CompleteAt *time.Time
}

/* PUBLIC */
func (t Task) GetTableName() string {
	return "tasks"
}

func (t Task) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":          t.Id,
		"created_at":  t.CreatedAt,
		"title":       t.Title,
		"complete_at": t.CompleteAt,
	}
}
