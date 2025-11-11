FROM golang:alpine

WORKDIR /message_service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/run

EXPOSE 12911

CMD ["./main"]