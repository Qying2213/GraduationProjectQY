# 智能人才招聘平台 - 部署文档

## 目录
1. [系统要求](#系统要求)
2. [环境准备](#环境准备)
3. [本地开发部署](#本地开发部署)
4. [Docker部署](#docker部署)
5. [生产环境部署](#生产环境部署)
6. [常见问题](#常见问题)

---

## 系统要求

### 硬件要求
- CPU: 2核以上
- 内存: 4GB以上
- 硬盘: 20GB以上

### 软件要求
- 操作系统: Linux/macOS/Windows
- Node.js: 18.x 或更高版本
- Go: 1.21 或更高版本
- PostgreSQL: 14.x 或更高版本
- Redis: 6.x 或更高版本
- Docker: 20.x 或更高版本（可选）

---

## 环境准备

### 1. 安装 Node.js
```bash
# macOS (使用 Homebrew)
brew install node

# Ubuntu/Debian
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 验证安装
node --version
npm --version
```

### 2. 安装 Go
```bash
# macOS
brew install go

# Ubuntu/Debian
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 验证安装
go version
```

### 3. 安装 PostgreSQL
```bash
# macOS
brew install postgresql@14
brew services start postgresql@14

# Ubuntu/Debian
sudo apt-get install postgresql-14

# 创建数据库
psql -U postgres
CREATE DATABASE talent_platform;
CREATE USER talent_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE talent_platform TO talent_user;
```

### 4. 安装 Redis
```bash
# macOS
brew install redis
brew services start redis

# Ubuntu/Debian
sudo apt-get install redis-server
sudo systemctl start redis
```

---

## 本地开发部署

### 1. 克隆项目
```bash
git clone <repository-url>
cd talent-platform
```

### 2. 配置环境变量
```bash
# 复制环境变量模板
cp backend/.env.example backend/.env

# 编辑配置文件
vim backend/.env
```

`.env` 文件内容：
```env
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=talent_user
DB_PASSWORD=your_password
DB_NAME=talent_platform

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT配置
JWT_SECRET=your-jwt-secret-key

# 服务端口
GATEWAY_PORT=8080
USER_SERVICE_PORT=8081
JOB_SERVICE_PORT=8082
TALENT_SERVICE_PORT=8083
RESUME_SERVICE_PORT=8084
```

### 3. 初始化数据库
```bash
# 执行数据库脚本
psql -U talent_user -d talent_platform -f backend/database/schema.sql
psql -U talent_user -d talent_platform -f backend/database/mock_data.sql
```

### 4. 启动后端服务
```bash
# 方式一：使用启动脚本
chmod +x start-backend.sh
./start-backend.sh

# 方式二：手动启动各服务
cd backend/gateway && go run main.go &
cd backend/user-service && go run main.go &
cd backend/job-service && go run main.go &
cd backend/talent-service && go run main.go &
cd backend/resume-service && go run main.go &
```

### 5. 启动前端
```bash
cd frontend
npm install
npm run dev
```

### 6. 访问系统
- 前台求职端: http://localhost:3000/portal
- 后台管理: http://localhost:3000/login
- API网关: http://localhost:8080

### 默认账号
| 角色 | 用户名 | 密码 |
|------|--------|------|
| 管理员 | admin | password123 |
| HR主管 | hr_zhang | password123 |
| 只读用户 | viewer_test | password123 |

---

## Docker部署

### 1. 构建镜像
```bash
# 构建所有服务
docker-compose build
```

### 2. 启动服务
```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 3. 停止服务
```bash
docker-compose down
```

### docker-compose.yml 配置说明
```yaml
version: '3.8'
services:
  # PostgreSQL 数据库
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: talent_platform
      POSTGRES_USER: talent_user
      POSTGRES_PASSWORD: your_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  # Redis 缓存
  redis:
    image: redis:6-alpine
    ports:
      - "6379:6379"

  # API 网关
  gateway:
    build: ./backend/gateway
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis

  # 前端服务
  frontend:
    build: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - gateway
```

---

## 生产环境部署

### 1. 服务器配置

#### Nginx 配置
```nginx
# /etc/nginx/sites-available/talent-platform
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /var/www/talent-platform/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. 构建前端
```bash
cd frontend
npm run build
```

### 3. 构建后端
```bash
cd backend/gateway
go build -o gateway main.go

cd ../user-service
go build -o user-service main.go

# ... 其他服务类似
```

### 4. 使用 Systemd 管理服务
```ini
# /etc/systemd/system/talent-gateway.service
[Unit]
Description=Talent Platform Gateway
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/talent-platform/backend/gateway
ExecStart=/opt/talent-platform/backend/gateway/gateway
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
# 启动服务
sudo systemctl enable talent-gateway
sudo systemctl start talent-gateway
```

---

## 常见问题

### Q1: 数据库连接失败
```
检查步骤：
1. 确认 PostgreSQL 服务已启动
2. 检查 .env 中的数据库配置
3. 确认数据库用户权限
```

### Q2: 前端无法访问后端 API
```
检查步骤：
1. 确认后端服务已启动
2. 检查 CORS 配置
3. 确认 API 地址配置正确
```

### Q3: Docker 容器启动失败
```bash
# 查看容器日志
docker-compose logs <service-name>

# 重新构建
docker-compose build --no-cache
```

### Q4: 端口被占用
```bash
# 查找占用端口的进程
lsof -i :8080

# 终止进程
kill -9 <PID>
```

---

## 技术支持

如有问题，请提交 Issue 或联系开发者。
