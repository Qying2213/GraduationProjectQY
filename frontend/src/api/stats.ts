import request from '@/utils/request'

export interface DashboardStats {
  total_talents: number
  total_jobs: number
  total_applications: number
  total_interviews: number
  match_rate: number
  talent_trend: number
  job_trend: number
  application_trend: number
}

export interface RecruitmentFunnel {
  resumes: number
  screened: number
  interviewed: number
  passed: number
  hired: number
}

export interface ChannelStats {
  name: string
  count: number
  rate: number
}

export interface DepartmentProgress {
  department: string
  target: number
  hired: number
  progress: number
}

export interface InterviewerRank {
  name: string
  department: string
  interviews: number
  pass_rate: number
  avg_score: number
}

export interface TrendData {
  date: string
  resumes: number
  interviews: number
  hired: number
}

export const statsApi = {
  // 获取仪表板统计
  getDashboardStats() {
    return request.get<{ code: number; message: string; data: DashboardStats }>('/stats/dashboard')
  },

  // 获取招聘漏斗数据
  getRecruitmentFunnel(period?: string) {
    return request.get<{ code: number; message: string; data: RecruitmentFunnel }>('/stats/funnel', {
      params: { period }
    })
  },

  // 获取渠道统计
  getChannelStats() {
    return request.get<{ code: number; message: string; data: ChannelStats[] }>('/stats/channels')
  },

  // 获取部门招聘进度
  getDepartmentProgress() {
    return request.get<{ code: number; message: string; data: DepartmentProgress[] }>('/stats/department-progress')
  },

  // 获取面试官排行
  getInterviewerRank() {
    return request.get<{ code: number; message: string; data: InterviewerRank[] }>('/stats/interviewer-rank')
  },

  // 获取趋势数据
  getTrendData(startDate?: string, endDate?: string) {
    return request.get<{ code: number; message: string; data: TrendData[] }>('/stats/trend', {
      params: { start_date: startDate, end_date: endDate }
    })
  },

  // 获取职位热度排行
  getJobRank(limit?: number) {
    return request.get<{ code: number; message: string; data: { title: string; count: number }[] }>('/stats/job-rank', {
      params: { limit }
    })
  }
}
