FROM golang:latest

WORKDIR ${GOPATH}/avito-shop/
COPY . ${GOPATH}/avito-shop/

RUN go mod download

RUN go mod tidy

COPY . ./

RUN go build -o main ./cmd/app/main.go \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["./main"]