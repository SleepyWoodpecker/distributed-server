build: 
	go build -o server

run: build
	./server

test:
	@# this matches all subdirectories in the current directory
	go test ./... -v