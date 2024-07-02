export CGO_ENABLED=0

build:
	go build -o bin/ cmd/fsdupscan/fsdupscan.go

run: build
	./bin/fsdupscan 

test:
	go test -v ./... 
