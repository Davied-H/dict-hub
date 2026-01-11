.PHONY: build frontend backend clean dev help docker-build docker-push release

# Docker 配置（可通过环境变量覆盖）
DOCKER_IMAGE ?= dongdonghe122/dict-hub
VERSION ?= latest

# 默认目标
all: build

# 帮助信息
help:
	@echo "Dict-Hub 构建命令："
	@echo ""
	@echo "本地构建："
	@echo "  make build       - 完整构建（前端 + 后端）"
	@echo "  make frontend    - 仅构建前端"
	@echo "  make backend     - 仅构建后端"
	@echo "  make clean       - 清理构建产物"
	@echo "  make dev         - 开发模式（分别启动前后端）"
	@echo ""
	@echo "Docker 命令："
	@echo "  make docker-build           - 构建 Docker 镜像"
	@echo "  make docker-push            - 推送镜像到 Docker Hub"
	@echo "  make release VERSION=x.x.x  - 一键构建并推送新版本"
	@echo ""
	@echo "示例："
	@echo "  make release VERSION=1.0.0"
	@echo "  make docker-build DOCKER_IMAGE=myuser/myapp VERSION=2.0.0"

# 构建前端
frontend:
	@echo "==> 构建前端..."
	cd frontend && npm ci && npm run build
	@echo "==> 前端构建完成"

# 构建后端
backend:
	@echo "==> 构建后端..."
	cd backend && go build -o dict-hub ./cmd/server
	@echo "==> 后端构建完成: backend/dict-hub"

# 完整构建：先前端后后端
build: frontend backend
	@echo ""
	@echo "=========================================="
	@echo "构建完成！"
	@echo "运行: ./backend/dict-hub"
	@echo "=========================================="

# 清理构建产物
clean:
	@echo "==> 清理构建产物..."
	rm -rf backend/web/dist/*
	rm -f backend/dict-hub
	touch backend/web/dist/.gitkeep
	@echo "==> 清理完成"

# 开发模式
dev:
	@echo "开发模式：请分别在两个终端运行："
	@echo "  终端1: cd frontend && npm run dev"
	@echo "  终端2: cd backend && go run ./cmd/server"

# =============================================================================
# Docker 命令
# =============================================================================

# 构建 Docker 镜像
docker-build:
	@echo "==> 构建 Docker 镜像: $(DOCKER_IMAGE):$(VERSION)"
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	@if [ "$(VERSION)" != "latest" ]; then \
		echo "==> 同时标记为 latest"; \
		docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest; \
	fi
	@echo "==> Docker 镜像构建完成"

# 推送镜像到 Docker Hub
docker-push:
	@echo "==> 推送镜像到 Docker Hub: $(DOCKER_IMAGE):$(VERSION)"
	docker push $(DOCKER_IMAGE):$(VERSION)
	@if [ "$(VERSION)" != "latest" ]; then \
		echo "==> 同时推送 latest 标签"; \
		docker push $(DOCKER_IMAGE):latest; \
	fi
	@echo "==> 镜像推送完成"

# 一键构建并推送新版本
release: docker-build docker-push
	@echo ""
	@echo "=========================================="
	@echo "发布完成！"
	@echo "镜像: $(DOCKER_IMAGE):$(VERSION)"
	@echo ""
	@echo "用户可以通过以下命令拉取："
	@echo "  docker pull $(DOCKER_IMAGE):$(VERSION)"
	@echo "  docker pull $(DOCKER_IMAGE):latest"
	@echo "=========================================="
