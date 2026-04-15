.PHONY: all build run dev test clean docker-up docker-down swag

# 变量
APP_NAME := rental-management
MAIN_FILE := cmd/server/main.go
BUILD_DIR := bin
DOCKER_COMPOSE := docker-compose.yml

# Go参数
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# 构建目录
BUILD_TIME := $(shell date +%Y%m%d_%H%M%S)
LDFLAGS := -ldflags "-s -w"

all: clean deps build

# 编译
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

# 开发运行
dev:
	@echo "Running in development mode..."
	$(GOCMD) run $(MAIN_FILE)

# 运行
run: build
	@echo "Running..."
	./$(BUILD_DIR)/$(APP_NAME)

# 测试
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# 清理
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	$(GOCLEAN)

# 依赖
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Docker
docker-up:
	@echo "Starting Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE) up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE) down

docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .

docker-logs:
	docker-compose -f $(DOCKER_COMPOSE) logs -f

# Swagger文档
swag:
	@echo "Generating Swagger docs..."
	swag init -g $(MAIN_FILE) -o api/swagger

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# 格式化代码
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# 数据库迁移
migrate-up:
	@echo "Running migrations..."
	@mysql -h localhost -u root -p < scripts/init_db.sql

# 帮助
help:
	@echo "Available targets:"
	@echo "  make build      - Build the application"
	@echo "  make dev        - Run in development mode"
	@echo "  make run        - Build and run"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make deps       - Download dependencies"
	@echo "  make docker-up  - Start Docker containers"
	@echo "  make docker-down- Stop Docker containers"
	@echo "  make swag       - Generate Swagger docs"
	@echo "  make lint       - Run linter"
	@echo "  make fmt        - Format code"
