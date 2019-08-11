package memory

import (
	"context"
	"sync/atomic"

	"github.com/hashicorp/go-memdb"
	"github.com/squaresun/go-grpc-http-rest-microservice-tutorial/pkg/repo"
)

const (
	todoTableName  = "todo"
	todoTableIndex = "id"
)

type todoDB struct {
	*memdb.MemDB
	id *int64
}

// NewToDoDB returns a memory-based ToDo database
func NewToDoDB() (repo.ToDoRepo, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			todoTableName: &memdb.TableSchema{
				Name: todoTableName,
				Indexes: map[string]*memdb.IndexSchema{
					todoTableIndex: &memdb.IndexSchema{
						Name:    todoTableIndex,
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "ID"},
					},
				},
			},
		},
	}

	id := int64(0)

	db, err := memdb.NewMemDB(schema)
	return &todoDB{
		MemDB: db,
		id:    &id,
	}, err
}

func (r *todoDB) Create(_ context.Context, t *repo.ToDo) (repo.ToDoID, error) {
	t.ID = atomic.AddInt64(r.id, 1)

	txn := r.Txn(true)

	err := txn.Insert(todoTableName, t)
	if err != nil {
		return 0, err
	}

	txn.Commit()
	return repo.ToDoID(t.ID), nil
}

func (r *todoDB) Read(_ context.Context, id repo.ToDoID) (*repo.ToDo, error) {
	txn := r.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(todoTableName, todoTableIndex, id)
	if err != nil {
		return nil, err
	}

	if todo := iter.Next(); todo != nil {
		return todo.(*repo.ToDo), nil
	}
	return nil, memdb.ErrNotFound
}

func (r *todoDB) ReadAll(_ context.Context) ([]*repo.ToDo, error) {
	txn := r.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(todoTableName, todoTableIndex)
	if err != nil {
		return nil, err
	}

	todos := []*repo.ToDo{}
	for todo := iter.Next(); todo != nil; todo = iter.Next() {
		todos = append(todos, todo.(*repo.ToDo))
	}

	return todos, nil
}

func (r *todoDB) Update(_ context.Context, t *repo.ToDo) error {
	txn := r.Txn(true)

	err := txn.Insert(todoTableName, t)
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (r *todoDB) Delete(_ context.Context, id repo.ToDoID) error {
	txn := r.Txn(true)

	iter, err := txn.Get(todoTableName, todoTableIndex, id)
	if err != nil {
		return err
	}

	todo := iter.Next()
	if todo == nil {
		return memdb.ErrNotFound
	}
	err = txn.Delete(todoTableName, todo.(*repo.ToDo))

	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}
