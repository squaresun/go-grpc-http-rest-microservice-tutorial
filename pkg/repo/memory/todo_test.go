package memory

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/squaresun/go-grpc-http-rest-microservice-tutorial/pkg/repo"
)

func TestCRUD(t *testing.T) {
	ctx := context.Background()

	db, err := NewToDoDB()
	if err != nil {
		t.Error(err)
	}

	id, err := db.Create(ctx, &repo.ToDo{})
	if err != nil {
		t.Error(err)
	}

	todo := &repo.ToDo{
		ID:          int64(id),
		Title:       "hello",
		Description: "world",
		Reminder:    time.Now(),
	}
	err = db.Update(ctx, todo)
	if err != nil {
		t.Error(err)
	}

	todo, err = db.Read(ctx, id)
	if err != nil {
		t.Error(err)
	}

	todos, err := db.ReadAll(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(todos) != 1 {
		t.Errorf("Error length of todos: %d", len(todos))
		return
	}
	if todos[0].ID != todo.ID {
		t.Errorf("ID not match between Read: %d and ReadAll: %d", todos[0].ID, todo.ID)
	}

	err = db.Delete(ctx, id)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Read(ctx, id)
	if err != memdb.ErrNotFound {
		t.Errorf("ID: %d is not deleted", id)
	}
}
