.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

proto: ## Convert protobuf to golang code
	protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=plugins=grpc:pkg/api/v1 todo-service.proto
	protoc --proto_path=api/proto/v1 --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api/v1 todo-service.proto
	protoc --proto_path=api/proto/v1 --proto_path=third_party --swagger_out=logtostderr=true:api/swagger/v1 todo-service.proto

docker-mariadb: ## Docker run a mariadb service for this server
	docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=grpc mariadb

server: ## Build and run server
	go build ./cmd/server
	# Assumed that there is root:password@tcp(localhost:3306)/grpc
	./server -grpc-port=9090 -http-port=8080 -db-host=localhost:3306 -db-user=root -db-password=password -db-schema=grpc -log-level=-1 -log-time-format=2006-01-02T15:04:05.999999999Z07:00

client-grpc: ## Build and run grpc client
	go build ./cmd/client-grpc
	./client-grpc -server=localhost:9090

client-rest: ## Build and run rest client
	go build ./cmd/client-rest
	./client-rest -server=http://localhost:8080

depgraph: ## Show dependencies of server in this repo
	godepgraph -s -o github.com/squaresun/go-grpc-http-rest-microservice-tutorial ./cmd/server | dot -Tpng -o godepgraph.png
	open godepgraph.png
