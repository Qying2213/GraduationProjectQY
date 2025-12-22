import { describe, it, expect, vi, beforeEach } from 'vitest'
import { exportToExcel, exportToCsv, talentExportColumns, jobExportColumns } from '../export'

// Mock XLSX
vi.mock('xlsx', () => ({
  utils: {
    aoa_to_sheet: vi.fn(() => ({})),
    book_new: vi.fn(() => ({})),
    book_append_sheet: vi.fn(),
  },
  writeFile: vi.fn(),
}))

describe('Export Utils', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('exportToExcel', () => {
    it('should export data to Excel format', () => {
      const mockData = [
        { name: '张三', email: 'zhangsan@example.com', skills: ['Go', 'Vue'] },
        { name: '李四', email: 'lisi@example.com', skills: ['Java', 'React'] },
      ]

      const columns = [
        { key: 'name', title: '姓名' },
        { key: 'email', title: '邮箱' },
        { key: 'skills', title: '技能', formatter: (v: string[]) => v.join(', ') },
      ]

      expect(() => {
        exportToExcel({
          filename: 'test',
          columns,
          data: mockData,
        })
      }).not.toThrow()
    })

    it('should handle empty data', () => {
      expect(() => {
        exportToExcel({
          filename: 'empty',
          columns: [{ key: 'name', title: '姓名' }],
          data: [],
        })
      }).not.toThrow()
    })
  })

  describe('exportToCsv', () => {
    it('should create CSV content with proper formatting', () => {
      // Mock URL and document methods
      const mockCreateObjectURL = vi.fn(() => 'blob:test')
      const mockRevokeObjectURL = vi.fn()
      global.URL.createObjectURL = mockCreateObjectURL
      global.URL.revokeObjectURL = mockRevokeObjectURL

      const mockClick = vi.fn()
      const mockAppendChild = vi.fn()
      const mockRemoveChild = vi.fn()

      vi.spyOn(document, 'createElement').mockReturnValue({
        href: '',
        download: '',
        click: mockClick,
      } as any)
      vi.spyOn(document.body, 'appendChild').mockImplementation(mockAppendChild)
      vi.spyOn(document.body, 'removeChild').mockImplementation(mockRemoveChild)

      const mockData = [
        { name: '张三', phone: '13800138000' },
      ]

      exportToCsv({
        filename: 'test',
        columns: [
          { key: 'name', title: '姓名' },
          { key: 'phone', title: '电话' },
        ],
        data: mockData,
      })

      expect(mockCreateObjectURL).toHaveBeenCalled()
      expect(mockClick).toHaveBeenCalled()
    })

    it('should escape special characters in CSV', () => {
      const mockCreateObjectURL = vi.fn(() => 'blob:test')
      global.URL.createObjectURL = mockCreateObjectURL
      global.URL.revokeObjectURL = vi.fn()

      vi.spyOn(document, 'createElement').mockReturnValue({
        href: '',
        download: '',
        click: vi.fn(),
      } as any)
      vi.spyOn(document.body, 'appendChild').mockImplementation(vi.fn())
      vi.spyOn(document.body, 'removeChild').mockImplementation(vi.fn())

      const mockData = [
        { name: '张三,李四', description: '包含"引号"的内容' },
      ]

      expect(() => {
        exportToCsv({
          filename: 'test',
          columns: [
            { key: 'name', title: '姓名' },
            { key: 'description', title: '描述' },
          ],
          data: mockData,
        })
      }).not.toThrow()
    })
  })

  describe('Export Columns Configuration', () => {
    it('talentExportColumns should have required fields', () => {
      const requiredKeys = ['name', 'email', 'phone', 'skills', 'status']
      requiredKeys.forEach(key => {
        expect(talentExportColumns.some(col => col.key === key)).toBe(true)
      })
    })

    it('jobExportColumns should have required fields', () => {
      const requiredKeys = ['title', 'department', 'location', 'status']
      requiredKeys.forEach(key => {
        expect(jobExportColumns.some(col => col.key === key)).toBe(true)
      })
    })

    it('status formatter should return correct Chinese text', () => {
      const statusColumn = talentExportColumns.find(col => col.key === 'status')
      expect(statusColumn?.formatter).toBeDefined()
      expect(statusColumn?.formatter?.('active', {})).toBe('在职看机会')
    })

    it('skills formatter should join array', () => {
      const skillsColumn = talentExportColumns.find(col => col.key === 'skills')
      expect(skillsColumn?.formatter).toBeDefined()
      expect(skillsColumn?.formatter?.(['Go', 'Vue'], {})).toBe('Go, Vue')
    })
  })
})
