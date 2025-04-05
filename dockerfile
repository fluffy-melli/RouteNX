FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o routenx .

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache certbot openssl

COPY --from=builder /app/routenx .
COPY --from=builder /app/RouteNX.json ./RouteNX.json
COPY --from=builder /app/dist ./dist

EXPOSE 80
EXPOSE 443
EXPOSE 3000

CMD ["./routenx"]