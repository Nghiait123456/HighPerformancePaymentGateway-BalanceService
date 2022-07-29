# Build stage
# Build stage
FROM golang:1.18-alpine3.15 AS builder
WORKDIR /app
ADD go.mod go.sum /app/
RUN go mod download

COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .

CMD [ "/app/main" ]
