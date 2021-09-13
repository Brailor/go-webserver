FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN GOOS=linux go build server.go

CMD [ "./server" ]
