FROM golang:1.27-bookworm AS build

WORKDIR /service
COPY go.mod go.sum /service
RUN go mod download
COPY . /service
RUN go build -o server cmd/main.go


FROM debian:bookworm-slim

WORKDIR /service
COPY --from=build /service/server /service/bin/server

RUN chmod +x /service/bin/server

EXPOSE 8080

ENTRYPOINT ["/service/bin/server"]
