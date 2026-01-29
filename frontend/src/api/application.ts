import request from '@/utils/request'
import type { Application, ApiResponse } from '@/types'

export const applicationApi = {
    // 创建申请
    create(data: {
        job_id: number
        talent_id: number
        resume_id?: number
        cover_letter?: string
    }) {
        return request.post<ApiResponse<Application>>('/applications', data)
    },

    // 获取申请列表
    list(params?: {
        page?: number
        page_size?: number
        job_id?: number
        talent_id?: number
        status?: string
    }) {
        return request.get<ApiResponse>('/applications', { params })
    },

    // 获取申请详情
    get(id: number) {
        return request.get<ApiResponse<Application>>(`/applications/${id}`)
    },

    // 更新申请状态
    update(id: number, data: { status?: string; notes?: string }) {
        return request.put<ApiResponse<Application>>(`/applications/${id}`, data)
    },

    // 删除申请
    delete(id: number) {
        return request.delete<ApiResponse>(`/applications/${id}`)
    },

    // 获取我的申请（候选人视角）
    getMyApplications(params?: { page?: number; page_size?: number; status?: string }) {
        return request.get<ApiResponse>('/applications', {
            params: { ...params, talent_id: 'me' }
        })
    },

    // 获取职位的申请（HR视角）
    getJobApplications(jobId: number, params?: { page?: number; page_size?: number; status?: string }) {
        return request.get<ApiResponse>(`/jobs/${jobId}/applications`, { params })
    }
}
