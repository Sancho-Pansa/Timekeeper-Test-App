FROM golang:1.27-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o ./cmd ./cmd/server/main.go
RUN go build -o ./timezones ./internal/timezones.go

FROM alpine:3

WORKDIR /app

COPY --from=builder /usr/src/app/cmd .
COPY --from=builder /usr/src/app/timezones .

EXPOSE "8080"

CMD [ "sh", "-c", "./main" ]