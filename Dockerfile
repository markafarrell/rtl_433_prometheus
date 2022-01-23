ARG GO_VERSION=1.14
ARG DEBIAN_VERSION=11.2-slim

FROM golang:${GO_VERSION} as gobuilder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.* ./
RUN go mod download -x

# Copy local code to the container image.
COPY . ./

RUN go env GOOS GOARCH GOARM

# Build the binary.
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -mod=readonly -a -v rtl_433_prometheus.go

FROM debian:${DEBIAN_VERSION}

RUN apt update && \
    apt upgrade && \
    apt install -y --no-install-recommends \
        rtl-433 \
        udev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /
COPY --from=gobuilder /app/rtl_433_prometheus /

ENV FREQUENCY=433M
ENV RTL_SDR_DEV_PATH=/dev/rtl_sdr
ENV RTL_433_FLAGS=
ENV MATCHERS=

COPY entrypoint.sh /entrypoint.sh

EXPOSE 9550
ENTRYPOINT ["/entrypoint.sh"]
