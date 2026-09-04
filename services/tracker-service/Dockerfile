FROM golang:1.27-bookworm AS build

WORKDIR /service

COPY proto/ ./proto/
COPY services/tracker-service/ ./services/tracker-service/

WORKDIR /service/services/tracker-service
RUN go mod download
RUN go build -o tracker cmd/main.go


FROM debian:bookworm-slim

WORKDIR /service

COPY --from=build /service/services/tracker-service/tracker /service/bin/tracker

RUN chmod +x /service/bin/tracker

EXPOSE 8080

ENTRYPOINT ["/service/bin/tracker"]
