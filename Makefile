BINARY_NAME=lets-go-chat

build:
	go build -o build/${BINARY_NAME} cmd/lets-go-chat/main.go

run:
	go build -o build/${BINARY_NAME} cmd/lets-go-chat/main.go
	bin/${BINARY_NAME}

test:
	go test ./...

clean:
	go clean
	rm build/${BINARY_NAME}