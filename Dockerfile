FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /app/server cmd/main.go


FROM alpine:latest

WORKDIR /app

LABEL org.opencontainers.image.title="Proxy M3U8" \
  org.opencontainers.image.description="Proxy server for M3U8" \
  org.opencontainers.image.maintainer="dovakiin0@kitsunee.online" \
  org.opencontainers.image.source="https://github.com/dovakiin0/proxy-m3u8.git" \
  org.opencontainers.image.vendor="kitsunee.online"

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/server /app/server

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 4040

ENTRYPOINT ["/app/server"]
