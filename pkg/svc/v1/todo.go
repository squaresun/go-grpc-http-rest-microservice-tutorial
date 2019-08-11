package v1

import (
	"context"

	"github.com/squaresun/go-grpc-http-rest-microservice-tutorial/pkg/repo"

	"github.com/golang/protobuf/ptypes"
	v1 "github.com/squaresun/go-grpc-http-rest-microservice-tutorial/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type toDoServiceServer struct {
	repo.ToDoRepo
}

// NewToDoServiceServer creates ToDo service
func NewToDoServiceServer(r repo.ToDoRepo) v1.ToDoServiceServer {
	return &toDoServiceServer{ToDoRepo: r}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *toDoServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// Create new todo task
func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	// insert ToDo entity data
	id, err := s.ToDoRepo.Create(ctx, &repo.ToDo{
		Title:       req.ToDo.Title,
		Description: req.ToDo.Description,
		Reminder:    reminder,
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into ToDo-> "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id.Int64(),
	}, nil
}

// Read todo task
func (s *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// query ToDo by ID
	todo, err := s.ToDoRepo.Read(ctx, repo.ToDoID(req.Id))
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}

	v1ToDo := v1.ToDo{
		Id:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
	}
	v1ToDo.Reminder, err = ptypes.TimestampProto(todo.Reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
	}

	return &v1.ReadResponse{
		Api:  apiVersion,
		ToDo: &v1ToDo,
	}, nil

}

// Update todo task
func (s *toDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	// update ToDo
	err = s.ToDoRepo.Update(ctx, &repo.ToDo{
		Title:       req.ToDo.Title,
		Description: req.ToDo.Description,
		Reminder:    reminder,
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update ToDo-> "+err.Error())
	}

	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: 1,
	}, nil
}

// Delete todo task
func (s *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// delete ToDo
	err := s.ToDoRepo.Delete(ctx, repo.ToDoID(req.Id))
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete ToDo-> "+err.Error())
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: 1,
	}, nil
}

// Read all todo tasks
func (s *toDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get ToDo list
	todos, err := s.ToDoRepo.ReadAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}

	list := make([]*v1.ToDo, len(todos))
	for i, todo := range todos {

		v1ToDo := v1.ToDo{
			Id:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
		}
		v1ToDo.Reminder, err = ptypes.TimestampProto(todo.Reminder)
		if err != nil {
			return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
		}

		list[i] = &v1ToDo
	}

	return &v1.ReadAllResponse{
		Api:   apiVersion,
		ToDos: list,
	}, nil
}
