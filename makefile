BINARY=healthy
DIST_DIRS := find * -type d -exec
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
	mkdir -p "${GOPATH}/bin"

	
bootstrap: clean bootstrap-gox bootstrap-glide
	mkdir ./release
	glide install 

build: bootstrap
	go build -o ${BINARY} ${LDFLAGS} 

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./${BINARY} ${DESTDIR}/usr/local/bin/${BINARY}

build-all: bootstrap 
	gox -verbose  \
	-os="darwin linux freebsd windows" \
	-arch="amd64" \
	${LDFLAGS} \
	-output="release/{{.OS}}-{{.Arch}}/{{.Dir}}" 

bootstrap-gox:
	go get -u github.com/mitchellh/gox
	cd ${GOPATH}/src/github.com/mitchellh/gox && go install

bootstrap-glide:
	go get -u github.com/Masterminds/glide
	cd ${GOPATH}/src/github.com/Masterminds/glide && go install


dist: build-all
	cd release && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf ${BINARY}-${VERSION_TAG}-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r ${BINARY}-${VERSION_TAG}-{}.zip {} \; && \
	cd ..


.PHONY: build all clean release install linux darwin windos freebsd






