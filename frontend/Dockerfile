# 使用最新的 Ubuntu LTS 版本作為基礎映像
FROM ubuntu:22.04

# 設定維護者資訊 (可選)
LABEL maintainer="your_email@example.com"

# 更新 apt 軟體包列表並安裝 nginx
# --no-install-recommends 減少安裝不必要的推薦軟體包
# rm -rf /var/lib/apt/lists/* 清理 apt 快取以減小映像大小
RUN apt-get update && \
    apt-get install -y --no-install-recommends nginx && \
    rm -rf /var/lib/apt/lists/*

# 移除 Nginx 的預設網站配置，我們將提供自己的配置
RUN rm /etc/nginx/sites-enabled/default

# 將自定義的 nginx 配置檔案複製到容器中
# 這個檔案將定義如何服務靜態檔案和反向代理 API 請求
COPY nginx.conf /etc/nginx/sites-available/default.conf

# 啟用剛複製的 Nginx 配置 (創建軟連結)
RUN ln -s /etc/nginx/sites-available/default.conf /etc/nginx/sites-enabled/default

# 在容器中創建一個目錄來存放你的靜態檔案
# Nginx 將會從這個目錄服務您的 HTML, CSS, JS, Images
RUN mkdir -p /usr/share/nginx/html/static

# 將所有前端靜態檔案複製到 Nginx 服務的目錄中
# 注意：這些路徑是相對於 Dockerfile 所在目錄 (即 frontend/)
# 您需要確保這些檔案在構建時都位於 frontend/ 目錄內或其子目錄中
COPY index.html /usr/share/nginx/html/
COPY css/ /usr/share/nginx/html/static/css/
COPY js/ /usr/share/nginx/html/static/js/
COPY images/ /usr/share/nginx/html/static/images/
# 如果您還有其他靜態檔案或目錄，可以按需添加更多 COPY 指令

# 暴露 80 埠，這是 Nginx 預設監聽 HTTP 請求的埠口
EXPOSE 80

# 定義容器啟動時執行的命令
# "daemon off;" 確保 Nginx 在前台運行，這樣 Docker 才能監控它
CMD ["nginx", "-g", "daemon off;"]
