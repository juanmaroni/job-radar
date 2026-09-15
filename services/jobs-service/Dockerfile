FROM golang:1.27-bookworm AS build

WORKDIR /service

COPY events/ ./events/
COPY services/jobs-service/ ./services/jobs-service/

WORKDIR /service/services/jobs-service
RUN go mod download
RUN go build -o jobs ./cmd


FROM debian:bookworm-slim

WORKDIR /service

COPY --from=build /service/services/jobs-service/jobs /service/bin/jobs

RUN chmod +x /service/bin/jobs

ENTRYPOINT ["/service/bin/jobs"]
