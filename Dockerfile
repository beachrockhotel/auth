FROM golang:1.20.3-alpine AS builder

WORKDIR /app
COPY --from=builder /app/bin/crud_server .

RUN go mod download
RUN go build -o ./bin/crud_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/beachrockhotel/auth/source/bin/crud_server .

CMD ["./crud_server"]