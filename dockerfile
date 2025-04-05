FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o routenx .

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache certbot openssl

RUN mkdir -p /etc/letsencrypt/live && \
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
	-keyout /etc/letsencrypt/live/privkey.pem \
	-out /etc/letsencrypt/live/fullchain.pem \
	-subj "/C=US/ST=State/L=City/O=Organization/OU=Unit/CN=localhost"

COPY --from=builder /app/routenx .
COPY --from=builder /app/RouteNX.json ./
COPY --from=builder /app/dist ./dist

EXPOSE 80
EXPOSE 443
EXPOSE 3000

CMD ["./routenx"]