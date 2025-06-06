# web_ordering_system/backend/Dockerfile
# 使用 Go 官方映像作為基礎映像，適合編譯 Go 應用程式
FROM golang:1.23-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 將 go.mod 和 go.sum 複製到容器中，並下載依賴
# 這樣做可以利用 Docker 的層快取，如果 go.mod/go.sum 沒變，就不會重複下載依賴
COPY go/go.mod .
COPY go/go.sum .
# maybe can utilize (optimize) , because copy go/ of full file, contain go/go.mod and go.sum of files
RUN go mod download

# 將所有 Go 原始碼複製到容器中
COPY go/ .

# 編譯 Go 應用程式
# CGO_ENABLED=0 和 -ldflags -s -w 用於創建較小的、不依賴 C 庫的靜態連結可執行檔
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# --- 第二階段：輕量級運行時映像 ---
# 使用一個更小的、僅包含運行時必要的映像
FROM alpine:latest

# 設定工作目錄
WORKDIR /app

# 從第一階段的 builder 映像中複製編譯好的可執行檔
COPY --from=builder /app/main .

# 將資料庫檔案的目錄掛載點創建好，以便卷宗掛載
# 確保這個目錄在容器內部存在，並與 docker-compose.yaml 中的 volumes 映射路徑一致
# 這裡假設 SQLite 檔案將存放在 /app/data
RUN mkdir -p /app/data

# 暴露 8080 埠，這是 Go 應用程式監聽的埠口
EXPOSE 8080

# 定義容器啟動時執行的命令
# 運行編譯好的 Go 應用程式
CMD ["/app/main"]
