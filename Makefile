
run:
	go run cmd/server/main.go

test:
	go test -v ./...

build:
	go build -o build/ cmd/server/main.go
