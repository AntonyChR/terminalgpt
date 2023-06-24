build:
	go build -o bin/gpt
dist:
	go build -o bin/gpt && sudo cp bin/gpt /bin/gpt
