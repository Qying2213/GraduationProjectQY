import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { authApi } from '@/api/auth'

// 安全的 localStorage 操作。
// 浏览器隐私模式或存储异常时不让应用直接崩溃，而是降级为未持久化登录态。
const safeStorage = {
    setItem(key: string, value: string) {
        try {
            localStorage.setItem(key, value)
        } catch (e) {
            console.error('Failed to save to localStorage:', e)
        }
    },
    getItem(key: string): string | null {
        try {
            return localStorage.getItem(key)
        } catch (e) {
            console.error('Failed to read from localStorage:', e)
            return null
        }
    },
    removeItem(key: string) {
        try {
            localStorage.removeItem(key)
        } catch (e) {
            console.error('Failed to remove from localStorage:', e)
        }
    }
}

export const useUserStore = defineStore('user', () => {
    // user/token 是前端权限判断的核心状态。
    // token 用于请求拦截器追加 Authorization，user.role 用于路由和菜单分流。
    const user = ref<User | null>(null)
    const token = ref<string>('')

    // 初始化时从 localStorage 恢复登录态，让刷新页面后仍保持登录。
    const initFromStorage = () => {
        const storedToken = safeStorage.getItem('token')
        const storedUser = safeStorage.getItem('user')

        if (storedToken) {
            token.value = storedToken
        }

        if (storedUser) {
            try {
                user.value = JSON.parse(storedUser)
            } catch (e) {
                console.error('Failed to parse user data:', e)
                safeStorage.removeItem('user')
            }
        }
    }

    // 是否已登录
    const isLoggedIn = computed(() => !!token.value && !!user.value)

    // 用户角色
    const role = computed(() => user.value?.role || '')

    // 是否是管理员
    const isAdmin = computed(() => user.value?.role === 'admin')

    // 是否是HR
    const isHR = computed(() => ['hr', 'hr_manager', 'recruiter'].includes(user.value?.role || ''))

    // 是否是候选人
    const isCandidate = computed(() => user.value?.role === 'candidate')

    // 登录：调用后端 user-service，成功后同时更新内存状态和本地缓存。
    const login = async (username: string, password: string) => {
        console.log('[Login] 开始登录, 用户名:', username)
        try {
            const res = await authApi.login({ username, password })
            console.log('[Login] 响应:', res.data)
            
            if (res.data.code === 0 && res.data.data) {
                console.log('[Login] 登录成功, 用户:', res.data.data.user)
                token.value = res.data.data.token
                user.value = res.data.data.user

                safeStorage.setItem('token', res.data.data.token)
                safeStorage.setItem('user', JSON.stringify(res.data.data.user))
                return res.data
            } else {
                // 登录失败，抛出错误
                console.error('[Login] 登录失败:', res.data.message || res.data.error)
                throw new Error(res.data.message || res.data.error || '用户名或密码错误')
            }
        } catch (error: any) {
            console.error('[Login] 请求异常:', error)
            // 如果是 axios 错误，提取错误信息
            if (error.response?.data) {
                throw new Error(error.response.data.message || error.response.data.error || '登录失败')
            }
            throw error
        }
    }

    // 注册：供后台注册页和求职端注册流程复用。
    const register = async (data: any) => {
        const res = await authApi.register(data)
        return res.data
    }

    // 登出：清理内存和 localStorage，后续请求拦截器就不会再带旧 token。
    const logout = () => {
        user.value = null
        token.value = ''
        safeStorage.removeItem('token')
        safeStorage.removeItem('user')
    }

    // 更新用户信息：资料保存成功后同步刷新本地缓存，避免页面显示旧资料。
    const updateProfile = async (data: Partial<User>) => {
        const res = await authApi.updateProfile(data)
        if (res.data.code === 0 && res.data.data) {
            user.value = res.data.data
            safeStorage.setItem('user', JSON.stringify(res.data.data))
        }
        return res.data
    }

    // 获取用户信息
    const fetchProfile = async () => {
        const res = await authApi.getProfile()
        if (res.data.code === 0 && res.data.data) {
            user.value = res.data.data
            safeStorage.setItem('user', JSON.stringify(res.data.data))
        }
        return res.data
    }

    // 初始化
    initFromStorage()

    return {
        user,
        token,
        isLoggedIn,
        role,
        isAdmin,
        isHR,
        isCandidate,
        login,
        register,
        logout,
        updateProfile,
        fetchProfile
    }
})
