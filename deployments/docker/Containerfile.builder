# syntax=docker/dockerfile:1
#
# GoREveal reproducible build image.
#
# Produces a single static binary in /out. `docker build --target export
# --output type=local,dest=./dist` extracts it without an intermediate image.
#
# Base image is digest-pinned. Refresh it deliberately.
FROM docker.io/library/golang:1.26-trixie@sha256:4ee9ffa999b4583ce281939cdff828763083610292f252279a0cee77473bd9a7 AS builder

LABEL org.opencontainers.image.title="GoREveal Builder" \
      org.opencontainers.image.description="Reproducible build image for GoREveal" \
      org.opencontainers.image.source="https://github.com/ioplane/goreveal" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.vendor="ioplane"

ENV GOFLAGS=-buildvcs=false \
    CGO_ENABLED=0

WORKDIR /src

# Dependency layer first so source edits do not invalidate the module cache.
COPY go.mod go.sum go.work go.work.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

COPY . .

ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown
ARG TARGETOS
ARG TARGETARCH

# -trimpath and an empty build ID keep the output byte-identical across hosts for
# a given (VERSION, GIT_COMMIT, BUILD_DATE) triple.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS="${TARGETOS:-linux}" GOARCH="${TARGETARCH:-amd64}" \
    go build -trimpath \
      -ldflags="-s -w -buildid= \
        -X github.com/ioplane/goreveal/internal/version.Version=${VERSION} \
        -X github.com/ioplane/goreveal/internal/version.GitCommit=${GIT_COMMIT} \
        -X github.com/ioplane/goreveal/internal/version.BuildDate=${BUILD_DATE}" \
      -o /out/goreveal ./cmd/goreveal \
    && /out/goreveal version

# Scratch stage for `--output type=local`: contains only the built artifact.
FROM scratch AS export
COPY --from=builder /out/goreveal /goreveal
