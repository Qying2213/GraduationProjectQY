import * as XLSX from 'xlsx'

export interface ExportColumn {
  key: string
  title: string
  width?: number
  formatter?: (value: any, row: any) => string
}

export interface ExportOptions {
  filename: string
  sheetName?: string
  columns: ExportColumn[]
  data: any[]
}

/**
 * Export data to Excel file
 */
export function exportToExcel(options: ExportOptions): void {
  const { filename, sheetName = 'Sheet1', columns, data } = options

  // Prepare header row
  const headers = columns.map(col => col.title)

  // Prepare data rows
  const rows = data.map(item => {
    return columns.map(col => {
      const value = item[col.key]
      if (col.formatter) {
        return col.formatter(value, item)
      }
      return value ?? ''
    })
  })

  // Create worksheet
  const wsData = [headers, ...rows]
  const ws = XLSX.utils.aoa_to_sheet(wsData)

  // Set column widths
  const colWidths = columns.map(col => ({
    wch: col.width || Math.max(
      col.title.length * 2,
      ...rows.map(row => String(row[columns.indexOf(col)] || '').length)
    )
  }))
  ws['!cols'] = colWidths

  // Create workbook
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, sheetName)

  // Generate filename with timestamp
  const timestamp = new Date().toISOString().slice(0, 10)
  const fullFilename = `${filename}_${timestamp}.xlsx`

  // Download file
  XLSX.writeFile(wb, fullFilename)
}

/**
 * Export data to CSV file
 */
export function exportToCsv(options: ExportOptions): void {
  const { filename, columns, data } = options

  // Prepare header row
  const headers = columns.map(col => col.title)

  // Prepare data rows
  const rows = data.map(item => {
    return columns.map(col => {
      const value = item[col.key]
      if (col.formatter) {
        return col.formatter(value, item)
      }
      // Escape quotes and wrap in quotes if contains comma
      const strValue = String(value ?? '')
      if (strValue.includes(',') || strValue.includes('"') || strValue.includes('\n')) {
        return `"${strValue.replace(/"/g, '""')}"`
      }
      return strValue
    })
  })

  // Combine into CSV string
  const csvContent = [
    headers.join(','),
    ...rows.map(row => row.join(','))
  ].join('\n')

  // Add BOM for Excel to recognize UTF-8
  const BOM = '\uFEFF'
  const blob = new Blob([BOM + csvContent], { type: 'text/csv;charset=utf-8;' })

  // Generate filename with timestamp
  const timestamp = new Date().toISOString().slice(0, 10)
  const fullFilename = `${filename}_${timestamp}.csv`

  // Download file
  downloadBlob(blob, fullFilename)
}

/**
 * Helper function to download blob
 */
function downloadBlob(blob: Blob, filename: string): void {
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

/**
 * Predefined export configurations for different data types
 */
export const talentExportColumns: ExportColumn[] = [
  { key: 'name', title: '姓名', width: 12 },
  { key: 'gender', title: '性别', width: 8 },
  { key: 'age', title: '年龄', width: 8 },
  { key: 'phone', title: '电话', width: 15 },
  { key: 'email', title: '邮箱', width: 25 },
  { key: 'education', title: '学历', width: 10 },
  { key: 'experience', title: '工作经验', width: 12 },
  { key: 'currentCompany', title: '当前公司', width: 20 },
  { key: 'currentPosition', title: '当前职位', width: 15 },
  { key: 'expectedSalary', title: '期望薪资', width: 12 },
  { key: 'skills', title: '技能', width: 30, formatter: (v) => Array.isArray(v) ? v.join(', ') : v },
  { key: 'status', title: '状态', width: 10, formatter: (v) => {
    const statusMap: Record<string, string> = {
      active: '在职看机会',
      looking: '积极找工作',
      passive: '暂不考虑'
    }
    return statusMap[v] || v
  }},
  { key: 'source', title: '来源', width: 12 },
  { key: 'createTime', title: '创建时间', width: 20 }
]

export const jobExportColumns: ExportColumn[] = [
  { key: 'title', title: '职位名称', width: 20 },
  { key: 'department', title: '部门', width: 15 },
  { key: 'location', title: '工作地点', width: 12 },
  { key: 'type', title: '工作类型', width: 10, formatter: (v) => {
    const typeMap: Record<string, string> = {
      fulltime: '全职',
      parttime: '兼职',
      intern: '实习',
      contract: '合同'
    }
    return typeMap[v] || v
  }},
  { key: 'salaryMin', title: '最低薪资', width: 10 },
  { key: 'salaryMax', title: '最高薪资', width: 10 },
  { key: 'experience', title: '经验要求', width: 12 },
  { key: 'education', title: '学历要求', width: 10 },
  { key: 'headcount', title: '招聘人数', width: 10 },
  { key: 'status', title: '状态', width: 10, formatter: (v) => {
    const statusMap: Record<string, string> = {
      open: '招聘中',
      paused: '暂停',
      closed: '已关闭',
      filled: '已满员'
    }
    return statusMap[v] || v
  }},
  { key: 'urgent', title: '是否紧急', width: 10, formatter: (v) => v ? '是' : '否' },
  { key: 'publishDate', title: '发布日期', width: 15 },
  { key: 'deadline', title: '截止日期', width: 15 }
]

export const resumeExportColumns: ExportColumn[] = [
  { key: 'candidateName', title: '候选人姓名', width: 12 },
  { key: 'position', title: '应聘职位', width: 20 },
  { key: 'phone', title: '电话', width: 15 },
  { key: 'email', title: '邮箱', width: 25 },
  { key: 'education', title: '学历', width: 10 },
  { key: 'experience', title: '工作经验', width: 12 },
  { key: 'status', title: '状态', width: 10, formatter: (v) => {
    const statusMap: Record<string, string> = {
      pending: '待筛选',
      reviewing: '筛选中',
      interviewed: '已面试',
      offered: '已发offer',
      hired: '已录用',
      rejected: '已拒绝'
    }
    return statusMap[v] || v
  }},
  { key: 'matchScore', title: '匹配度', width: 10, formatter: (v) => v ? `${v}%` : '' },
  { key: 'source', title: '来源', width: 12 },
  { key: 'submitDate', title: '投递日期', width: 15 }
]

export const interviewExportColumns: ExportColumn[] = [
  { key: 'candidateName', title: '候选人', width: 12 },
  { key: 'position', title: '应聘职位', width: 20 },
  { key: 'date', title: '面试日期', width: 12 },
  { key: 'time', title: '面试时间', width: 10 },
  { key: 'type', title: '面试类型', width: 10, formatter: (v) => {
    const typeMap: Record<string, string> = {
      initial: '初试',
      second: '复试',
      final: '终面',
      hr: 'HR面'
    }
    return typeMap[v] || v
  }},
  { key: 'interviewer', title: '面试官', width: 12 },
  { key: 'method', title: '面试方式', width: 10, formatter: (v) => {
    const methodMap: Record<string, string> = {
      onsite: '现场',
      video: '视频',
      phone: '电话'
    }
    return methodMap[v] || v
  }},
  { key: 'location', title: '地点/链接', width: 25 },
  { key: 'status', title: '状态', width: 10, formatter: (v) => {
    const statusMap: Record<string, string> = {
      scheduled: '待进行',
      completed: '已完成',
      cancelled: '已取消'
    }
    return statusMap[v] || v
  }},
  { key: 'notes', title: '备注', width: 30 }
]
