# Multi-stage, multi-arch-aware Dockerfile for hello-docker.
#
# Cross-compiles via Go's native TARGETARCH support (no QEMU emulation needed
# for the build stage), then ships a scratch-based final image.

# syntax=docker/dockerfile:1.7
ARG GO_VERSION=1.27

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION} AS build
ARG TARGETARCH
ARG VERSION=unknown
WORKDIR /src
# Copy the whole module (go.mod + app/) so `go build` can resolve the module.
COPY . .
# CGO disabled + GOARCH set per TARGETARCH -> static binary per platform.
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOARCH=${TARGETARCH} GOOS=linux \
    go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /out/hello ./app

FROM scratch
# No shell in scratch; tini not needed for a single Go process.
COPY --from=build /out/hello /hello
EXPOSE 8080
ENTRYPOINT ["/hello"]
