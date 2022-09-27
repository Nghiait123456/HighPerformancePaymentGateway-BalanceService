# Build stage
FROM golang:1.18-alpine3.15 AS builder
WORKDIR /app
COPY . .

RUN apk add build-base
RUN go install github.com/google/wire/cmd/wire@latest
RUN cd ./balance/  && wire
RUN go build -o main main.go

# Run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD [ "./main" ]
