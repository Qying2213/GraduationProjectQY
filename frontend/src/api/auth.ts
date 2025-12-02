import request from '@/utils/request'
import type { LoginRequest, RegisterRequest, User, ApiResponse } from '@/types'

export const authApi = {
    // 登录
    login(data: LoginRequest) {
        return request.post<ApiResponse<{ token: string; user: User }>>('/login', data)
    },

    // 注册
    register(data: RegisterRequest) {
        return request.post<ApiResponse<User>>('/register', data)
    },

    // 获取当前用户信息
    getProfile() {
        return request.get<ApiResponse<User>>('/profile')
    },

    // 更新用户信息
    updateProfile(data: Partial<User>) {
        return request.put<ApiResponse<User>>('/profile', data)
    },

    // 获取用户列表
    listUsers(params?: any) {
        return request.get<ApiResponse>('/users', { params })
    }
}
