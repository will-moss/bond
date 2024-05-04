FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod  .
COPY go.sum  .
COPY main.go .
COPY default.env .

RUN CGO_ENABLED=0 GOOS=linux go build -o bond main.go

FROM gcr.io/distroless/base

WORKDIR /

COPY --from=builder /app/bond .
COPY --from=builder /app/default.env .

ENTRYPOINT ["./bond"]
