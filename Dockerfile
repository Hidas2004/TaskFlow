# --- GIAI ĐOẠN 1: Builder ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy file thư viện trước để tận dụng cache
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source code
COPY . .

# Build ra file thực thi tên là 'server'
RUN go build -o server ./cmd/api/main.go

# --- GIAI ĐOẠN 2: Runner ---
FROM alpine:latest

# Cài chứng chỉ bảo mật
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Lấy file 'server' từ giai đoạn 1 bỏ sang đây
COPY --from=builder /app/server .

# Copy file môi trường .env
COPY .env .

# Tạo thư mục uploads (nếu cần)
RUN mkdir uploads

# Mở cổng 8080
EXPOSE 8080

# Chạy server
CMD ["./server"]