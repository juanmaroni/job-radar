FROM golang:1.27-bookworm AS build

WORKDIR /service

COPY proto/ ./proto/
COPY events/ ./events/
COPY services/matching-service/ ./services/matching-service/

WORKDIR /service/services/matching-service
RUN go mod download
RUN go build -o matching ./cmd


FROM debian:bookworm-slim

WORKDIR /service

COPY --from=build /service/services/matching-service/matching /service/bin/matching

RUN chmod +x /service/bin/matching

EXPOSE 50051

ENTRYPOINT ["/service/bin/matching"]
