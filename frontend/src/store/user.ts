import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { authApi } from '@/api/auth'

// 安全的 localStorage 操作
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
    const user = ref<User | null>(null)
    const token = ref<string>('')

    // 初始化时从localStorage读取
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
    const isHR = computed(() => user.value?.role === 'hr')

    // 是否是候选人
    const isCandidate = computed(() => user.value?.role === 'candidate')

    // 登录
    const login = async (username: string, password: string) => {
        console.log('[Login] 开始登录, 用户名:', username)
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
            console.error('[Login] 登录失败:', res.data.message)
            throw new Error(res.data.message || '用户名或密码错误')
        }
    }

    // 注册
    const register = async (data: any) => {
        const res = await authApi.register(data)
        return res.data
    }

    // 登出
    const logout = () => {
        user.value = null
        token.value = ''
        safeStorage.removeItem('token')
        safeStorage.removeItem('user')
    }

    // 更新用户信息
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
