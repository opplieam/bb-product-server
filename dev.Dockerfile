FROM golang:1.22-alpine3.19 AS builder
ENV CGO_ENABLED 0
ARG BUILD_REF

# Copy the source code into the container.
RUN go env -w GOCACHE=/go-cache
COPY . /service
# Build the service binary.
WORKDIR /service/cmd/server
RUN --mount=type=cache,target=/go-cache go build -ldflags "-X main.build=${BUILD_REF}" -o server

# Run the Go Binary in Alpine.
FROM alpine:3.19
ARG BUILD_DATE
ARG BUILD_REF
WORKDIR /service
COPY --from=builder /service/cmd/server/server .
CMD ["./server"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="buy-better-product-server" \
      org.opencontainers.image.revision="${BUILD_REF}"
