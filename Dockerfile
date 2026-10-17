FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
  go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/timekeeper \
    ./cmd/server

FROM scratch

WORKDIR /app

COPY --from=builder /out/timekeeper /app/timekeeper

USER "65532:65532"

EXPOSE 8080

ENTRYPOINT [ "/app/timekeeper" ]