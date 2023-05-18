# lets-go-chat

## Getting Started

These instructions will get you a copy of the project up and running on your local machine.

### Prerequisites

- Go: [asdf](https://asdf-vm.com/), [asdf go plugin](https://github.com/kennyp/asdf-golang)
- GNU Make: [brew make](https://formulae.brew.sh/formula/make#default), [documentation](https://www.gnu.org/software/make/)

### Running App

```sh
$ go run cmd/lets-go-chat/main.go
```

Makefile:

```sh
$ make run
```

### Running Tests

```sh
$ go test ./...
```

Makefile:

```sh
$ make test
```

### Documentation

- Install GoDoc: `$ go install golang.org/x/tools/cmd/godoc@latest`
- Run: `go godoc -http=:6060`
- Open: `localhost:6060`