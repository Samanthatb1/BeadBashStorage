# Build Stage
FROM golang:1.18-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY start.sh .
COPY db/migration ./db/migration

ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /app/start.sh

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT ["sh", "/app/start.sh" ]
