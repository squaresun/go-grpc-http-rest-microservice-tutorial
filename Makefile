proto:
	protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=plugins=grpc:pkg/api/v1 todo-service.proto
	protoc --proto_path=api/proto/v1 --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api/v1 todo-service.proto
	protoc --proto_path=api/proto/v1 --proto_path=third_party --swagger_out=logtostderr=true:api/swagger/v1 todo-service.proto

server:
	go build ./cmd/server
	# Assumed that there is root:password@tcp(localhost:3306)/grpc
	./server -grpc-port=9090 -http-port=8080 -db-host=localhost:3306 -db-user=root -db-password=password -db-schema=grpc -log-level=-1 -log-time-format=2006-01-02T15:04:05.999999999Z07:00

client-grpc:
	go build ./cmd/client-grpc
	./client-grpc -server=localhost:9090

client-rest:
	go build ./cmd/client-rest
	./client-rest -server=http://localhost:8080
