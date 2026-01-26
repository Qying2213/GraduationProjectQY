# CI/CD 流水线实现教程

本教程教你如何为智能招聘系统配置 GitHub Actions CI/CD 流水线。

---

## 第一步：配置 GitHub Secrets

### 1.1 进入 Secrets 配置页面

1. 打开你的 GitHub 仓库：`https://github.com/Qying2213/GraduationProjectQY`
2. 点击 **Settings** (设置)
3. 左侧菜单找到 **Secrets and variables** → **Actions**
4. 点击 **New repository secret**

### 1.2 添加以下 Secrets

你需要添加这些 Secret（一个一个添加）：

| Name | Value | 说明 |
|------|-------|------|
| `DOCKER_USERNAME` | 你的 Docker Hub 用户名 | 用于推送镜像 |
| `DOCKER_PASSWORD` | Docker Hub Access Token | 在 Docker Hub 生成 |
| `SERVER_HOST` | 服务器 IP 地址 | 部署目标服务器 |
| `SERVER_USER` | 服务器用户名（如 root） | SSH 登录用户 |
| `SERVER_SSH_KEY` | SSH 私钥内容 | 用于免密登录服务器 |
| `COZE_TOKEN` | 你的 Coze Token | AI 评估用 |
| `COZE_WORKFLOW_ID` | Coze 工作流 ID | AI 评估用 |
| `ARK_API_KEY` | 火山引擎 API Key | Embedding 用 |

### 1.3 获取 Docker Hub Access Token

1. 登录 https://hub.docker.com
2. 点击右上角头像 → **Account Settings**
3. 左侧选择 **Security**
4. 点击 **New Access Token**
5. 输入名称，选择 **Read & Write** 权限
6. 复制生成的 Token，添加到 GitHub Secrets

---

## 第二步：创建 Workflow 文件

### 2.1 创建目录

在你的项目根目录创建文件夹：

```
.github/
└── workflows/
```

### 2.2 创建 CI 流水线文件

创建文件 `.github/workflows/ci.yml`，内容如下：

```yaml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Build Backend Services
        run: |
          cd backend/gateway && go build -o gateway .
          cd ../user-service && go build -o user-service .
          cd ../job-service && go build -o job-service .
          # 添加其他服务...
      
      - name: Build Frontend
        run: |
          cd frontend
          npm ci
          npm run build
```

### 2.3 创建 CD 流水线文件

创建文件 `.github/workflows/cd.yml`，内容如下：

```yaml
name: CD

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      
      - name: Build and Push Images
        run: |
          docker build -t ${{ secrets.DOCKER_USERNAME }}/talent-gateway:latest ./backend/gateway
          docker push ${{ secrets.DOCKER_USERNAME }}/talent-gateway:latest
          # 重复其他服务...
      
      - name: Deploy to Server
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SERVER_SSH_KEY }}
          script: |
            cd /opt/talent-platform
            docker-compose pull
            docker-compose up -d
```

---

## 第三步：配置服务器

### 3.1 安装 Docker

SSH 登录你的服务器，执行：

```bash
# 安装 Docker
curl -fsSL https://get.docker.com | sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 3.2 创建部署目录

```bash
sudo mkdir -p /opt/talent-platform
cd /opt/talent-platform
```

### 3.3 创建 docker-compose.yml

在服务器上创建 `/opt/talent-platform/docker-compose.yml`：

```yaml
version: '3.8'

services:
  postgres:
    image: pgvector/pgvector:pg16
    environment:
      POSTGRES_USER: qinyang
      POSTGRES_DB: talent_platform
      POSTGRES_HOST_AUTH_METHOD: trust
    volumes:
      - postgres_data:/var/lib/postgresql/data

  gateway:
    image: your-dockerhub-username/talent-gateway:latest
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres

  # 添加其他服务...

volumes:
  postgres_data:
```

### 3.4 配置 SSH 密钥

在本地生成密钥对：

```bash
ssh-keygen -t rsa -b 4096 -f ~/.ssh/deploy_key -N ""
```

将公钥添加到服务器：

```bash
cat ~/.ssh/deploy_key.pub | ssh user@your-server "cat >> ~/.ssh/authorized_keys"
```

将私钥内容添加到 GitHub Secrets 的 `SERVER_SSH_KEY`。

---

## 第四步：测试流水线

### 4.1 提交代码触发 CI

```bash
git add .github/workflows/
git commit -m "Add CI/CD workflows"
git push origin main
```

### 4.2 查看运行状态

1. 打开 GitHub 仓库
2. 点击 **Actions** 标签
3. 查看 workflow 运行状态

### 4.3 手动触发部署

在 Actions 页面，选择 CD workflow，点击 **Run workflow**。

---

## 第五步：添加状态徽章（可选）

在 README.md 中添加：

```markdown
![CI](https://github.com/Qying2213/GraduationProjectQY/actions/workflows/ci.yml/badge.svg)
```

---

## 常见问题

### Q: 构建失败怎么办？
A: 点击失败的 job，查看详细日志，根据错误信息修复。

### Q: Docker 推送失败？
A: 检查 DOCKER_USERNAME 和 DOCKER_PASSWORD 是否正确配置。

### Q: SSH 连接失败？
A: 确保私钥格式正确，服务器防火墙开放 22 端口。

---

## 参考资源

- [GitHub Actions 官方文档](https://docs.github.com/en/actions)
- [Docker Hub 文档](https://docs.docker.com/docker-hub/)
- [appleboy/ssh-action](https://github.com/appleboy/ssh-action)
