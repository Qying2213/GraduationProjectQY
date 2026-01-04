import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

const instance: AxiosInstance = axios.create({
    baseURL: '/api/v1',
    timeout: 30000, // 30秒超时
    headers: {
        'Content-Type': 'application/json'
    }
})

// Request interceptor
instance.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token')
        if (token) {
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
            console.error('[Response Error] Status:', status, 'Data:', data)

            switch (status) {
                case 401:
                    ElMessage.error('未授权，请登录')
                    localStorage.removeItem('token')
                    localStorage.removeItem('user')
                    window.location.href = '/login'
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
