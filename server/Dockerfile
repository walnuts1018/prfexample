# syntax=docker/dockerfile:1.22
FROM golang:1.26.6-trixie AS builder

ENV ROOT=/build
ARG BUILD_TAGS=""
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    sh -c '\
    if [ -n "$BUILD_TAGS" ]; then \
    GOOS=linux go build -trimpath -o app -tags "$BUILD_TAGS" $ROOT/main.go; \
    else \
    GOOS=linux go build -trimpath -o app $ROOT/main.go; \
    fi' && \
    chmod +x ./app

FROM debian:13.8-slim@sha256:8aa5a7bf5a1955bddbe7d74fcdd0b91ab30e74362a8b73af672b0481b35b44a0
WORKDIR /app

RUN rm -f /etc/apt/apt.conf.d/docker-clean; echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
RUN --mount=type=cache,target=/var/lib/apt,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get -y update && apt-get upgrade -y && apt-get install -y ca-certificates

COPY --from=builder /build/app ./
EXPOSE 8080

CMD ["./app"]
