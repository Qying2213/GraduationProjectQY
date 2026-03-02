import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
    plugins: [vue()],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
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
