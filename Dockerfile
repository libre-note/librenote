ARG GO_VERSION=1.17
ARG CURL_VERSION=7.81.0

FROM golang:${GO_VERSION}-alpine AS builder

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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o librenote-core .


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
COPY --from=curlimages/curl:${CURL_VERSION} /bin/curl /bin/curl
# Import the compiled executable from the first stage.
COPY --from=builder /build/librenote-core /app/librenote-core

# Open port
EXPOSE 8000

# Perform any further action as an unprivileged user.
USER nobody:nobody
WORKDIR /app

HEALTHCHECK --interval=1m --timeout=2s --retries=3 --start-period=2s CMD ["curl", "--fail", "http://localhost:8000/health"]

# Run the compiled binary.
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/librenote-core","--config","/app/config.yml"]
