FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o main .

FROM alpine:3.19
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]