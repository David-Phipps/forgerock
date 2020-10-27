FROM golang:1.14

WORKDIR /go/src/app

COPY ./app/main.go  .

EXPOSE 9000

RUN go get -d -v ./...

RUN go install -v ./...

CMD ["app"]
