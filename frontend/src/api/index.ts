// API 模块统一导出
export { authApi } from './auth'
export { talentApi } from './talent'
export { jobApi } from './job'
export { interviewApi } from './interview'
export { messageApi } from './message'
export { statsApi } from './stats'
export { recommendationApi } from './recommendation'
export { resumeApi, aiApi, evaluationApi } from './resume'
export { applicationApi } from './application'

// 类型导出
export type {
    Interview,
    InterviewFeedback,
    InterviewListResponse,
    InterviewStats,
    CreateInterviewRequest,
    UpdateInterviewRequest,
    SubmitFeedbackRequest,
    RescheduleRequest,
    InterviewListParams
} from './interview'
