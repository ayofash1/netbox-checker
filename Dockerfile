# Stage 1: Build
FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o checker ./cmd/checker

# Stage 2: Runtime
FROM ubuntu:22.04

WORKDIR /app

# Install CA certificates (for HTTPS)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/checker /app/checker
COPY compliance-rules.yaml /app/compliance-rules.yaml

ENV NETBOX_URL=http://netbox.argo.local
ENV NETBOX_TOKEN=REPLACE_ME

ENTRYPOINT ["/app/checker"]
