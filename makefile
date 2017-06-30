BINARY=healthy
VERSION_TAG=`git describe --tags 2>/dev/null | cut -f 1 -d '-' 2>/dev/null`
COMMIT_HASH=`git rev-parse --short=8 HEAD 2>/dev/null`
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-s -w \
	-X main.CommitHash=${COMMIT_HASH} \
	-X main.BuildTime=${BUILD_TIME} \
	-X main.VersionTag=${VERSION_TAG}"

all: build

release: clean install linux darwin windows freebsd

clean:
	go clean
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf ./release || true
	
install: clean
	mkdir ./release
	glide install
	

build:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	go build -o ${BINARY} ${LDFLAGS}

linux:
	GOOS=linux GOARCH=386 go build ${LDFLAGS} -o ./release/${BINARY}_linux_386
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ./release/${BINARY}_linux_amd64

darwin:
	GOOS=darwin GOARCH=386 go build ${LDFLAGS} -o ./release/${BINARY}_darwin_386
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ./release/${BINARY}_darwin_amd64
windows:
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o ./release/${BINARY}_windows_386.exe
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ./release/${BINARY}_windows_amd64.exe

freebsd:
	GOOS=freebsd GOARCH=386 go build ${LDFLAGS} -o ./release/${BINARY}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build ${LDFLAGS} -o ./release/${BINARY}_freebsd_amd64



.PHONY: build all clean release install linux darwin windos freebsd






