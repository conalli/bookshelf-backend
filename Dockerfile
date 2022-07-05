FROM golang:1.18-alpine as build

WORKDIR /go/src/github.com/conalli/bookshelf-backend

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd/bookshelf-server/main.go

FROM alpine:3.16

WORKDIR /app

COPY --from=build /go/src/github.com/conalli/bookshelf-backend/main .

EXPOSE 8080

CMD ["./main"]