# Build
FROM golang:1.19-alpine as builder
WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -o /transactions-proxy-backend ./cmd/transactions-proxy-backend/main.go

# Copy
FROM alpine:3
COPY --from=builder transactions-proxy-backend /bin/transactions-proxy-backend
ENTRYPOINT ["/bin/transactions-proxy-backend"]

