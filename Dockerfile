FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod vendor && \
    cd vendor/github.com/usenwep/nwep-go && bash setup.sh
RUN go build -mod=vendor -o server .

FROM debian:trixie-slim
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 3000
CMD ["./server"]
