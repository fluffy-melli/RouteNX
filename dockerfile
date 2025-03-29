FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o routenx .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/routenx .
COPY --from=builder /app/RouteNX.json ./

EXPOSE 8080

CMD ["./routenx"]