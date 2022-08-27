FROM golang:1.16-alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o auth-service .

EXPOSE 8081

ENTRYPOINT ["./auth-service"]
