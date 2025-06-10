# web_ordering_system/backend/Dockerfile
# 使用 Go 官方映像的 Alpine 版本作為基礎映像
# 這會讓 Go 編譯器連結到 musl libc，與 Alpine 運行時兼容
FROM golang:1.23 AS builder

WORKDIR /app

# 將 go.mod 和 go.sum 複製到容器中，並下載依賴
COPY go/go.mod .
COPY go/go.sum .
RUN go mod download

# 將所有 Go 原始碼複製到容器中
COPY go/ .

# 確保 CGO_ENABLED 為 1，讓 go-sqlite3 能正常工作
ENV CGO_ENABLED=1

# 【重要】在所有原始碼複製後，且編譯前，再次執行 go mod tidy
# 這確保 go.mod 和 go.sum 始終與當前的原始碼狀態同步
RUN go mod tidy

# 編譯 Go 應用程式。因為 builder 已經是 Alpine 基礎，所以編譯出來的執行檔
# 將會與 Alpine 運行時兼容。
# 不需再手動加 -a -installsuffix cgo，Go 會根據基礎映像自動處理
RUN COMPOSE_BAKE=true GOOS=linux go build -o main .

# --- 第二階段：輕量級運行時映像 ---
FROM ubuntu:latest

WORKDIR /app

RUN mkdir -p /app/data
# COPY ./sql/order_db.db ./data # just for test docker file 
# 從第一階段的 builder 映像中複製編譯好的可執行檔
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]