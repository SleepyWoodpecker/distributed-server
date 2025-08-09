
FILES = cmd/main.go
BUILD = out

build: 
	go build -o $(BUILD)/server $(FILES)

run: build
	./$(BUILD)/server

test:
	@# this matches all subdirectories in the current directory
	go test ./... -v

fmt:
	go fmt ./...