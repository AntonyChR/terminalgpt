build:
	go build -o bin/gpt

dist:
	go build -o bin/gpt && sudo cp bin/gpt /bin/gpt

run:
	go run main.go

test:
	go clean -testcache
	go test -v ./...