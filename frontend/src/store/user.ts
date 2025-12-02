import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { authApi } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
    const user = ref<User | null>(null)
    const token = ref<string>('')

    // 初始化时从localStorage读取
    const initFromStorage = () => {
        const storedToken = localStorage.getItem('token')
        const storedUser = localStorage.getItem('user')

        if (storedToken) {
            token.value = storedToken
        }

        if (storedUser) {
            try {
                user.value = JSON.parse(storedUser)
            } catch (e) {
                console.error('Failed to parse user data:', e)
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
        const res = await authApi.login({ username, password })
        if (res.data.code === 0 && res.data.data) {
            token.value = res.data.data.token
            user.value = res.data.data.user

            localStorage.setItem('token', res.data.data.token)
            localStorage.setItem('user', JSON.stringify(res.data.data.user))
        }
        return res.data
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
        localStorage.removeItem('token')
        localStorage.removeItem('user')
    }

    // 更新用户信息
    const updateProfile = async (data: Partial<User>) => {
        const res = await authApi.updateProfile(data)
        if (res.data.code === 0 && res.data.data) {
            user.value = res.data.data
            localStorage.setItem('user', JSON.stringify(res.data.data))
        }
        return res.data
    }

    // 获取用户信息
    const fetchProfile = async () => {
        const res = await authApi.getProfile()
        if (res.data.code === 0 && res.data.data) {
            user.value = res.data.data
            localStorage.setItem('user', JSON.stringify(res.data.data))
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
