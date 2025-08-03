# Stage 1: Build
FROM golang:1.21 as builder

WORKDIR /app
COPY . .
RUN go build -o checker ./cmd/checker

# Stage 2: Run
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/checker /app/checker
COPY compliance-rules.yaml /app/compliance-rules.yaml

# Optional: add ca-certificates if using TLS NetBox
RUN apk add --no-cache ca-certificates

ENV NETBOX_URL=http://netbox.argo.local
ENV NETBOX_TOKEN=REPLACE_ME

ENTRYPOINT ["/app/checker"]
