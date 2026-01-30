import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

// 解析 JWT token 获取过期时间
function isTokenExpired(token: string): boolean {
    try {
        const payload = token.split('.')[1]
        const decoded = JSON.parse(atob(payload))
        if (decoded.exp) {
            // 提前 60 秒认为过期，避免边界情况
            return decoded.exp * 1000 < Date.now() + 60000
        }
        return false
    } catch {
        return true
    }
}

const instance: AxiosInstance = axios.create({
    baseURL: '/api/v1',
    timeout: 300000, // 5分钟超时，AI评估需要较长时间
    headers: {
        'Content-Type': 'application/json'
    }
})

// Request interceptor
instance.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token')
        if (token) {
            // 检查 token 是否过期
            if (isTokenExpired(token)) {
                localStorage.removeItem('token')
                localStorage.removeItem('user')
                ElMessage.error('登录已过期，请重新登录')
                // 根据当前路径判断跳转到哪个登录页
                const isPortal = window.location.pathname.startsWith('/portal')
                window.location.href = isPortal ? '/portal/login' : '/login'
                return Promise.reject(new Error('Token expired'))
            }
            config.headers.Authorization = `Bearer ${token}`
        }

        // 如果是 FormData，删除默认的 Content-Type，让浏览器自动设置
        if (config.data instanceof FormData) {
            console.log('[Request] 检测到 FormData，删除默认 Content-Type')
            delete config.headers['Content-Type']
        }

        console.log('[Request]', config.method?.toUpperCase(), config.url)
        console.log('[Request] Headers:', JSON.stringify(config.headers))

        return config
    },
    (error) => {
        console.error('[Request] 拦截器错误:', error)
        return Promise.reject(error)
    }
)

// Response interceptor
instance.interceptors.response.use(
    (response: AxiosResponse) => {
        console.log('[Response]', response.status, response.config.url)
        return response
    },
    (error) => {
        console.error('[Response Error]', error.config?.url, error.message)
        
        if (error.response) {
            const { status, data } = error.response
            const url = error.config?.url || ''
            console.error('[Response Error] Status:', status, 'Data:', data)

            switch (status) {
                case 401:
                    // 登录接口返回 401 不跳转，让登录页面自己处理错误
                    if (url.includes('/login')) {
                        // 登录失败，不跳转，返回错误让页面处理
                        return Promise.reject(error)
                    }
                    ElMessage.error('未授权，请登录')
                    localStorage.removeItem('token')
                    localStorage.removeItem('user')
                    // 根据当前路径判断跳转到哪个登录页
                    const currentPath = window.location.pathname
                    // 如果已经在登录页，不要跳转
                    if (currentPath === '/login' || currentPath === '/portal/login') {
                        break
                    }
                    const isPortal = currentPath.startsWith('/portal')
                    window.location.href = isPortal ? '/portal/login' : '/login'
                    break
                case 403:
                    ElMessage.error('无权限访问')
                    break
                case 404:
                    ElMessage.error('请求的资源不存在')
                    break
                case 500:
                    ElMessage.error('服务器错误')
                    break
                default:
                    ElMessage.error(data?.message || '请求失败')
            }
        } else {
            console.error('[Response Error] No response:', error)
            ElMessage.error('网络错误，请检查网络连接')
        }

        return Promise.reject(error)
    }
)

export default instance
