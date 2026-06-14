FROM golang:1.26.3-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /tstun ./cmd/tstun

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -H -s /sbin/nologin tstun

COPY --from=builder /tstun /usr/bin/tstun

RUN mkdir -p /data/tsnet_state /etc/tstun && \
    chown -R tstun:tstun /data/tsnet_state

EXPOSE 80 443 2019

VOLUME ["/data/tsnet_state"]

USER tstun

ENTRYPOINT ["tstun"]
CMD ["run", "--config", "/etc/tstun/Caddyfile"]
