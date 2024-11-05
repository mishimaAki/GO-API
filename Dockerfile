FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

RUN apk --nocache add tzdata
ENV TZ=Asia/Tokyo

EXPOSE 8080

CMD ["./main"]