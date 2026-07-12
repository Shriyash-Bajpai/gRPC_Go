gen:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/laptop_service.proto proto/laptop_message.proto proto/screen_message.proto proto/processor_message.proto proto/memory_message.proto proto/storage_message.proto proto/keyboard_message.proto
clean:
	rm pb/*.go
run:
	go run Main.go
test:
	go test -cover -race ./...