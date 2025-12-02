mport request from '@/utils/request'
import type { Talent, ApiResponse } from '@/types'

export const talentApi = {
    // 创建人才
    create(data: Partial<Talent>) {
        return request.post<ApiResponse<Talent>>('/talents', data)
    },

    // 获取人才列表
    list(params?: any) {
        return request.get<ApiResponse>('/talents', { params })
    },

    // 获取人才详情
    get(id: number) {
        return request.get<ApiResponse<Talent>>(`/talents/${id}`)
    },

    // 更新人才
    update(id: number, data: Partial<Talent>) {
        return request.put<ApiResponse<Talent>>(`/talents/${id}`, data)
    },

    // 删除人才
    delete(id: number) {
        return request.delete<ApiResponse>(`/talents/${id}`)
    },

    // 搜索人才
    search(params: any) {
        return request.get<ApiResponse>('/talents/search', { params })
    }
}
