package repo

import (
	"context"
	"time"
)

// ToDoRepo for both memory and mysql
type ToDoRepo interface {
	Create(context.Context, *ToDo) (ToDoID, error)
	Read(context.Context, ToDoID) (*ToDo, error)
	ReadAll(context.Context) ([]*ToDo, error)
	Update(context.Context, *ToDo) error
	Delete(context.Context, ToDoID) error
}

// ToDoID for ToDo object which clarify ToDoRepo func parameter
type ToDoID int64

// Int64 from ToDoID
func (id ToDoID) Int64() int64 {
	return int64(id)
}

// ToDo only used in repo
type ToDo struct {
	ID          int64
	Title       string
	Description string
	Reminder    time.Time
}
