package = github.com/WianVos/healthy

.PHONY: install release 

install:
	go get -t -v ./...

release:
	mkdir -p release
	rm -rf release/*
	GOOS=linux GOARCH=amd64 go build -o release/healthy $(package)
	