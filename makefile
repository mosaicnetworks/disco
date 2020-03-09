



run:
	go run server/cmd/main.go

build: 
	go build -o build/disco main.go

.PHONY: run build