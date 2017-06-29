FROM golang:1.8.3-alpine3.6

WORKDIR /go/src/app
COPY . . 

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"] # ["app"]