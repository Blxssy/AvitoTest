FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go mod tidy

COPY . ./

RUN go build -o main ./cmd/app/main.go

CMD ["./main"]