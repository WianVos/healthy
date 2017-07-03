FROM alpine:3.6

ADD https://github.com/WianVos/healthy/releases/download/untagged-e3be289d9a955edbf5f3/healthy-untagged-linux-amd64.zip /tmp
RUN unzip /tmp/healthy-untagged-linux-amd64.zip 

