gen:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/laptop_service.proto proto/laptop_message.proto proto/screen_message.proto proto/processor_message.proto proto/memory_message.proto proto/storage_message.proto proto/keyboard_message.proto
clean:
	rm pb/*.go
server:
	go run cmd/server/main.go -port 8080
client:
	go run cmd/client/main.go -address 0.0.0.0:8080
test:
	go test -cover -race ./...