# -------------------------- *** Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o echo .

# -------------------------- *** Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -g 1001 -S echo && adduser -S echo -u 1001 -G echo
WORKDIR /app
COPY --from=builder /app/echo .
RUN chown echo:echo echo
USER echo
EXPOSE 6868

CMD ["./echo"]
