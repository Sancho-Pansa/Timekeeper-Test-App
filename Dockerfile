FROM golang:1.27-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -v -o /usr/local/bin/app ./...

FROM alpine:3

WORKDIR /app

COPY --from=builder /usr/local/bin/app .

EXPOSE "8080"

CMD [ "sh", "-c", "./app" ]