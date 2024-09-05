FROM golang:latest

WORKDIR /go/src/client

COPY . .

RUN go mod download

RUN go build -o main .

EXPOSE 8180

CMD ["./main"]

