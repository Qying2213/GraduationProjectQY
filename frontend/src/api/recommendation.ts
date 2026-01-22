import request from '@/utils/request'
import type { Recommendation, ApiResponse } from '@/types'

export const recommendationApi = {
    // 为人才推荐职位
    recommendJobsForTalent(talent: any) {
        return request.post<ApiResponse<Recommendation[]>>('/api/v1/recommendations/jobs-for-talent', talent)
    },

    // 为职位推荐人才
    recommendTalentsForJob(job: any) {
        return request.post<ApiResponse<Recommendation[]>>('/api/v1/recommendations/talents-for-job', job)
    },

    // 获取推荐统计
    getStats() {
        return request.get<ApiResponse>('/api/v1/recommendations/stats')
    },

    // 生成归因报告
    generateAttributionReport(talentId: number, jobId: number) {
        return request.post<ApiResponse>('/api/v1/recommendations/attribution-report', {
            talent_id: talentId,
            job_id: jobId
        })
    },

    // 语义匹配
    semanticMatch(text1: string, text2: string) {
        return request.post<ApiResponse>('/api/v1/recommendations/semantic-match', {
            text1,
            text2
        })
    },

    // 批量推荐
    batchRecommend(talentIds: number[], jobIds: number[]) {
        return request.post<ApiResponse>('/api/v1/recommendations/batch', {
            talent_ids: talentIds,
            job_ids: jobIds
        })
    }
}
