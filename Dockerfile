FROM golang:1.24.2 AS builder
WORKDIR /build
COPY ./src ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./main

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/main ./main
EXPOSE 8888
ENTRYPOINT ["./main"]