import request from '@/utils/request'

export interface Interview {
  id: number
  candidate_id: number
  candidate_name: string
  position_id: number
  position: string
  type: 'initial' | 'second' | 'final' | 'hr'
  date: string
  time: string
  duration: number
  interviewer_id: number
  interviewer: string
  method: 'onsite' | 'video' | 'phone'
  location: string
  status: 'scheduled' | 'completed' | 'cancelled' | 'no_show'
  notes: string
  feedback?: string
  rating?: number
  created_at: string
  updated_at: string
}

export interface InterviewFeedback {
  id: number
  interview_id: number
  interviewer_id: number
  rating: number
  strengths: string
  weaknesses: string
  comments: string
  recommendation: 'pass' | 'fail' | 'pending'
  created_at: string
}

export interface InterviewListResponse {
  interviews: Interview[]
  total: number
  page: number
  page_size: number
}

export interface InterviewStats {
  total_interviews: number
  scheduled_interviews: number
  completed_interviews: number
  cancelled_interviews: number
  today_interviews: number
  week_interviews: number
}

export interface CreateInterviewRequest {
  candidate_id: number
  candidate_name: string
  position_id: number
  position: string
  type: string
  date: string
  time: string
  duration: number
  interviewer_id: number
  interviewer: string
  method: string
  location: string
  notes?: string
}

export interface UpdateInterviewRequest {
  type?: string
  date?: string
  time?: string
  duration?: number
  interviewer_id?: number
  interviewer?: string
  method?: string
  location?: string
  status?: string
  notes?: string
  feedback?: string
  rating?: number
}

export interface SubmitFeedbackRequest {
  rating: number
  strengths?: string
  weaknesses?: string
  comments?: string
  recommendation: 'pass' | 'fail' | 'pending'
}

export interface RescheduleRequest {
  date: string
  time: string
  reason?: string
}

export interface InterviewListParams {
  page?: number
  page_size?: number
  status?: string
  date?: string
  start_date?: string
  end_date?: string
  interviewer_id?: number
  candidate_id?: number
}

export const interviewApi = {
  // 创建面试
  create(data: CreateInterviewRequest) {
    return request.post<{ code: number; message: string; data: Interview }>('/interviews', data)
  },

  // 获取面试列表
  list(params?: InterviewListParams) {
    return request.get<{ code: number; message: string; data: InterviewListResponse }>('/interviews', { params })
  },

  // 获取单个面试
  get(id: number) {
    return request.get<{ code: number; message: string; data: Interview }>(`/interviews/${id}`)
  },

  // 更新面试
  update(id: number, data: UpdateInterviewRequest) {
    return request.put<{ code: number; message: string; data: Interview }>(`/interviews/${id}`, data)
  },

  // 删除面试
  delete(id: number) {
    return request.delete<{ code: number; message: string }>(`/interviews/${id}`)
  },

  // 取消面试
  cancel(id: number) {
    return request.post<{ code: number; message: string }>(`/interviews/${id}/cancel`)
  },

  // 完成面试
  complete(id: number, feedback?: { feedback?: string; rating?: number }) {
    return request.post<{ code: number; message: string }>(`/interviews/${id}/complete`, feedback)
  },

  // 提交面试反馈
  submitFeedback(id: number, data: SubmitFeedbackRequest) {
    return request.post<{ code: number; message: string; data: InterviewFeedback }>(`/interviews/${id}/feedback`, data)
  },

  // 获取面试反馈
  getFeedback(id: number) {
    return request.get<{ code: number; message: string; data: InterviewFeedback[] }>(`/interviews/${id}/feedback`)
  },

  // 重新安排面试
  reschedule(id: number, data: RescheduleRequest) {
    return request.post<{ code: number; message: string; data: Interview }>(`/interviews/${id}/reschedule`, data)
  },

  // 获取候选人的所有面试
  getCandidateInterviews(candidateId: number) {
    return request.get<{ code: number; message: string; data: Interview[] }>(`/interviews/candidate/${candidateId}`)
  },

  // 获取面试统计
  getStats() {
    return request.get<{ code: number; message: string; data: InterviewStats }>('/interviews/stats')
  },

  // 获取今日面试
  getToday() {
    return request.get<{ code: number; message: string; data: Interview[] }>('/interviews/today')
  },

  // 获取面试官日程
  getInterviewerSchedule(interviewerId: number, startDate?: string, endDate?: string) {
    return request.get<{ code: number; message: string; data: Interview[] }>(
      `/interviews/interviewer/${interviewerId}`,
      { params: { start_date: startDate, end_date: endDate } }
    )
  }
}
