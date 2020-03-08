



run:
	go run main.go

build: 
	go build -o build/disco main.go

.PHONY: run build