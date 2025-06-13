FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

WORKDIR /app/cmd/server
RUN go build -o /app/go-web-analyzer

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/go-web-analyzer .
COPY web /app/web

EXPOSE 8080
EXPOSE 6060

CMD ["./go-web-analyzer"]
