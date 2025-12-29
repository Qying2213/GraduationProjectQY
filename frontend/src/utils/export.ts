/**
 * 报表导出工具
 */

export interface ExportColumn {
  key: string
  title: string
  formatter?: (value: any, row: any) => string
}

export interface ExportOptions {
  filename: string
  sheetName?: string
  columns: ExportColumn[]
  data: any[]
}

// 下载Blob
function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

// 格式化单元格值
function formatCellValue(value: any, column: ExportColumn, row: any): string {
  if (column.formatter) {
    return column.formatter(value, row)
  }
  if (value === null || value === undefined) return ''
  if (Array.isArray(value)) return value.join(', ')
  return String(value)
}

// 导出为CSV
export function exportToCsv(options: ExportOptions) {
  const { filename, columns, data } = options
  if (!data.length) return

  const headerRow = columns.map(col => col.title)
  
  const csvContent = [
    headerRow.join(','),
    ...data.map(row => columns.map(col => {
      let value = formatCellValue(row[col.key], col, row)
      // 处理包含逗号、换行或引号的值
      if (value.includes(',') || value.includes('\n') || value.includes('"')) {
        value = `"${value.replace(/"/g, '""')}"`
      }
      return value
    }).join(','))
  ].join('\n')

  const BOM = '\uFEFF'
  const blob = new Blob([BOM + csvContent], { type: 'text/csv;charset=utf-8;' })
  downloadBlob(blob, `${filename}.csv`)
}

// 导出为Excel (使用CSV格式，兼容Excel打开)
export function exportToExcel(options: ExportOptions) {
  const { filename } = options
  exportToCsv({ ...options, filename: filename })
  // 实际导出为.xlsx需要引入xlsx库，这里简化为csv格式
}

// 导出为JSON
export function exportToJSON(data: any, filename: string) {
  const jsonContent = JSON.stringify(data, null, 2)
  const blob = new Blob([jsonContent], { type: 'application/json' })
  downloadBlob(blob, `${filename}.json`)
}

// 旧版兼容函数
export function exportToCSV(data: any[], filename: string, headers?: Record<string, string>) {
  if (!data.length) return
  const keys = Object.keys(headers || data[0])
  const columns: ExportColumn[] = keys.map(key => ({
    key,
    title: headers ? headers[key] : key
  }))
  exportToCsv({ filename, columns, data })
}

// 生成招聘报告数据
export function generateRecruitmentReport(data: {
  talents: any[]
  jobs: any[]
  interviews: any[]
  applications: any[]
}) {
  const now = new Date()
  const report = {
    title: '招聘数据报告',
    generatedAt: now.toLocaleString('zh-CN'),
    summary: {
      totalTalents: data.talents.length,
      totalJobs: data.jobs.length,
      totalInterviews: data.interviews.length,
      totalApplications: data.applications.length,
      openJobs: data.jobs.filter((j: any) => j.status === 'open').length,
      completedInterviews: data.interviews.filter((i: any) => i.status === 'completed').length
    },
    talents: data.talents,
    jobs: data.jobs,
    interviews: data.interviews
  }
  return report
}

// 人才列表导出列配置
export const talentExportColumns: ExportColumn[] = [
  { key: 'id', title: 'ID' },
  { key: 'name', title: '姓名' },
  { key: 'email', title: '邮箱' },
  { key: 'phone', title: '电话' },
  { key: 'skills', title: '技能', formatter: (v) => Array.isArray(v) ? v.join(', ') : v },
  { key: 'experience', title: '工作经验(年)' },
  { key: 'education', title: '学历' },
  { key: 'status', title: '状态', formatter: (v) => {
    const map: Record<string, string> = { active: '活跃', hired: '已雇佣', pending: '待处理', rejected: '已拒绝' }
    return map[v] || v
  }},
  { key: 'location', title: '所在地' },
  { key: 'salary', title: '期望薪资' }
]

// 职位列表导出列配置
export const jobExportColumns: ExportColumn[] = [
  { key: 'id', title: 'ID' },
  { key: 'title', title: '职位名称' },
  { key: 'department', title: '部门' },
  { key: 'location', title: '工作地点' },
  { key: 'salary', title: '薪资范围' },
  { key: 'type', title: '工作类型' },
  { key: 'status', title: '状态' }
]

// 面试列表导出列配置
export const interviewExportColumns: ExportColumn[] = [
  { key: 'id', title: 'ID' },
  { key: 'candidate_name', title: '候选人' },
  { key: 'position', title: '应聘职位' },
  { key: 'interviewer', title: '面试官' },
  { key: 'date', title: '面试日期' },
  { key: 'time', title: '面试时间' },
  { key: 'method', title: '面试方式' },
  { key: 'status', title: '状态' }
]

// 旧版兼容表头（向后兼容）
export const talentExportHeaders = {
  id: 'ID', name: '姓名', email: '邮箱', phone: '电话',
  experience: '工作经验(年)', education: '学历', status: '状态',
  location: '所在地', salary: '期望薪资'
}

export const jobExportHeaders = {
  id: 'ID', title: '职位名称', department: '部门', location: '工作地点',
  salary: '薪资范围', type: '工作类型', status: '状态'
}

export const interviewExportHeaders = {
  id: 'ID', candidate_name: '候选人', position: '应聘职位',
  interviewer: '面试官', date: '面试日期', time: '面试时间',
  method: '面试方式', status: '状态'
}
