# Build stage
# Build stage
FROM golang:1.18-alpine3.15 AS builder
WORKDIR /App
ADD go.mod go.sum /App/
RUN go mod download

COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.15
WORKDIR /App
COPY --from=builder /App/main .

CMD [ "/App/main" ]
