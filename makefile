



run:
	go run server/cmd/main.go

build: 
	go build -o build/disco server/cmd/main.go

.PHONY: run build