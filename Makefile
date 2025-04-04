BINARY_NAME = gpt
OUTPUT_DIR = bin
INSTALL_DIR = /usr/local/bin

build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME)

install: build
	sudo cp $(OUTPUT_DIR)/$(BINARY_NAME) $(INSTALL_DIR)

run:
	go run main.go

test:
	go clean -testcache
	go test -v ./...

clean:
	rm -rf $(OUTPUT_DIR)

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64
