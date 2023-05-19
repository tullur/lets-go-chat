BINARY_NAME=lets-go-chat

build:
	go build -o bin/${BINARY_NAME} cmd/lets-go-chat/main.go

run:
	go build -o bin/${BINARY_NAME} cmd/lets-go-chat/main.go
	bin/${BINARY_NAME}

test:
	go test ./...

clean:
	go clean
	rm bin/${BINARY_NAME}