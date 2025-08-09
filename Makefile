
FILES = cmd/main.go

build: 
	go build -o server $(FILES)

run: build
	./server

test:
	@# this matches all subdirectories in the current directory
	go test ./... -v