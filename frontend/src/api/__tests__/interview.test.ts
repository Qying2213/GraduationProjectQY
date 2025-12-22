import { describe, it, expect, vi, beforeEach } from 'vitest'
import { interviewApi } from '../interview'

// Mock axios request
vi.mock('../request', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

import request from '../request'

describe('Interview API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('create', () => {
    it('should call POST /interviews with correct data', async () => {
      const mockData = {
        candidate_id: 1,
        candidate_name: '张三',
        position_id: 1,
        position: 'Go开发工程师',
        type: 'initial',
        date: '2024-12-25',
        time: '14:00',
        duration: 60,
        interviewer_id: 1,
        interviewer: '李四',
        method: 'onsite',
        location: '会议室A',
      }

      vi.mocked(request.post).mockResolvedValue({ data: { code: 0, data: mockData } })

      await interviewApi.create(mockData)

      expect(request.post).toHaveBeenCalledWith('/interviews', mockData)
    })
  })

  describe('list', () => {
    it('should call GET /interviews with params', async () => {
      const params = { page: 1, page_size: 20, status: 'scheduled' }

      vi.mocked(request.get).mockResolvedValue({
        data: { code: 0, data: { interviews: [], total: 0 } },
      })

      await interviewApi.list(params)

      expect(request.get).toHaveBeenCalledWith('/interviews', { params })
    })

    it('should work without params', async () => {
      vi.mocked(request.get).mockResolvedValue({
        data: { code: 0, data: { interviews: [], total: 0 } },
      })

      await interviewApi.list()

      expect(request.get).toHaveBeenCalledWith('/interviews', { params: undefined })
    })
  })

  describe('get', () => {
    it('should call GET /interviews/:id', async () => {
      vi.mocked(request.get).mockResolvedValue({
        data: { code: 0, data: { id: 1 } },
      })

      await interviewApi.get(1)

      expect(request.get).toHaveBeenCalledWith('/interviews/1')
    })
  })

  describe('update', () => {
    it('should call PUT /interviews/:id with data', async () => {
      const updateData = { status: 'completed', rating: 4 }

      vi.mocked(request.put).mockResolvedValue({
        data: { code: 0, data: { id: 1, ...updateData } },
      })

      await interviewApi.update(1, updateData)

      expect(request.put).toHaveBeenCalledWith('/interviews/1', updateData)
    })
  })

  describe('delete', () => {
    it('should call DELETE /interviews/:id', async () => {
      vi.mocked(request.delete).mockResolvedValue({
        data: { code: 0, message: 'success' },
      })

      await interviewApi.delete(1)

      expect(request.delete).toHaveBeenCalledWith('/interviews/1')
    })
  })

  describe('cancel', () => {
    it('should call POST /interviews/:id/cancel', async () => {
      vi.mocked(request.post).mockResolvedValue({
        data: { code: 0, message: 'success' },
      })

      await interviewApi.cancel(1)

      expect(request.post).toHaveBeenCalledWith('/interviews/1/cancel')
    })
  })

  describe('submitFeedback', () => {
    it('should call POST /interviews/:id/feedback with data', async () => {
      const feedbackData = {
        rating: 4,
        strengths: '技术扎实',
        weaknesses: '经验略少',
        comments: '整体良好',
        recommendation: 'pass' as const,
      }

      vi.mocked(request.post).mockResolvedValue({
        data: { code: 0, data: feedbackData },
      })

      await interviewApi.submitFeedback(1, feedbackData)

      expect(request.post).toHaveBeenCalledWith('/interviews/1/feedback', feedbackData)
    })
  })

  describe('getFeedback', () => {
    it('should call GET /interviews/:id/feedback', async () => {
      vi.mocked(request.get).mockResolvedValue({
        data: { code: 0, data: [] },
      })

      await interviewApi.getFeedback(1)

      expect(request.get).toHaveBeenCalledWith('/interviews/1/feedback')
    })
  })

  describe('reschedule', () => {
    it('should call POST /interviews/:id/reschedule with data', async () => {
      const rescheduleData = {
        date: '2024-12-26',
        time: '15:00',
        reason: '面试官有事',
      }

      vi.mocked(request.post).mockResolvedValue({
        data: { code: 0, data: { id: 1 } },
      })

      await interviewApi.reschedule(1, rescheduleData)

      expect(request.post).toHaveBeenCalledWith('/interviews/1/reschedule', rescheduleData)
    })
  })

  describe('getStats', () => {
    it('should call GET /interviews/stats', async () => {
      vi.mocked(request.get).mockResolvedValue({
        data: {
          code: 0,
          data: {
            total_interviews: 100,
            scheduled_interviews: 20,
            completed_interviews: 70,
            cancelled_interviews: 10,
          },
        },
      })

      await interviewApi.getStats()

      expect(request.get).toHaveBeenCalledWith('/interviews/stats')
    })
  })

  describe('getToday', () => {
    it('should call GET /interviews/today', async () => {
      vi.mocked(request.get).mockResolvedValue({
        data: { code: 0, data: [] },
      })

      await interviewApi.getToday()

      expect(request.get).toHaveBeenCalledWith('/interviews/today')
    })
  })

  describe('getInterviewerSchedule', () => {
    it('should call GET /interviews/interviewer/:id with date params', async () => {
      vi.mocked(request.get).mockResolvedValue({
        data: { code: 0, data: [] },
      })

      await interviewApi.getInterviewerSchedule(1, '2024-12-20', '2024-12-27')

      expect(request.get).toHaveBeenCalledWith('/interviews/interviewer/1', {
        params: { start_date: '2024-12-20', end_date: '2024-12-27' },
      })
    })
  })

  describe('getCandidateInterviews', () => {
    it('should call GET /interviews/candidate/:id', async () => {
      vi.mocked(request.get).mockResolvedValue({
        data: { code: 0, data: [] },
      })

      await interviewApi.getCandidateInterviews(1)

      expect(request.get).toHaveBeenCalledWith('/interviews/candidate/1')
    })
  })
})
