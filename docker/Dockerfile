# syntax=docker/dockerfile:1

FROM golang:1.20 as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . ./
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux make build

FROM scratch

COPY --from=builder /app/bin/cb-tracker /cb-tracker
COPY --from=builder /app/configuration/local.yaml /config.yml
ENTRYPOINT ["/cb-tracker", "--config", "/config.yml"]