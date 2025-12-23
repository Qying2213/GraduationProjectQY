import { describe, it, expect } from 'vitest'

// 格式化工具函数
const formatFileSize = (bytes: number): string => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

const formatSalary = (min: number, max: number): string => {
  return `${min}-${max}K`
}

const formatExperience = (years: number): string => {
  if (years === 0) return '应届生'
  if (years < 1) return '1年以下'
  if (years <= 3) return '1-3年'
  if (years <= 5) return '3-5年'
  if (years <= 10) return '5-10年'
  return '10年以上'
}

const maskPhone = (phone: string): string => {
  if (!phone || phone.length !== 11) return phone
  return phone.slice(0, 3) + '****' + phone.slice(7)
}

const maskEmail = (email: string): string => {
  if (!email || !email.includes('@')) return email
  const [name, domain] = email.split('@')
  if (name.length <= 2) return email
  return name.slice(0, 2) + '***@' + domain
}

describe('格式化工具函数', () => {
  describe('formatFileSize', () => {
    it('应该正确格式化字节', () => {
      expect(formatFileSize(500)).toBe('500 B')
    })

    it('应该正确格式化KB', () => {
      expect(formatFileSize(1024)).toBe('1.0 KB')
      expect(formatFileSize(2048)).toBe('2.0 KB')
    })

    it('应该正确格式化MB', () => {
      expect(formatFileSize(1024 * 1024)).toBe('1.0 MB')
      expect(formatFileSize(5 * 1024 * 1024)).toBe('5.0 MB')
    })
  })

  describe('formatDate', () => {
    it('应该正确格式化日期', () => {
      const result = formatDate('2024-01-15')
      expect(result).toMatch(/2024/)
      expect(result).toMatch(/01/)
      expect(result).toMatch(/15/)
    })
  })

  describe('formatSalary', () => {
    it('应该正确格式化薪资范围', () => {
      expect(formatSalary(20, 30)).toBe('20-30K')
      expect(formatSalary(15, 25)).toBe('15-25K')
    })
  })

  describe('formatExperience', () => {
    it('应该正确格式化经验年限', () => {
      expect(formatExperience(0)).toBe('应届生')
      expect(formatExperience(2)).toBe('1-3年')
      expect(formatExperience(4)).toBe('3-5年')
      expect(formatExperience(7)).toBe('5-10年')
      expect(formatExperience(12)).toBe('10年以上')
    })
  })

  describe('maskPhone', () => {
    it('应该正确脱敏手机号', () => {
      expect(maskPhone('13812345678')).toBe('138****5678')
    })

    it('无效手机号应该原样返回', () => {
      expect(maskPhone('123')).toBe('123')
      expect(maskPhone('')).toBe('')
    })
  })

  describe('maskEmail', () => {
    it('应该正确脱敏邮箱', () => {
      expect(maskEmail('zhangsan@example.com')).toBe('zh***@example.com')
    })

    it('短用户名应该原样返回', () => {
      expect(maskEmail('ab@example.com')).toBe('ab@example.com')
    })

    it('无效邮箱应该原样返回', () => {
      expect(maskEmail('invalid')).toBe('invalid')
    })
  })
})
