# GitHub Actions 配置指南

## 1. 配置 GitHub Secrets

在 GitHub 仓库中配置以下 Secrets：

### 1.1 进入 Settings → Secrets and variables → Actions

添加以下 Secrets：

| Secret 名称 | 说明 | 示例值 |
|------------|------|--------|
| `DOCKER_USERNAME` | Docker Hub 用户名 | your-dockerhub-username |
| `DOCKER_PASSWORD` | Docker Hub 密码/Token | your-access-token |
| `SERVER_HOST` | 部署服务器 IP | 192.168.1.100 |
| `SERVER_USER` | 服务器用户名 | root |
| `SERVER_SSH_KEY` | SSH 私钥 | -----BEGIN RSA PRIVATE KEY----- |
| `DB_PASSWORD` | 数据库密码 | your-db-password |
| `COZE_TOKEN` | Coze API Token | pat_xxx |
| `COZE_WORKFLOW_ID` | Coze 工作流 ID | 7583886563420373019 |
| `ARK_API_KEY` | 火山引擎 API Key | your-ark-api-key |

### 1.2 配置步骤截图说明

```
GitHub 仓库页面
    │
    ├── Settings (设置)
    │       │
    │       ├── Secrets and variables
    │       │       │
    │       │       └── Actions
    │       │               │
    │       │               ├── New repository secret
    │       │               │
    │       │               └── 添加上述 Secrets
```

## 2. 创建 Workflow 文件

### 2.1 CI 流水线 (.github/workflows/ci.yml)

```yaml
name: CI Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  GO_VERSION: '1.21'
  NODE_VERSION: '18'

jobs:
  # 后端代码检查和测试
  backend-lint-test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Cache Go modules
        uses: actions/cache@v4
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Install dependencies
        run: |
          cd backend/common && go mod download
          cd ../user-service && go mod download
          cd ../job-service && go mod download
          cd ../resume-service && go mod download
          cd ../talent-service && go mod download
          cd ../interview-service && go mod download
          cd ../message-service && go mod download
          cd ../recommendation-service && go mod download
          cd ../gateway && go mod download

      - name: Run Go Vet
        run: |
          cd backend/user-service && go vet ./...
          cd ../job-service && go vet ./...
          cd ../resume-service && go vet ./...

      - name: Run Tests
        run: |
          cd backend/user-service && go test -v ./...
          cd ../interview-service && go test -v ./...
          cd ../recommendation-service && go test -v ./...

  # 前端代码检查和测试
  frontend-lint-test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json

      - name: Install dependencies
        run: cd frontend && npm ci

      - name: Run Lint
        run: cd frontend && npm run lint || true

      - name: Run Tests
        run: cd frontend && npm run test || true

      - name: Build
        run: cd frontend && npm run build

  # 构建检查
  build-check:
    needs: [backend-lint-test, frontend-lint-test]
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Build all services
        run: |
          cd backend/gateway && go build -o gateway .
          cd ../user-service && go build -o user-service .
          cd ../job-service && go build -o job-service .
          cd ../resume-service && go build -o resume-service .
          cd ../talent-service && go build -o talent-service .
          cd ../interview-service && go build -o interview-service .
          cd ../message-service && go build -o message-service .
          cd ../recommendation-service && go build -o recommendation-service .

      - name: Build success
        run: echo "All services built successfully!"
```

### 2.2 CD 流水线 (.github/workflows/cd.yml)

```yaml
name: CD Pipeline

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]
  workflow_dispatch:

env:
  REGISTRY: docker.io
  IMAGE_PREFIX: ${{ secrets.DOCKER_USERNAME }}/talent-platform

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        service:
          - gateway
          - user-service
          - job-service
          - resume-service
          - talent-service
          - interview-service
          - message-service
          - recommendation-service
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.IMAGE_PREFIX }}-${{ matrix.service }}
          tags: |
            type=ref,event=branch
            type=ref,event=tag
            type=sha,prefix=
            type=raw,value=latest,enable=${{ github.ref == 'refs/heads/main' }}

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: ./backend/${{ matrix.service }}
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  build-frontend:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and push frontend
        uses: docker/build-push-action@v5
        with:
          context: ./frontend
          push: true
          tags: |
            ${{ env.IMAGE_PREFIX }}-frontend:latest
            ${{ env.IMAGE_PREFIX }}-frontend:${{ github.sha }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy:
    needs: [build-and-push, build-frontend]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Deploy to server
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SERVER_SSH_KEY }}
          script: |
            cd /opt/talent-platform
            docker-compose pull
            docker-compose up -d
            docker system prune -f
```

## 3. 手动触发部署

在 GitHub Actions 页面，选择 CD Pipeline，点击 "Run workflow" 即可手动触发部署。

## 4. 查看构建状态

### 4.1 添加状态徽章到 README

```markdown
![CI](https://github.com/Qying2213/GraduationProjectQY/actions/workflows/ci.yml/badge.svg)
![CD](https://github.com/Qying2213/GraduationProjectQY/actions/workflows/cd.yml/badge.svg)
```

### 4.2 查看构建日志

1. 进入仓库的 Actions 页面
2. 选择对应的 workflow
3. 点击具体的 run 查看详细日志

## 5. 分支策略

```
main (生产环境)
  │
  ├── develop (开发环境)
  │     │
  │     ├── feature/xxx (功能分支)
  │     │
  │     └── bugfix/xxx (修复分支)
  │
  └── hotfix/xxx (紧急修复)
```

- `main`: 生产环境，只接受 PR 合并
- `develop`: 开发环境，日常开发
- `feature/*`: 新功能开发
- `bugfix/*`: Bug 修复
- `hotfix/*`: 紧急修复，直接合并到 main
