# 优化后的 Dockerfile
# 使用多阶段构建，最小化最终镜像体积

# Go 后端构建阶段
FROM golang:1.21-alpine AS golang
WORKDIR /backend

# 启用 Go Modules 缓存，静态编译
ENV GO111MODULE=on CGO_ENABLED=0

# 先复制 go.mod 和 go.sum，利用 Docker 层缓存
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

# 再复制源码进行构建
COPY ./backend .
RUN go build -ldflags="-s -w" -o rustdesk-api-server-pro .

# Node 前端构建阶段
FROM node:20-alpine AS node
WORKDIR /frontend

COPY ./soybean-admin ./
RUN npm install -g pnpm && pnpm i --frozen-lockfile

RUN pnpm build

# 最终运行阶段
FROM alpine:3.20.3
ENV ADMIN_USER=
ENV ADMIN_PASS=
WORKDIR /app

# 添加非 root 用户
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# 复制构建产物
COPY --from=golang /backend/rustdesk-api-server-pro .
COPY --from=golang /backend/server.yaml .
COPY --from=node /frontend/dist ./dist
COPY --chmod=755 ./docker/start.sh .

# 安装时区数据 + 进程降权工具
RUN apk add --no-cache tzdata su-exec

# 设置文件权限
RUN mkdir -p /app/data /var/log && chown -R appuser:appgroup /app /var/log

EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080 || exit 1

CMD [ "sh", "/app/start.sh"]