FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /timekeeper ./cmd/server

FROM alpine:3

WORKDIR /app

COPY --from=builder /timekeeper /app/timekeeper

EXPOSE "8080"

ENV APP_USERNAME="user"
ENV APP_PASSWORD="password"
ENV APP_VERSION="0.0.10"

CMD [ "/app/timekeeper" ]