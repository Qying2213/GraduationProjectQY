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
    },

    // 生成归因报告
    generateAttributionReport(talentId: number, jobId: number) {
        return request.post<ApiResponse>('/recommendations/attribution-report', {
            talent_id: talentId,
            job_id: jobId
        })
    },

    // 语义匹配
    semanticMatch(text1: string, text2: string) {
        return request.post<ApiResponse>('/recommendations/semantic-match', {
            text1,
            text2
        })
    },

    // 批量推荐
    batchRecommend(talentIds: number[], jobIds: number[]) {
        return request.post<ApiResponse>('/recommendations/batch', {
            talent_ids: talentIds,
            job_ids: jobIds
        })
    },

    // RAG 查询
    ragQuery(query: string, type: 'talent' | 'job' = 'talent', topK: number = 5) {
        return request.post<ApiResponse>('/recommendations/rag/query', {
            query,
            query_type: type,
            top_k: topK
        })
    },

    // 索引人才到向量数据库
    indexTalent(talentId: number) {
        return request.post<ApiResponse>('/recommendations/rag/index-talent', { talent_id: talentId })
    },

    // 索引职位到向量数据库
    indexJob(jobId: number) {
        return request.post<ApiResponse>('/recommendations/rag/index-job', { job_id: jobId })
    },

    // 批量索引所有数据
    indexAll() {
        return request.post<ApiResponse>('/recommendations/rag/index-all')
    },

    // RAG 匹配
    ragMatch(talentId: number, jobId: number) {
        return request.post<ApiResponse>('/recommendations/rag/match', {
            talent_id: talentId,
            job_id: jobId
        })
    }
}
