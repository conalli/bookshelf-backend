FROM golang:1.17-alpine

WORKDIR /go/src/github.com/bookshelf-backend

COPY . .

RUN go mod download

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]