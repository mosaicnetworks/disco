



run:
	go run server/cmd/main.go --cert-file=test_data/cert.pem \
							  --key-file=test_data/key.pem

build: 
	go build -o build/disco server/cmd/main.go

.PHONY: run build