import request from '@/utils/request'
import type { Recommendation, ApiResponse } from '@/types'

export const recommendationApi = {
    // 为人才推荐职位
    recommendJobsForTalent(talent: any) {
        return request.post<ApiResponse<Recommendation[]>>('/recommendations/jobs-for-talent', talent)
    },

    // 为职位推荐人才
    recommendTalentsForJob(job: any) {
        return request.post<ApiResponse<Recommendation[]>>('/recommendations/talents-for-job', job)
    },

    // 获取推荐统计
    getStats() {
        return request.get<ApiResponse>('/recommendations/stats')
    }
}
