.PHONY: run test tidy dao service

run:
	go run main.go

test:
	go test ./...

tidy:
	go mod tidy

dao:
	gf gen dao

service:
	gf gen service
