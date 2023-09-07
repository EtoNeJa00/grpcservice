# Build the manager binary
FROM golang:1.19 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/ cmd/
COPY internal/ internal/
COPY config/ config/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o grpcsevice cmd/serve/serve.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
ENV MEMCACHE_GRPC_PORT: ":8081"
ENV INTERNAL_STORAGE_GRPC_PORT: ":8080"
ENV SCYLLA_GRPC_PORT: ":8082"
ENV MEMCACHED_ADDR="localhost:11211"
ENV SCYLLA_DB="localhost:9042"
WORKDIR /
COPY --from=builder /workspace/grpcsevice .

ENTRYPOINT ["/grpcsevice"]
