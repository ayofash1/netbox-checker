# Stage 1: Build
FROM golang:1.21 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o checker ./cmd

# Stage 2: Run
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/checker /app/checker
COPY compliance-rules.yaml /app/compliance-rules.yaml

RUN apk add --no-cache ca-certificates

ENV NETBOX_URL=http://netbox.argo.local
ENV NETBOX_TOKEN=REPLACE_ME

ENTRYPOINT ["/app/checker"]
