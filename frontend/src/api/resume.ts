import request from '@/utils/request'
import type { Resume, ApiResponse } from '@/types'

// 在线简历数据类型
export interface BasicInfo {
    name: string
    phone: string
    email: string
    location: string
    avatar?: string
    gender?: string
    age?: number
    summary?: string
}

export interface WorkExperience {
    company: string
    position: string
    start_date: string
    end_date?: string
    is_current: boolean
    description?: string
    location?: string
}

export interface Education {
    school: string
    degree: string
    major: string
    start_date: string
    end_date?: string
    is_current: boolean
    gpa?: string
    activities?: string
}

export interface OnlineResumeRequest {
    basic_info: BasicInfo
    work_experience: WorkExperience[]
    education: Education[]
    skills: string[]
}

export interface OnlineResumeResponse {
    id: number
    basic_info: BasicInfo
    work_experience: WorkExperience[]
    education: Education[]
    skills: string[]
    is_complete: boolean
    updated_at: string
}

export const resumeApi = {
    // 上传简历文件
    upload(file: File, talentId?: number, jobId?: number) {
        const formData = new FormData()
        formData.append('file', file)
        if (talentId) formData.append('talent_id', talentId.toString())
        if (jobId) formData.append('job_id', jobId.toString())
        return request.post<ApiResponse<Resume>>('/resumes/upload', formData)
    },

    // 创建简历记录（JSON方式）
    create(data: Partial<Resume>) {
        return request.post<ApiResponse<Resume>>('/resumes', data)
    },

    // 获取简历列表
    list(params?: { page?: number; page_size?: number; talent_id?: number; status?: string; search?: string }) {
        return request.get<ApiResponse>('/resumes', { params })
    },

    // 获取简历详情
    get(id: number) {
        return request.get<ApiResponse<Resume>>(`/resumes/${id}`)
    },

    // 下载简历
    download(id: number) {
        return request.get(`/resumes/${id}/download`, { responseType: 'blob' })
    },

    // 删除简历
    delete(id: number) {
        return request.delete<ApiResponse>(`/resumes/${id}`)
    },

    // 更新简历状态
    updateStatus(id: number, status: string) {
        return request.put<ApiResponse>(`/resumes/${id}/status`, { status })
    },

    // 解析简历文本
    parse(text: string) {
        return request.post<ApiResponse>('/resumes/parse', { text })
    },

    // 计算简历与职位匹配度
    match(resumeText: string, jobSkills: string[], jobExperience: number, jobEducation: string) {
        return request.post<ApiResponse>('/resumes/match', {
            resume_text: resumeText,
            job_skills: jobSkills,
            job_experience: jobExperience,
            job_education: jobEducation
        })
    },

    // 获取用于评估的简历列表
    listForEvaluation(params?: { page?: number; page_size?: number; status?: string }) {
        return request.get<ApiResponse>('/resumes/evaluation', { params })
    },

    // 获取当前用户的在线简历
    getOnlineResume() {
        return request.get<ApiResponse<OnlineResumeResponse>>('/resumes/online')
    },

    // 保存在线简历
    saveOnlineResume(data: OnlineResumeRequest) {
        return request.put<ApiResponse<OnlineResumeResponse>>('/resumes/online', data)
    }
}

// AI 评估相关 API
export const aiApi = {
    // 检查 AI 配置
    checkConfig() {
        return request.get<ApiResponse>('/ai/config')
    },

    // 获取当前评估任务
    getCurrentTask() {
        return request.get<ApiResponse>('/ai/current-task')
    },

    // 根据简历 ID 评估
    evaluate(resumeId: number, jobId?: number) {
        return request.post<ApiResponse>('/ai/evaluate', { resume_id: resumeId, job_id: jobId })
    },

    // 评估上传的文件
    evaluateUpload(file: File, jobId?: number) {
        const formData = new FormData()
        formData.append('file', file)
        if (jobId) formData.append('job_id', jobId.toString())
        return request.post<ApiResponse>('/ai/evaluate/upload', formData)
    },

    // 批量评估
    batchEvaluate(resumeIds: number[], jobId?: number) {
        return request.post<ApiResponse>('/ai/evaluate/batch', { resume_ids: resumeIds, job_id: jobId })
    },

    // 获取评估结果
    getResult(id: number) {
        return request.get<ApiResponse>(`/ai/evaluate/${id}/result`)
    },

    // AI 智能解析简历
    parseResume(file: File) {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<ApiResponse>('/ai/parse', formData)
    },

    // OCR 文本提取
    ocrExtract(file: File) {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<ApiResponse>('/ai/ocr', formData)
    }
}

// 评估结果 API
export const evaluationApi = {
    // 获取评估结果列表
    list(params?: { page?: number; page_size?: number; status?: string }) {
        return request.get<ApiResponse>('/evaluations', { params })
    },

    // 获取评估统计
    getStats() {
        return request.get<ApiResponse>('/evaluations/stats')
    },

    // 获取评估详情
    get(id: number) {
        return request.get<ApiResponse>(`/evaluations/${id}`)
    },

    // 删除评估结果
    delete(id: number) {
        return request.delete<ApiResponse>(`/evaluations/${id}`)
    }
}
