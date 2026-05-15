# Kubernetes 学习与实战指南

这份文档是写给刚接触 Kubernetes 的 Go 后端实习生看的。目标不是背概念，而是让你能理解 Kubernetes 为什么存在、它解决什么问题，以及如何把本项目这种 Go 微服务系统逐步部署到 Kubernetes 中。

本项目是一个智能人才运营平台，后端由多个 Go 服务组成，包括 `gateway`、`user-service`、`job-service`、`resume-service`、`talent-service`、`recommendation-service`、`interview-service`、`message-service`、`log-service` 和 `evaluator-service`，还依赖 PostgreSQL、Redis、Elasticsearch 等基础设施。`k8s` 目录目前已经包含：

- `base/namespace.yaml`：创建 `talent` 命名空间。
- `base/configmap.yaml`：保存服务公共配置，例如数据库地址、Redis 地址、各微服务访问地址。
- `infra/postgres.yaml`：部署 PostgreSQL，包括 PVC、Deployment 和 Service。
- `infra/postgres-init-configmap.yaml`：PostgreSQL 初始化 SQL。

后续你可以以这份文档为路线，把当前 Docker Compose 部署逐步迁移到 Kubernetes 部署。

## 1. 先理解 Kubernetes 是什么

Kubernetes 通常简称 K8s，它是一个容器编排平台。你可以把它理解为“专门管理容器化应用的操作系统”。

在没有 Kubernetes 之前，你可能会这样部署 Go 服务：

```bash
go build -o user-service
./user-service
```

或者用 Docker：

```bash
docker build -t user-service:latest .
docker run -p 8081:8081 user-service:latest
```

如果只有一个服务，这样还能接受。但真实项目里会有很多问题：

- 服务挂了谁来拉起？
- 多个服务如何互相访问？
- 配置和密码怎么管理？
- 服务升级时如何减少中断？
- PostgreSQL 这种有状态服务的数据怎么保留？
- 多台机器时容器应该调度到哪台机器？
- 如何扩容、缩容、查看日志、排查健康状态？

Kubernetes 解决的就是这些问题。

在 Kubernetes 里，你通常不直接运行容器，而是声明“我希望系统变成什么样”。比如你声明：

```yaml
replicas: 3
image: user-service:latest
```

Kubernetes 就会努力让集群里始终有 3 个 `user-service` 实例在运行。如果某个实例挂了，它会自动创建新的实例补上。

## 2. Kubernetes 的核心工作方式

Kubernetes 的核心思想是“声明式配置 + 控制器持续纠偏”。

你写 YAML 文件告诉 Kubernetes 期望状态：

- 我要一个命名空间。
- 我要一个 Deployment。
- 我要 3 个 Pod。
- 我要一个 Service 暴露访问入口。
- 我要一个 ConfigMap 保存普通配置。
- 我要一个 Secret 保存敏感配置。

然后执行：

```bash
kubectl apply -f xxx.yaml
```

Kubernetes 会把实际状态调整成你声明的状态。

如果你声明一个 Deployment 有 2 个副本，但现在只剩 1 个 Pod，Kubernetes 会自动再创建 1 个 Pod。

这和传统手动运维最大的区别是：你不再一步一步手工操作服务器，而是把系统结构写成 YAML，让 Kubernetes 持续维护这个结构。

## 3. 你必须掌握的几个核心对象

### 3.1 Namespace：命名空间

Namespace 用来隔离资源。比如本项目使用：

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: talent
```

这表示创建一个叫 `talent` 的命名空间。后续 PostgreSQL、Redis、后端服务、前端服务都可以放在这个命名空间里。

常用命令：

```bash
kubectl get ns
kubectl create namespace talent
kubectl get all -n talent
```

建议你所有本项目资源都加：

```yaml
metadata:
  namespace: talent
```

这样不会和其他项目混在一起。

### 3.2 Pod：Kubernetes 里最小的运行单元

Pod 是 Kubernetes 里最小的调度单位。一个 Pod 里面可以有一个或多个容器。

对于 Go 微服务来说，通常一个 Pod 里只放一个业务容器。例如：

- 一个 `user-service` Pod 里运行一个 `user-service` 容器。
- 一个 `job-service` Pod 里运行一个 `job-service` 容器。
- 一个 `postgres` Pod 里运行一个 PostgreSQL 容器。

你平时不会直接手写 Pod，因为 Pod 挂了以后不会自动重建。生产中更常用 Deployment 来管理 Pod。

查看 Pod：

```bash
kubectl get pods -n talent
kubectl describe pod <pod-name> -n talent
kubectl logs <pod-name> -n talent
```

进入 Pod：

```bash
kubectl exec -it <pod-name> -n talent -- sh
```

### 3.3 Deployment：管理无状态服务

Deployment 用来部署无状态应用，例如 Go API 服务、前端 Nginx 服务。

Deployment 负责：

- 创建 Pod。
- 保持副本数。
- 滚动更新。
- 失败自动拉起。
- 回滚版本。

一个典型的 Go 服务 Deployment 长这样：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  namespace: talent
spec:
  replicas: 1
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
        - name: user-service
          image: user-service:latest
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8081
          envFrom:
            - configMapRef:
                name: talent-config
            - secretRef:
                name: talent-secret
```

这里要注意三点：

1. `selector.matchLabels` 必须和 `template.metadata.labels` 对上。
2. `containerPort` 是容器内部监听端口，不等于外部访问端口。
3. `envFrom` 可以把 ConfigMap 和 Secret 里的键值注入为环境变量。

常用命令：

```bash
kubectl get deploy -n talent
kubectl describe deploy user-service -n talent
kubectl rollout status deploy/user-service -n talent
kubectl rollout restart deploy/user-service -n talent
kubectl rollout undo deploy/user-service -n talent
```

### 3.4 Service：让 Pod 有稳定访问地址

Pod 的 IP 会变化，所以服务之间不要直接访问 Pod IP。Kubernetes 用 Service 提供稳定的内部 DNS 名称。

例如本项目 `base/configmap.yaml` 里有：

```yaml
USER_SERVICE_URL: http://user-service:8081
JOB_SERVICE_URL: http://job-service:8082
```

这里的 `user-service`、`job-service` 就应该是 Kubernetes Service 的名字。

一个 Service 示例：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: user-service
  namespace: talent
spec:
  selector:
    app: user-service
  ports:
    - port: 8081
      targetPort: 8081
```

含义：

- `metadata.name: user-service`：集群内可以用 `http://user-service:8081` 访问。
- `selector.app: user-service`：把流量转发到标签为 `app=user-service` 的 Pod。
- `port`：Service 暴露的端口。
- `targetPort`：容器实际监听的端口。

常用命令：

```bash
kubectl get svc -n talent
kubectl describe svc user-service -n talent
```

### 3.5 ConfigMap：保存普通配置

ConfigMap 用来保存非敏感配置，比如数据库 host、端口、服务 URL。

本项目已有：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: talent-config
  namespace: talent
data:
  DB_HOST: postgres
  DB_PORT: "5432"
  DB_USER: postgres
  DB_NAME: talent_platform
  DB_SSLMODE: disable
  REDIS_HOST: redis
  REDIS_PORT: "6379"
  ES_URL: http://elasticsearch:9200
```

Go 服务读取这些配置时，通常用：

```go
os.Getenv("DB_HOST")
os.Getenv("DB_PORT")
```

使用 ConfigMap 的好处是：镜像不用因为配置变化重新构建。你只需要改 YAML，然后重启 Deployment。

常用命令：

```bash
kubectl get configmap -n talent
kubectl describe configmap talent-config -n talent
```

### 3.6 Secret：保存敏感配置

Secret 用来保存密码、Token、JWT 密钥等敏感信息。

当前 `infra/postgres.yaml` 引用了一个 Secret：

```yaml
valueFrom:
  secretKeyRef:
    name: talent-secret
    key: DB_PASSWORD
```

这说明你必须创建 `talent-secret`，否则 PostgreSQL Pod 会启动失败。

可以写一个 `base/secret.yaml`，示例：

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: talent-secret
  namespace: talent
type: Opaque
stringData:
  DB_PASSWORD: postgres
  JWT_SECRET: talent-platform-jwt-secret-2024
```

注意：这里用了 `stringData`，可以直接写明文，Kubernetes 会帮你转成 base64。实际生产不要把真实密码提交到 Git 仓库。

也可以用命令创建：

```bash
kubectl create secret generic talent-secret \
  --from-literal=DB_PASSWORD=postgres \
  --from-literal=JWT_SECRET=talent-platform-jwt-secret-2024 \
  -n talent
```

查看 Secret：

```bash
kubectl get secret -n talent
kubectl describe secret talent-secret -n talent
```

### 3.7 PVC：持久化数据

Pod 是临时的，删除后容器文件系统也会消失。数据库不能这样，所以 PostgreSQL 需要持久化卷。

本项目已有：

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-data
  namespace: talent
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
```

PVC 表示“我需要一块 2Gi 的持久化存储”。具体这块存储来自哪里，由集群的 StorageClass 决定。

本地 Docker Desktop、minikube、kind 的存储行为可能不同。你可以用：

```bash
kubectl get pvc -n talent
kubectl describe pvc postgres-data -n talent
kubectl get storageclass
```

如果 PVC 一直 `Pending`，说明集群没有可用的默认 StorageClass。

### 3.8 Ingress：对外暴露 HTTP 服务

Service 默认主要用于集群内部访问。如果你想用域名访问前端或 Gateway，一般会使用 Ingress。

例如：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: talent-ingress
  namespace: talent
spec:
  rules:
    - host: talent.local
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend
                port:
                  number: 80
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: gateway
                port:
                  number: 8080
```

Ingress 需要 Ingress Controller，例如 Nginx Ingress。没有 Controller，Ingress 对象创建了也不会生效。

本地学习阶段也可以先不用 Ingress，直接用端口转发：

```bash
kubectl port-forward svc/gateway 8080:8080 -n talent
kubectl port-forward svc/frontend 3000:80 -n talent
```

## 4. 本项目从 Docker Compose 到 Kubernetes 的映射

当前 `docker-compose.yml` 里主要有这些服务：

| Docker Compose 服务 | Kubernetes 资源建议 | 说明 |
|---|---|---|
| `postgres` | `Deployment` 或 `StatefulSet` + `PVC` + `Service` | 当前已有 Deployment + PVC |
| `redis` | `Deployment` + `PVC` + `Service` | 当前 K8s 目录还没有 |
| `elasticsearch` | `Deployment`/`StatefulSet` + `PVC` + `Service` | 当前 K8s 目录还没有 |
| `backend` | 多个 Go 服务 Deployment + Service | Compose 里是一个容器跑全部服务，K8s 更推荐拆分 |
| `frontend` | `Deployment` + `Service` | 当前 K8s 目录还没有 |

如果你只是毕业设计演示，可以先用“一个 backend 镜像跑全部服务”的方式迁移，简单但不够云原生。

如果你想真正学习 Kubernetes，建议每个 Go 微服务拆成独立 Deployment：

- `gateway`
- `user-service`
- `job-service`
- `resume-service`
- `talent-service`
- `recommendation-service`
- `interview-service`
- `message-service`
- `log-service`
- `evaluator-service`

这样你才能真正理解服务发现、独立扩缩容、滚动发布和故障隔离。

## 5. 推荐的目录结构

当前目录：

```text
k8s/
  base/
    namespace.yaml
    configmap.yaml
  infra/
    postgres.yaml
    postgres-init-configmap.yaml
```

建议逐步扩展成：

```text
k8s/
  README.md
  base/
    namespace.yaml
    configmap.yaml
    secret.example.yaml
  infra/
    postgres.yaml
    postgres-init-configmap.yaml
    redis.yaml
    elasticsearch.yaml
  services/
    gateway.yaml
    user-service.yaml
    job-service.yaml
    resume-service.yaml
    talent-service.yaml
    recommendation-service.yaml
    interview-service.yaml
    message-service.yaml
    log-service.yaml
    evaluator-service.yaml
  frontend/
    frontend.yaml
  ingress/
    ingress.yaml
```

这样分层更清晰：

- `base`：基础配置。
- `infra`：数据库、缓存、搜索引擎等基础设施。
- `services`：Go 微服务。
- `frontend`：前端应用。
- `ingress`：外部访问入口。

## 6. 第一次部署：先跑 PostgreSQL

### 6.1 创建命名空间

```bash
kubectl apply -f k8s/base/namespace.yaml
```

验证：

```bash
kubectl get ns
```

你应该看到：

```text
talent
```

### 6.2 创建 Secret

当前仓库还没有 `secret.yaml`，你可以先用命令创建：

```bash
kubectl create secret generic talent-secret \
  --from-literal=DB_PASSWORD=postgres \
  --from-literal=JWT_SECRET=talent-platform-jwt-secret-2024 \
  -n talent
```

验证：

```bash
kubectl get secret -n talent
```

### 6.3 创建 ConfigMap

```bash
kubectl apply -f k8s/base/configmap.yaml
```

验证：

```bash
kubectl get configmap -n talent
```

### 6.4 创建 PostgreSQL 初始化 SQL

```bash
kubectl apply -f k8s/infra/postgres-init-configmap.yaml
```

验证：

```bash
kubectl get configmap postgres-init-sql -n talent
```

### 6.5 部署 PostgreSQL

```bash
kubectl apply -f k8s/infra/postgres.yaml
```

查看状态：

```bash
kubectl get pods -n talent
kubectl get pvc -n talent
kubectl get svc -n talent
```

如果正常，你应该看到：

```text
postgres-xxxxx   Running
postgres-data    Bound
postgres         ClusterIP
```

查看日志：

```bash
kubectl logs deploy/postgres -n talent
```

进入 PostgreSQL：

```bash
kubectl exec -it deploy/postgres -n talent -- psql -U postgres -d talent_platform
```

查看表：

```sql
\dt
```

如果能看到 `users`、`jobs`、`talents` 等表，说明初始化 SQL 成功执行。

## 7. 部署一个 Go 微服务的完整思路

假设你要部署 `user-service`，完整步骤是：

1. 构建镜像。
2. 推送镜像到镜像仓库，或者本地集群直接使用本地镜像。
3. 写 Deployment。
4. 写 Service。
5. 应用 YAML。
6. 查看 Pod 日志。
7. 通过 Service 或 port-forward 访问。

### 7.1 构建镜像

本项目每个服务都有自己的 Dockerfile，例如：

```text
backend/user-service/Dockerfile
```

在 `backend` 目录构建镜像：

```bash
cd backend
docker build -t user-service:latest -f user-service/Dockerfile .
```

注意这里构建上下文是 `backend`，因为 Dockerfile 里有：

```dockerfile
COPY common ./common
COPY user-service ./user-service
```

如果你在错误目录构建，会找不到 `common`。

### 7.2 写 Deployment 和 Service

示例：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  namespace: talent
spec:
  replicas: 1
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
        - name: user-service
          image: user-service:latest
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8081
          envFrom:
            - configMapRef:
                name: talent-config
            - secretRef:
                name: talent-secret
          readinessProbe:
            httpGet:
              path: /health
              port: 8081
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /health
              port: 8081
            initialDelaySeconds: 15
            periodSeconds: 20
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
  namespace: talent
spec:
  selector:
    app: user-service
  ports:
    - port: 8081
      targetPort: 8081
```

这里的 `/health` 要看你的 Go 服务是否真的提供了这个接口。如果没有，就先不要加 HTTP Probe，或者改成实际存在的健康检查路径。

### 7.3 应用 YAML

```bash
kubectl apply -f k8s/services/user-service.yaml
```

查看：

```bash
kubectl get pods -n talent
kubectl logs deploy/user-service -n talent
```

如果 Pod 起不来，用：

```bash
kubectl describe pod <pod-name> -n talent
```

常见错误：

- `ImagePullBackOff`：镜像拉不到。
- `CrashLoopBackOff`：程序启动后崩溃。
- `CreateContainerConfigError`：ConfigMap 或 Secret 缺失。
- `Readiness probe failed`：健康检查失败。

## 8. Go 服务在 Kubernetes 里的配置原则

Go 服务部署到 Kubernetes 时，建议遵守这几个原则。

### 8.1 不要把配置写死在代码里

不要这样：

```go
dbHost := "localhost"
```

在 Kubernetes 里，服务之间不是 `localhost` 访问。`localhost` 只代表当前容器自己。

应该这样：

```go
dbHost := os.Getenv("DB_HOST")
```

然后在 ConfigMap 中配置：

```yaml
DB_HOST: postgres
```

### 8.2 服务之间用 Service 名称访问

在 Docker Compose 里你可以用服务名访问，比如 `postgres`。

在 Kubernetes 里也是类似的：

```text
postgres
user-service
job-service
resume-service
```

同一个 Namespace 下可以直接访问：

```text
http://user-service:8081
```

跨 Namespace 要写完整 DNS：

```text
http://user-service.talent.svc.cluster.local:8081
```

### 8.3 日志输出到 stdout/stderr

Go 服务不要只写本地文件。Kubernetes 最容易采集的是标准输出。

推荐：

```go
log.Println("server started")
```

然后查看：

```bash
kubectl logs deploy/user-service -n talent
```

如果一定要写文件，例如上传文件、临时解析文件，要明确挂载 Volume。

### 8.4 优雅退出

Kubernetes 删除 Pod 时会先发送 SIGTERM。Go 服务应该处理这个信号，优雅关闭 HTTP Server 和数据库连接。

示意：

```go
srv := &http.Server{Addr: ":8081", Handler: router}

go func() {
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
}()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
_ = srv.Shutdown(ctx)
```

这样滚动更新时不容易中断正在处理的请求。

## 9. 健康检查：Liveness 和 Readiness

Kubernetes 常见两个健康检查：

- `livenessProbe`：判断容器是否活着。失败后会重启容器。
- `readinessProbe`：判断容器是否可以接流量。失败后 Service 不会把流量转给它。

Go 服务最好提供：

```text
GET /health
```

返回：

```json
{"status":"ok"}
```

Kubernetes 配置：

```yaml
readinessProbe:
  httpGet:
    path: /health
    port: 8081
  initialDelaySeconds: 5
  periodSeconds: 10
livenessProbe:
  httpGet:
    path: /health
    port: 8081
  initialDelaySeconds: 15
  periodSeconds: 20
```

注意：健康检查接口不要依赖太重。例如数据库偶尔抖动时，是否要让容器重启，需要谨慎设计。

一般建议：

- readiness 可以检查数据库连接，因为数据库不可用时不该接流量。
- liveness 不要因为数据库短暂不可用就失败，否则可能造成服务反复重启。

## 10. 资源限制：requests 和 limits

每个服务最好设置资源请求和限制。

```yaml
resources:
  requests:
    cpu: "100m"
    memory: "128Mi"
  limits:
    cpu: "500m"
    memory: "512Mi"
```

含义：

- `requests`：调度时至少需要多少资源。
- `limits`：最多能用多少资源。

如果不设置，某个服务可能占满节点资源，影响其他服务。

Go 服务通常内存占用比较稳定，但要注意：

- 简历解析服务可能需要更多内存。
- AI 评估服务如果处理大文本，也可能需要更多内存。
- Elasticsearch 本身比较吃内存，本地学习时至少给 Docker Desktop 分配足够内存。

## 11. 滚动更新和回滚

更新镜像后：

```bash
kubectl set image deploy/user-service user-service=user-service:v2 -n talent
```

查看更新状态：

```bash
kubectl rollout status deploy/user-service -n talent
```

回滚：

```bash
kubectl rollout undo deploy/user-service -n talent
```

查看历史：

```bash
kubectl rollout history deploy/user-service -n talent
```

这就是 Kubernetes 比手动部署强的地方：它可以一边启动新版本，一边逐步替换旧版本，出问题还能回滚。

## 12. 扩容和缩容

如果 `gateway` 压力变大，可以扩容：

```bash
kubectl scale deploy/gateway --replicas=3 -n talent
```

查看：

```bash
kubectl get pods -n talent -l app=gateway
```

Service 会自动把流量分发到多个 Pod。

但是要注意：

- 无状态服务适合直接扩容。
- 有状态服务如 PostgreSQL 不能随便 `replicas: 3`。
- WebSocket 服务扩容时要考虑连接保持和会话路由。

本项目里，比较适合扩容的是：

- `gateway`
- `user-service`
- `job-service`
- `talent-service`
- `recommendation-service`

需要谨慎扩容的是：

- `message-service`，因为可能涉及 WebSocket。
- `resume-service`，因为可能涉及上传文件和本地临时文件。
- PostgreSQL、Redis、Elasticsearch，这些有状态组件需要专门方案。

## 13. 排错流程

Kubernetes 排错不要乱猜，按顺序看。

### 13.1 看资源是否存在

```bash
kubectl get all -n talent
```

### 13.2 看 Pod 状态

```bash
kubectl get pods -n talent
```

常见状态：

- `Running`：正在运行。
- `Pending`：调度不了，可能资源不足或 PVC 未绑定。
- `CrashLoopBackOff`：程序反复崩溃。
- `ImagePullBackOff`：镜像拉取失败。
- `CreateContainerConfigError`：配置缺失。

### 13.3 看详细事件

```bash
kubectl describe pod <pod-name> -n talent
```

重点看最后的 `Events`。

### 13.4 看日志

```bash
kubectl logs <pod-name> -n talent
```

如果 Pod 有多个容器：

```bash
kubectl logs <pod-name> -c <container-name> -n talent
```

如果容器刚崩溃过：

```bash
kubectl logs <pod-name> --previous -n talent
```

### 13.5 进入容器检查

```bash
kubectl exec -it <pod-name> -n talent -- sh
```

检查环境变量：

```bash
env | sort
```

检查 DNS：

```bash
nslookup postgres
```

Alpine 镜像可能没有 `nslookup`，可以临时创建 debug Pod：

```bash
kubectl run debug -it --rm --image=busybox:1.36 -n talent -- sh
```

在 debug Pod 里：

```sh
nslookup postgres
wget -qO- http://user-service:8081/health
```

## 14. 本项目的推荐实战路线

你可以按下面路线学习，不要一口气全部上 K8s。

### 第一阶段：理解基础设施

目标：让 PostgreSQL 在 K8s 中跑起来。

步骤：

```bash
kubectl apply -f k8s/base/namespace.yaml
kubectl create secret generic talent-secret \
  --from-literal=DB_PASSWORD=postgres \
  --from-literal=JWT_SECRET=talent-platform-jwt-secret-2024 \
  -n talent
kubectl apply -f k8s/base/configmap.yaml
kubectl apply -f k8s/infra/postgres-init-configmap.yaml
kubectl apply -f k8s/infra/postgres.yaml
kubectl get pods -n talent
```

验收标准：

- PostgreSQL Pod 是 `Running`。
- PVC 是 `Bound`。
- 能进入数据库并看到业务表。

### 第二阶段：补 Redis 和 Elasticsearch

目标：把 Docker Compose 里的基础设施迁移完整。

你需要新增：

- `k8s/infra/redis.yaml`
- `k8s/infra/elasticsearch.yaml`

Redis 最小版本可以先用 Deployment + Service。

Elasticsearch 本地学习可以用单节点模式：

```yaml
env:
  - name: discovery.type
    value: single-node
  - name: xpack.security.enabled
    value: "false"
  - name: ES_JAVA_OPTS
    value: "-Xms512m -Xmx512m"
```

验收标准：

```bash
kubectl port-forward svc/elasticsearch 9200:9200 -n talent
curl http://localhost:9200/_cluster/health
```

### 第三阶段：部署 Gateway

目标：先让一个 Go 服务跑起来。

建议先部署 `gateway`，因为它是对外入口。

构建镜像：

```bash
cd backend
docker build -t gateway:latest -f gateway/Dockerfile .
```

部署 Deployment + Service。

访问：

```bash
kubectl port-forward svc/gateway 8080:8080 -n talent
```

浏览器或 curl：

```bash
curl http://localhost:8080
```

### 第四阶段：部署一个业务服务

建议从 `user-service` 开始，因为认证和用户是基础。

构建：

```bash
cd backend
docker build -t user-service:latest -f user-service/Dockerfile .
```

部署：

```bash
kubectl apply -f k8s/services/user-service.yaml
```

验证：

```bash
kubectl logs deploy/user-service -n talent
kubectl get svc -n talent
```

### 第五阶段：逐步拆分全部服务

服务端口参考：

| 服务 | 端口 | 说明 |
|---|---:|---|
| `gateway` | 8080 | API 网关 |
| `user-service` | 8081 | 用户认证与权限 |
| `job-service` | 8082 | 职位管理 |
| `interview-service` | 8083 | 面试管理 |
| `resume-service` | 8084 | 简历上传与解析 |
| `message-service` | 8085 | 消息与 WebSocket |
| `talent-service` | 8086 | 人才库 |
| `recommendation-service` | 8087 | 推荐匹配 |
| `log-service` | 8088 | 操作日志 |
| `evaluator-service` | 8090 | AI 评估 |

每个服务都应该有：

- Deployment
- Service
- 环境变量
- 日志查看方式
- 健康检查
- 资源限制

### 第六阶段：部署前端

构建：

```bash
cd frontend
docker build -t talent-frontend:latest .
```

部署：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
  namespace: talent
spec:
  replicas: 1
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      containers:
        - name: frontend
          image: talent-frontend:latest
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: frontend
  namespace: talent
spec:
  selector:
    app: frontend
  ports:
    - port: 80
      targetPort: 80
```

访问：

```bash
kubectl port-forward svc/frontend 3000:80 -n talent
```

浏览器打开：

```text
http://localhost:3000
```

## 15. 你应该重点练习的 kubectl 命令

基础查看：

```bash
kubectl get ns
kubectl get all -n talent
kubectl get pods -n talent -o wide
kubectl get svc -n talent
kubectl get deploy -n talent
kubectl get pvc -n talent
```

详情排错：

```bash
kubectl describe pod <pod-name> -n talent
kubectl describe deploy <deploy-name> -n talent
kubectl describe svc <svc-name> -n talent
```

日志：

```bash
kubectl logs <pod-name> -n talent
kubectl logs deploy/user-service -n talent
kubectl logs deploy/user-service -n talent --tail=100 -f
```

进入容器：

```bash
kubectl exec -it <pod-name> -n talent -- sh
```

应用配置：

```bash
kubectl apply -f k8s/base/namespace.yaml
kubectl apply -f k8s/base/
kubectl apply -f k8s/infra/
kubectl apply -f k8s/services/
```

删除资源：

```bash
kubectl delete -f k8s/services/user-service.yaml
kubectl delete namespace talent
```

端口转发：

```bash
kubectl port-forward svc/postgres 5432:5432 -n talent
kubectl port-forward svc/gateway 8080:8080 -n talent
kubectl port-forward svc/frontend 3000:80 -n talent
```

滚动更新：

```bash
kubectl rollout status deploy/user-service -n talent
kubectl rollout restart deploy/user-service -n talent
kubectl rollout undo deploy/user-service -n talent
```

扩缩容：

```bash
kubectl scale deploy/user-service --replicas=2 -n talent
```

## 16. 常见问题

### 16.1 为什么服务里不能访问 localhost:5432？

因为在容器里，`localhost` 指的是当前容器自己，不是 PostgreSQL 容器。

应该访问：

```text
postgres:5432
```

前提是你有一个名为 `postgres` 的 Service。

### 16.2 为什么 Pod 是 ImagePullBackOff？

原因通常是：

- 镜像名写错。
- 本地集群找不到本地镜像。
- 私有镜像仓库没有配置认证。

如果用 kind，需要把镜像加载进 kind：

```bash
kind load docker-image user-service:latest
```

如果用 minikube，可以在 minikube 的 Docker 环境里构建：

```bash
eval $(minikube docker-env)
docker build -t user-service:latest -f backend/user-service/Dockerfile backend
```

### 16.3 为什么 Pod 是 CrashLoopBackOff？

说明容器启动后程序崩了。

看日志：

```bash
kubectl logs <pod-name> -n talent --previous
```

常见原因：

- 数据库连接失败。
- 环境变量缺失。
- 配置里的服务地址不对。
- 程序监听端口和 YAML 里写的不一致。

### 16.4 为什么 PostgreSQL 初始化 SQL 没执行？

PostgreSQL 官方镜像只会在数据目录为空时执行 `/docker-entrypoint-initdb.d`。

如果 PVC 已经有旧数据，重新 apply ConfigMap 不会重新初始化。

解决学习环境里可以删除 PVC：

```bash
kubectl delete deploy postgres -n talent
kubectl delete pvc postgres-data -n talent
kubectl apply -f k8s/infra/postgres.yaml
```

注意：删除 PVC 会删除数据库数据，生产环境不要随便做。

### 16.5 为什么 ConfigMap 改了服务没生效？

环境变量是在容器启动时注入的。ConfigMap 改了以后，已经运行的 Pod 不会自动刷新环境变量。

你需要重启 Deployment：

```bash
kubectl rollout restart deploy/user-service -n talent
```

### 16.6 为什么 Service 访问不到 Pod？

检查 Service selector 是否和 Pod label 匹配：

```bash
kubectl describe svc user-service -n talent
kubectl get pods -n talent --show-labels
```

Service 写的是：

```yaml
selector:
  app: user-service
```

Pod 必须有：

```yaml
labels:
  app: user-service
```

否则 Service 后面没有 endpoints。

检查：

```bash
kubectl get endpoints -n talent
```

## 17. 生产环境要考虑什么

学习环境能跑起来只是第一步。真正生产环境还要考虑：

- 镜像版本不能一直用 `latest`，要用明确版本号。
- Secret 不要明文提交 Git。
- 数据库不要用单副本 Deployment，应考虑云数据库或 StatefulSet。
- 服务要有资源限制。
- 要有日志采集和监控告警。
- Ingress 要配置 TLS。
- 敏感接口要有限流和鉴权。
- 上传文件要考虑对象存储，而不是只放容器本地目录。
- 发布要有 CI/CD 流程。

本项目如果继续往生产化演进，建议路线：

1. 镜像推送到镜像仓库。
2. Kubernetes YAML 使用固定镜像 tag。
3. 用 GitHub Actions 或其他 CI/CD 自动构建镜像。
4. 使用 Helm 或 Kustomize 管理多环境配置。
5. 数据库迁移从 ConfigMap 初始化 SQL 改为版本化 migration。
6. 监控接入 Prometheus + Grafana。
7. 日志接入 Loki 或 Elasticsearch。

## 18. 给 Go 实习生的学习建议

不要一开始就追求“全部上云原生最佳实践”。你可以按下面顺序练：

1. 会写 Dockerfile。
2. 会构建镜像并运行容器。
3. 会用 Kubernetes 跑一个 Deployment。
4. 会用 Service 做服务发现。
5. 会用 ConfigMap 和 Secret 管理配置。
6. 会用 PVC 跑 PostgreSQL。
7. 会看 Pod 日志和 Events。
8. 会处理 `CrashLoopBackOff`、`ImagePullBackOff`。
9. 会滚动更新和回滚。
10. 会把一个真实 Go 微服务部署到 Kubernetes。

最重要的是：每次只解决一个问题。先让 PostgreSQL 跑起来，再让一个 Go 服务连上数据库，再部署 Gateway，最后再接前端和 Ingress。

## 19. 最小实战清单

你可以用下面清单判断自己是否真的掌握了本项目的 K8s 部署。

- [ ] 我能解释 Pod、Deployment、Service 的区别。
- [ ] 我能解释 ConfigMap 和 Secret 的区别。
- [ ] 我能解释为什么数据库需要 PVC。
- [ ] 我能在 `talent` 命名空间中部署 PostgreSQL。
- [ ] 我能进入 PostgreSQL Pod 并查看业务表。
- [ ] 我能构建一个 Go 服务镜像。
- [ ] 我能写出一个 Go 服务的 Deployment 和 Service。
- [ ] 我能通过 Service 名称访问另一个服务。
- [ ] 我能通过 `kubectl logs` 查服务日志。
- [ ] 我能用 `kubectl describe pod` 看失败原因。
- [ ] 我能处理一次 `CrashLoopBackOff`。
- [ ] 我能处理一次 `ImagePullBackOff`。
- [ ] 我能执行一次滚动更新。
- [ ] 我能回滚到上一个版本。
- [ ] 我能用 `port-forward` 在本地访问 Gateway 或前端。

完成这些，你就不是“只知道 Kubernetes 名词”，而是真的能把 Go 项目部署到 Kubernetes 里了。
