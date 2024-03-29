ARG VERSION=1.20

FROM golang:${VERSION}-alpine AS builder

# See https://github.com/hadolint/hadolint/wiki/DL4006
SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

# Hadolint expects you to specify exact versions. Since this is a build container ignore this check
# hadolint ignore=DL3018
RUN apk add --update --no-cache \
    # Git is required for fetching the dependencies.
    git \
    # For CGO_ENABLED=1 we need these gcc and alpine SDK packages
    gcc alpine-sdk ca-certificates

ARG UID=10001
ARG USER=appuser

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR $GOPATH/src/package

COPY . .

RUN go get \
 && go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init --output swagger -g cmd/qq/main.go \
 && go test ./... -timeout=60s -parallel=32

WORKDIR $GOPATH/src/package/cmd/qq

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/server

############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /etc/passwd /etc/group /etc/
COPY --from=builder /go/bin/server /go/bin/server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/pki/tls/certs/

# Use an unprivileged user.
USER appuser:appuser

# Run the binary
ENTRYPOINT ["/go/bin/server"]
