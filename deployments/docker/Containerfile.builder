FROM docker.io/golang:1.26-trixie AS builder

LABEL org.opencontainers.image.title="GoREveal Builder"
LABEL org.opencontainers.image.description="Reproducible build image for GoREveal"
LABEL org.opencontainers.image.source="https://github.com/dantte-lp/goreveal"

ENV GOFLAGS=-buildvcs=false \
    CGO_ENABLED=0

WORKDIR /src
COPY go.mod go.work ./
COPY internal ./internal
COPY . .

ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

RUN if [ -d ./cmd/goreveal ]; then \
      go build -trimpath \
      -ldflags="-s -w -X github.com/dantte-lp/goreveal/internal/version.Version=${VERSION} -X github.com/dantte-lp/goreveal/internal/version.GitCommit=${GIT_COMMIT} -X github.com/dantte-lp/goreveal/internal/version.BuildDate=${BUILD_DATE}" \
      -o /out/goreveal ./cmd/goreveal ; \
    else \
      mkdir -p /out ; \
      echo "goreveal builder scaffold ready; cmd/goreveal not implemented yet" > /out/README.txt ; \
    fi
