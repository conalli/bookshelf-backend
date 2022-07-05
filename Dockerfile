FROM golang:1.18-alpine as build

WORKDIR /go/src/github.com/bookshelf-backend

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:3.15

WORKDIR /app

COPY --from=build /go/src/github.com/bookshelf-backend/main .

EXPOSE 8080

CMD ["./main"]