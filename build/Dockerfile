FROM golang:latest

WORKDIR /opt/app

COPY . .

RUN go build cmd/app/main.go

EXPOSE 8080

CMD ["./main"]