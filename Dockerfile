# === Development ===
FROM golang:alpine AS dev
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# 代码通过 volume 挂载, 修改后运行: docker compose restart app
CMD ["go", "run", "./cmd/server"]

# === Build ===
FROM golang:alpine AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

# === Production ===
FROM alpine:latest AS prod
WORKDIR /app
RUN apk add --no-cache tzdata
COPY --from=builder /server .
COPY configs/ ./configs/
EXPOSE 8080
CMD ["./server"]
