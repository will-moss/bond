FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /app
COPY go.mod  .
COPY go.sum  .
COPY main.go .
COPY default.env .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bond main.go

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /

COPY --from=builder /app/bond .
COPY --from=builder /app/default.env .

ENTRYPOINT ["./bond"]
