import request from '@/utils/request'
import type { Job, ApiResponse } from '@/types'

export const jobApi = {
    // 创建职位
    create(data: Partial<Job>) {
        return request.post<ApiResponse<Job>>('/jobs', data)
    },

    // 获取职位列表
    list(params?: any) {
        return request.get<ApiResponse>('/jobs', { params })
    },

    // 获取职位详情
    get(id: number) {
        return request.get<ApiResponse<Job>>(`/jobs/${id}`)
    },

    // 更新职位
    update(id: number, data: Partial<Job>) {
        return request.put<ApiResponse<Job>>(`/jobs/${id}`, data)
    },

    // 删除职位
    delete(id: number) {
        return request.delete<ApiResponse>(`/jobs/${id}`)
    },

    // 获取职位统计
    getStats() {
        return request.get<ApiResponse>('/jobs/stats')
    }
}
