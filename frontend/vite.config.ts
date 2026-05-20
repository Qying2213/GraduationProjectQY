import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import fs from 'node:fs'
import path from 'node:path'

const projectRoot = fileURLToPath(new URL('..', import.meta.url))
const devLogDir = path.join(projectRoot, 'logs/dev')
const devPidDir = path.join(projectRoot, 'tmp/pids')

const devServices = [
    { name: 'user-service', port: 8081, url: 'http://localhost:8081' },
    { name: 'job-service', port: 8082, url: 'http://localhost:8082' },
    { name: 'interview-service', port: 8083, url: 'http://localhost:8083' },
    { name: 'resume-service', port: 8084, url: 'http://localhost:8084' },
    { name: 'message-service', port: 8085, url: 'http://localhost:8085' },
    { name: 'talent-service', port: 8086, url: 'http://localhost:8086' },
    { name: 'recommendation-service', port: 8087, url: 'http://localhost:8087' },
    { name: 'log-service', port: 8088, url: 'http://localhost:8088' },
    { name: 'evaluator-service', port: 8090, url: 'http://localhost:8090' },
    { name: 'gateway', port: 8080, url: 'http://localhost:8080/health' },
    { name: 'frontend', port: 5173, url: 'http://localhost:5173' }
]

const sendJson = (res: any, statusCode: number, body: unknown) => {
    res.statusCode = statusCode
    res.setHeader('content-type', 'application/json; charset=utf-8')
    res.end(JSON.stringify(body, null, 2))
}

const fileStat = (filePath: string) => {
    try {
        return fs.statSync(filePath)
    } catch {
        return null
    }
}

const readPid = (pidFile: string) => {
    try {
        const raw = fs.readFileSync(pidFile, 'utf-8').trim()
        const pid = Number(raw)
        return Number.isFinite(pid) ? pid : null
    } catch {
        return null
    }
}

const isPidRunning = (pid: number | null) => {
    if (!pid) return false
    try {
        process.kill(pid, 0)
        return true
    } catch {
        return false
    }
}

const readLastLines = (filePath: string, lineCount: number) => {
    try {
        const content = fs.readFileSync(filePath, 'utf-8')
        return sanitizeLogContent(content.split(/\r?\n/).slice(-lineCount).join('\n'))
    } catch {
        return ''
    }
}

const sanitizeLogContent = (content: string) => {
    return content
        .replace(/([?&]token=)[^"'\s&]+/gi, '$1<redacted>')
        .replace(/(authorization:\s*bearer\s+)[^"'\s]+/gi, '$1<redacted>')
        .replace(/("token"\s*:\s*")[^"]+(")/gi, '$1<redacted>$2')
        .replace(/("password"\s*:\s*")[^"]+(")/gi, '$1<redacted>$2')
}

const buildServiceSummary = () => {
    return devServices.map((service) => {
        const logFile = path.join(devLogDir, `${service.name}.log`)
        const pidFile = path.join(devPidDir, `${service.name}.pid`)
        const pid = service.name === 'frontend' ? process.pid : readPid(pidFile)
        const stat = fileStat(logFile)

        return {
            ...service,
            pid,
            running: isPidRunning(pid),
            logFile,
            pidFile,
            logSize: stat?.size ?? 0,
            updatedAt: stat?.mtime?.toISOString() ?? null
        }
    })
}

const devLogsPlugin = () => ({
    name: 'local-dev-logs-api',
    configureServer(server: any) {
        server.middlewares.use('/__dev_logs/api/services', (_req: any, res: any) => {
            sendJson(res, 200, {
                generatedAt: new Date().toISOString(),
                logDir: devLogDir,
                pidDir: devPidDir,
                services: buildServiceSummary()
            })
        })

        server.middlewares.use('/__dev_logs/api/log', (req: any, res: any) => {
            const requestUrl = new URL(req.url || '', 'http://localhost')
            const name = requestUrl.searchParams.get('name') || ''
            const lines = Math.min(Number(requestUrl.searchParams.get('lines') || 300), 2000)
            const service = devServices.find((item) => item.name === name)

            if (!service) {
                sendJson(res, 404, { message: `unknown service: ${name}` })
                return
            }

            const logFile = path.join(devLogDir, `${service.name}.log`)
            sendJson(res, 200, {
                service: service.name,
                logFile,
                lines,
                content: readLastLines(logFile, lines)
            })
        })
    }
})

export default defineConfig({
    plugins: [vue(), devLogsPlugin()],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
    },
    build: {
        // 交给 Vite/Rollup 自动分包，避免手动拆分 Element Plus 与 Vue 造成循环初始化白屏。
        chunkSizeWarningLimit: 1200
    },
    server: {
        port: 5173,
        proxy: {
            // 用户服务
            '/api/v1/login': {
                target: 'http://localhost:8081',
                changeOrigin: true
            },
            '/api/v1/register': {
                target: 'http://localhost:8081',
                changeOrigin: true
            },
            '/api/v1/profile': {
                target: 'http://localhost:8081',
                changeOrigin: true
            },
            '/api/v1/users': {
                target: 'http://localhost:8081',
                changeOrigin: true
            },
            // 职位服务
            '/api/v1/jobs': {
                target: 'http://localhost:8082',
                changeOrigin: true
            },
            // 面试服务
            '/api/v1/interviews': {
                target: 'http://localhost:8083',
                changeOrigin: true
            },
            // 简历服务
            '/api/v1/resumes': {
                target: 'http://localhost:8084',
                changeOrigin: true,
                // 配置代理以支持大文件上传
                configure: (proxy, options) => {
                    proxy.on('proxyReq', (proxyReq, req, res) => {
                        // 不要修改 content-length
                    })
                }
            },
            '/api/v1/applications': {
                target: 'http://localhost:8084',
                changeOrigin: true
            },
            '/api/v1/ai': {
                target: 'http://localhost:8084',
                changeOrigin: true
            },
            // 评估结果服务
            '/api/v1/evaluations': {
                target: 'http://localhost:8084',
                changeOrigin: true
            },
            // 消息服务
            '/api/v1/messages': {
                target: 'http://localhost:8085',
                changeOrigin: true
            },
            // 公告服务 (message-service)
            '/api/v1/notices': {
                target: 'http://localhost:8085',
                changeOrigin: true
            },
            // 人才服务
            '/api/v1/talents': {
                target: 'http://localhost:8086',
                changeOrigin: true
            },
            // 推荐服务
            '/api/v1/recommendations': {
                target: 'http://localhost:8087',
                changeOrigin: true
            },
            // 统计服务 (gateway)
            '/api/v1/stats': {
                target: 'http://localhost:8080',
                changeOrigin: true
            },
            // 日志服务（统一走 Gateway，避免生产环境断链）
            '/api/v1/logs': {
                target: 'http://localhost:8080',
                changeOrigin: true
            },
            // 聊天会话服务 (message-service)
            '/api/v1/conversations': {
                target: 'http://localhost:8085',
                changeOrigin: true
            },
            // WebSocket 代理 (通过 Gateway)
            '/api/v1/ws': {
                target: 'ws://localhost:8080',
                ws: true,
                changeOrigin: true
            }
        }
    }
})
