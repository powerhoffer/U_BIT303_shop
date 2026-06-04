GF ?= /usr/local/go/gopath_dir/bin/gf

.PHONY: run test tidy dao service

run:
	go run main.go

test:
	go test ./...

tidy:
	go mod tidy

dao:
	$(GF) gen dao

service:
	$(GF) gen service
