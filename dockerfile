FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o routenx .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/dist ./dist
COPY --from=builder /app/routenx ./routenx
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/RouteNX.json ./RouteNX.json

EXPOSE 80
EXPOSE 443
EXPOSE 3000

CMD ["./routenx"]