/**
 * 智能人才运营平台 - 前端API集成测试
 * 测试所有API接口的调用
 */

import { describe, it, expect, beforeAll } from 'vitest'

// 各服务地址
const SERVICES = {
  user: 'http://localhost:8081/api/v1',
  job: 'http://localhost:8082/api/v1',
  interview: 'http://localhost:8083/api/v1',
  resume: 'http://localhost:8084/api/v1',
  message: 'http://localhost:8085/api/v1',
  talent: 'http://localhost:8086/api/v1',
}

// 辅助函数
async function fetchService(service: keyof typeof SERVICES, endpoint: string, options?: RequestInit) {
  const response = await fetch(`${SERVICES[service]}${endpoint}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    ...options,
  })
  return {
    status: response.status,
    data: await response.json().catch(() => null),
  }
}

describe('用户服务 (user-service)', () => {
  describe('认证接口', () => {
    it('POST /login - 正确密码登录成功', async () => {
      const res = await fetchService('user', '/login', {
        method: 'POST',
        body: JSON.stringify({
          username: 'admin',
          password: 'password123',
        }),
      })
      expect(res.status).toBe(200)
      expect(res.data).toHaveProperty('data')
    })

    it('POST /login - 错误密码登录失败', async () => {
      const res = await fetchService('user', '/login', {
        method: 'POST',
        body: JSON.stringify({
          username: 'admin',
          password: 'wrongpassword',
        }),
      })
      expect(res.status).toBe(401)
    })
  })

  describe('用户接口', () => {
    it('GET /users - 获取用户列表', async () => {
      const res = await fetchService('user', '/users')
      expect(res.status).toBe(200)
      expect(res.data).toHaveProperty('data')
    })

    it('GET /users - 分页获取', async () => {
      const res = await fetchService('user', '/users?page=1&page_size=5')
      expect(res.status).toBe(200)
    })
  })
})

describe('职位服务 (job-service)', () => {
  describe('职位列表', () => {
    it('GET /jobs - 获取职位列表', async () => {
      const res = await fetchService('job', '/jobs')
      expect(res.status).toBe(200)
      expect(res.data).toHaveProperty('data')
    })

    it('GET /jobs?page=1&page_size=5 - 分页获取', async () => {
      const res = await fetchService('job', '/jobs?page=1&page_size=5')
      expect(res.status).toBe(200)
    })

    it('GET /jobs?status=open - 按状态筛选', async () => {
      const res = await fetchService('job', '/jobs?status=open')
      expect(res.status).toBe(200)
    })

    it('GET /jobs?location=北京 - 按地点筛选', async () => {
      const res = await fetchService('job', '/jobs?location=北京')
      expect(res.status).toBe(200)
    })
  })

  describe('职位详情', () => {
    it('GET /jobs/1 - 获取职位详情', async () => {
      const res = await fetchService('job', '/jobs/1')
      expect(res.status).toBe(200)
      expect(res.data.data).toHaveProperty('id', 1)
    })

    it('GET /jobs/9999 - 获取不存在的职位', async () => {
      const res = await fetchService('job', '/jobs/9999')
      expect(res.status).toBe(404)
    })
  })

  describe('职位统计', () => {
    it('GET /jobs/stats - 获取职位统计', async () => {
      const res = await fetchService('job', '/jobs/stats')
      expect(res.status).toBe(200)
    })
  })

  describe('创建职位', () => {
    it('POST /jobs - 创建新职位', async () => {
      const res = await fetchService('job', '/jobs', {
        method: 'POST',
        body: JSON.stringify({
          title: '测试职位_' + Date.now(),
          description: '这是一个测试职位',
          location: '北京',
          salary: '20-30K',
          type: 'full-time',
          status: 'open',
          department: '技术部',
          created_by: 1,
        }),
      })
      expect(res.status).toBe(201)
    })
  })
})

describe('面试服务 (interview-service)', () => {
  describe('面试列表', () => {
    it('GET /interviews - 获取面试列表', async () => {
      const res = await fetchService('interview', '/interviews')
      expect(res.status).toBe(200)
      expect(res.data).toHaveProperty('data')
    })

    it('GET /interviews?status=scheduled - 按状态筛选', async () => {
      const res = await fetchService('interview', '/interviews?status=scheduled')
      expect(res.status).toBe(200)
    })
  })

  describe('面试详情', () => {
    it('GET /interviews/1 - 获取面试详情', async () => {
      const res = await fetchService('interview', '/interviews/1')
      expect(res.status).toBe(200)
      expect(res.data.data).toHaveProperty('ID', 1)
    })
  })

  describe('面试统计', () => {
    it('GET /interviews/stats - 获取面试统计', async () => {
      const res = await fetchService('interview', '/interviews/stats')
      expect(res.status).toBe(200)
    })

    it('GET /interviews/today - 获取今日面试', async () => {
      const res = await fetchService('interview', '/interviews/today')
      expect(res.status).toBe(200)
    })
  })
})

describe('简历服务 (resume-service)', () => {
  describe('简历列表', () => {
    it('GET /resumes - 获取简历列表', async () => {
      const res = await fetchService('resume', '/resumes')
      expect(res.status).toBe(200)
      expect(res.data).toHaveProperty('data')
    })

    it('GET /resumes?page=1&page_size=5 - 分页获取', async () => {
      const res = await fetchService('resume', '/resumes?page=1&page_size=5')
      expect(res.status).toBe(200)
    })

    it('GET /resumes?status=pending - 按状态筛选', async () => {
      const res = await fetchService('resume', '/resumes?status=pending')
      expect(res.status).toBe(200)
    })

    it('GET /resumes?sort_by=created_at&sort_order=desc - 排序', async () => {
      const res = await fetchService('resume', '/resumes?sort_by=created_at&sort_order=desc')
      expect(res.status).toBe(200)
    })
  })

  describe('简历详情', () => {
    it('GET /resumes/1 - 获取简历详情', async () => {
      const res = await fetchService('resume', '/resumes/1')
      expect(res.status).toBe(200)
    })
  })

  describe('AI配置', () => {
    it('GET /ai/config - 获取AI配置状态', async () => {
      const res = await fetchService('resume', '/ai/config')
      expect(res.status).toBe(200)
    })
  })
})

describe('人才服务 (talent-service)', () => {
  describe('人才列表', () => {
    it('GET /talents - 获取人才列表', async () => {
      const res = await fetchService('talent', '/talents')
      expect(res.status).toBe(200)
      expect(res.data).toHaveProperty('data')
    })

    it('GET /talents?page=1&page_size=5 - 分页获取', async () => {
      const res = await fetchService('talent', '/talents?page=1&page_size=5')
      expect(res.status).toBe(200)
    })
  })

  describe('人才详情', () => {
    it('GET /talents/1 - 获取人才详情', async () => {
      const res = await fetchService('talent', '/talents/1')
      expect(res.status).toBe(200)
    })
  })

  describe('人才搜索', () => {
    it('GET /talents/search - 搜索人才', async () => {
      const res = await fetchService('talent', '/talents/search?keyword=前端')
      expect(res.status).toBe(200)
    })
  })
})

describe('消息服务 (message-service)', () => {
  describe('消息列表', () => {
    it('GET /messages?user_id=1 - 获取消息列表', async () => {
      const res = await fetchService('message', '/messages?user_id=1')
      expect(res.status).toBe(200)
    })
  })

  describe('未读消息', () => {
    it('GET /messages/unread-count?user_id=1 - 获取未读消息数', async () => {
      const res = await fetchService('message', '/messages/unread-count?user_id=1')
      expect(res.status).toBe(200)
    })
  })

  describe('发送消息', () => {
    it('POST /messages - 发送新消息', async () => {
      const res = await fetchService('message', '/messages', {
        method: 'POST',
        body: JSON.stringify({
          sender_id: 1,
          receiver_id: 2,
          type: 'chat',
          title: '测试消息',
          content: '这是一条测试消息',
        }),
      })
      expect(res.status).toBe(201)
    })
  })
})
