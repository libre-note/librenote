ARG GO_VERSION=1.17

FROM golang:${GO_VERSION}-alpine AS builder

ARG BUILD_VERSION=0.0.1
ARG BUILD_TIME=00000000-000000

# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates
RUN apk add --no-cache ca-certificates git tini tzdata

# Working directory outside $GOPATH
WORKDIR /build

# Copy go module files and download dependencies
COPY go.* ./
RUN go mod download

# Copy source files
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X librenote/app.Version=$BUILD_VERSION -X librenote/app.BuildTime=$BUILD_TIME" -o librenote .


# Minimal image for running the application
FROM scratch as final

LABEL org.opencontainers.image.source="https://github.com/libre-note/core" \
      org.opencontainers.image.url="https://github.com/libre-note/core" \
      org.opencontainers.image.title="A note taking applications"

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import time zone info for golang time package
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Import the init binary
COPY --from=builder /sbin/tini /sbin/tini
# Import curl from curl repository image
COPY --from="curlimages/curl:7.81.0" /usr/bin/curl /usr/bin/curl
# Import the compiled executable from the first stage.
COPY --from=builder /build/librenote /app/librenote

# Open port
EXPOSE 8000

# Perform any further action as an unprivileged user.
USER nobody:nobody
WORKDIR /app

HEALTHCHECK --interval=5m --timeout=2s --retries=3 --start-period=2s CMD ["curl", "--fail", "http://localhost:8000/h34l7h"]

# Run the compiled binary.
ENTRYPOINT ["/sbin/tini", "--"]
ENTRYPOINT ["/app/librenote","--config","/app/config.yml", "serve"]
