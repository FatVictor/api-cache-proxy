FROM golang AS builder
RUN mkdir /build
COPY . /build
RUN cd /build && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM alpine:latest
RUN mkdir /app
COPY --from=builder /build/main /app/main
RUN chmod 711 /app/main
WORKDIR /app
CMD ["/app/main"]

