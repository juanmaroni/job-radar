FROM golang:1.27-bookworm AS build

WORKDIR /service

COPY events/ ./events/
COPY services/ingestion-service/ ./services/ingestion-service/

WORKDIR /service/services/ingestion-service
RUN go mod download
RUN go build -o ingestion ./cmd


FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /service

COPY --from=build /service/services/ingestion-service/ingestion /service/bin/ingestion

RUN chmod +x /service/bin/ingestion

ENTRYPOINT ["/service/bin/ingestion"]
