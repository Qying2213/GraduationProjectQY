package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractPhone(t *testing.T) {
	parser := NewResumeParser()

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "标准手机号",
			text:     "联系电话：13812345678",
			expected: "13812345678",
		},
		{
			name:     "手机号在文本中间",
			text:     "张三，手机13912345678，北京",
			expected: "13912345678",
		},
		{
			name:     "无手机号",
			text:     "这是一段没有手机号的文本",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.extractPhone(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractEmail(t *testing.T) {
	parser := NewResumeParser()

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "标准邮箱",
			text:     "邮箱：zhangsan@example.com",
			expected: "zhangsan@example.com",
		},
		{
			name:     "复杂邮箱",
			text:     "联系方式：zhang.san_123@company.com.cn",
			expected: "zhang.san_123@company.com.cn",
		},
		{
			name:     "无邮箱",
			text:     "这是一段没有邮箱的文本",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.extractEmail(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractSkills(t *testing.T) {
	parser := NewResumeParser()

	tests := []struct {
		name        string
		text        string
		minExpected int
	}{
		{
			name:        "前端技能",
			text:        "熟练掌握Vue3、TypeScript、React，了解Webpack和Vite",
			minExpected: 4,
		},
		{
			name:        "后端技能",
			text:        "精通Java、Spring Boot，熟悉MySQL、Redis",
			minExpected: 4,
		},
		{
			name:        "全栈技能",
			text:        "Go语言开发，使用Gin框架，Docker部署，Kubernetes编排",
			minExpected: 3,
		},
		{
			name:        "无技能关键词",
			text:        "这是一段普通文本",
			minExpected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.extractSkills(tt.text)
			assert.GreaterOrEqual(t, len(result), tt.minExpected)
		})
	}
}

func TestExtractEducation(t *testing.T) {
	parser := NewResumeParser()

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "本科学历",
			text:     "学历：本科，毕业于某大学",
			expected: "本科",
		},
		{
			name:     "硕士学历",
			text:     "教育背景：硕士研究生",
			expected: "硕士",
		},
		{
			name:     "博士学历",
			text:     "博士毕业，专业方向为计算机科学",
			expected: "博士",
		},
		{
			name:     "无学历信息",
			text:     "这是一段没有学历信息的文本",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.extractEducation(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractExperience(t *testing.T) {
	parser := NewResumeParser()

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "X年工作经验",
			text:     "5年工作经验，熟悉互联网行业",
			expected: "5年",
		},
		{
			name:     "工作经验：X年",
			text:     "工作经验：3年",
			expected: "3年",
		},
		{
			name:     "无经验信息",
			text:     "应届毕业生",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.extractExperience(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractLocation(t *testing.T) {
	parser := NewResumeParser()

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "北京",
			text:     "现居北京，期望工作地点北京",
			expected: "北京",
		},
		{
			name:     "上海",
			text:     "上海市浦东新区",
			expected: "上海",
		},
		{
			name:     "无城市信息",
			text:     "这是一段没有城市信息的文本",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.extractLocation(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateMatchScore(t *testing.T) {
	parser := NewResumeParser()

	tests := []struct {
		name          string
		resume        *ParsedResume
		jobSkills     []string
		jobExperience int
		jobEducation  string
		minScore      int
		maxScore      int
	}{
		{
			name: "完全匹配",
			resume: &ParsedResume{
				Skills:     []string{"Vue", "TypeScript", "React"},
				Education:  "本科",
				Experience: "5年",
			},
			jobSkills:     []string{"Vue", "TypeScript", "React"},
			jobExperience: 3,
			jobEducation:  "本科",
			minScore:      80,
			maxScore:      100,
		},
		{
			name: "部分匹配",
			resume: &ParsedResume{
				Skills:     []string{"Vue", "JavaScript"},
				Education:  "本科",
				Experience: "2年",
			},
			jobSkills:     []string{"Vue", "TypeScript", "React", "Node.js"},
			jobExperience: 3,
			jobEducation:  "本科",
			minScore:      30,
			maxScore:      60,
		},
		{
			name: "低匹配",
			resume: &ParsedResume{
				Skills:     []string{"Python"},
				Education:  "大专",
				Experience: "1年",
			},
			jobSkills:     []string{"Vue", "TypeScript", "React"},
			jobExperience: 5,
			jobEducation:  "本科",
			minScore:      0,
			maxScore:      30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := parser.CalculateMatchScore(tt.resume, tt.jobSkills, tt.jobExperience, tt.jobEducation)
			assert.GreaterOrEqual(t, score, tt.minScore)
			assert.LessOrEqual(t, score, tt.maxScore)
		})
	}
}

func TestParse(t *testing.T) {
	parser := NewResumeParser()

	resumeText := `
姓名：张三
手机：13812345678
邮箱：zhangsan@example.com
现居：北京

教育背景：
本科 - 计算机科学与技术

工作经验：5年

技能：
- 精通Vue3、TypeScript、React
- 熟悉Node.js、Webpack
- 了解Docker、Kubernetes
`

	result, err := parser.Parse(resumeText)

	assert.NoError(t, err)
	assert.Equal(t, "张三", result.Name)
	assert.Equal(t, "13812345678", result.Phone)
	assert.Equal(t, "zhangsan@example.com", result.Email)
	assert.Equal(t, "北京", result.Location)
	assert.Equal(t, "本科", result.Education)
	assert.Equal(t, "5年", result.Experience)
	assert.GreaterOrEqual(t, len(result.Skills), 4)
}
