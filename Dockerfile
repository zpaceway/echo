# -------------------------- *** Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ohce .

# -------------------------- *** Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -g 1001 -S ohce && adduser -S ohce -u 1001 -G ohce
WORKDIR /app
COPY --from=builder /app/ohce .
RUN chown ohce:ohce ohce
USER ohce
EXPOSE 6868

CMD ["./ohce"]
