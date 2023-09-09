# Build the manager binary
FROM golang:1.20 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY migrations/migrations.go migrations/migrations.go
COPY config/ config/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o migrator migrations/migrations.go

FROM gcr.io/distroless/static:nonroot
ENV SCYLLA_DB="localhost:9042"
WORKDIR /
COPY --from=builder /workspace/migrator .
COPY /migrations/cql ./migrations/cql

ENTRYPOINT ["/migrator"]
