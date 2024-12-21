# Build stage
FROM golang:1.23-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.21
WORKDIR /app

COPY --from=builder /app/main .

COPY start.sh wait-for.sh ./

EXPOSE 9090
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
