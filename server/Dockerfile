# syntax=docker/dockerfile:1.22
FROM golang:1.26.2-trixie AS builder

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

FROM debian:13.3-slim
WORKDIR /app

RUN rm -f /etc/apt/apt.conf.d/docker-clean; echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
RUN --mount=type=cache,target=/var/lib/apt,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get -y update && apt-get upgrade -y && apt-get install -y ca-certificates

COPY --from=builder /build/app ./
EXPOSE 8080

CMD ["./app"]
