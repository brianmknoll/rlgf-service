FROM golang:1.23.7 AS builder

WORKDIR /app

COPY ./src .
RUN go mod download && go mod tidy
RUN go build -o main .

FROM ubuntu

WORKDIR /app

COPY --chmod=0755 --from=builder /app/main .

EXPOSE 8888

CMD ["./main"]
