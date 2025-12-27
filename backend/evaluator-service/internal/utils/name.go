package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	chineseName = regexp.MustCompile(`[-]*?】\s*([\p{Han}]{2,4})|\s([\p{Han}]{2,4})[\s\d]`)
	chineseWord = regexp.MustCompile(`[\p{Han}]{2,4}`)
)

func ExtractCandidateName(filename string) string {
	namePart := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	exclude := []string{"简历", "个人", "求职", "应聘", "后端", "前端", "开发", "工程师", "实习", "校招", "社招", "深圳", "北京", "上海", "广州", "杭州", "年以上", "年经验", "本科", "硕士", "博士"}
	if m := chineseName.FindStringSubmatch(namePart); len(m) > 1 {
		n := firstNonEmpty(m[1:])
		if n != "" && !inSlice(n, exclude) {
			return n
		}
	}
	if ms := chineseWord.FindAllString(namePart, -1); len(ms) > 0 {
		for _, n := range ms {
			if !inSlice(n, exclude) {
				return n
			}
		}
	}
	s := namePart
	s = strings.NewReplacer("_", " ", "-", " ", "[", " ", "]", " ", "【", " ", "】", " ", "(", " ", ")", " ", "（", " ", "）", " ").Replace(s)
	s = chineseWord.ReplaceAllString(s, "$0")
	parts := strings.Fields(s)
	if len(parts) > 0 {
		return trimLen(parts[0], 20)
	}
	return "候选人"
}

func firstNonEmpty(s []string) string {
	for _, v := range s {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
func inSlice(s string, arr []string) bool {
	for _, v := range arr {
		if s == v {
			return true
		}
	}
	return false
}
func trimLen(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}
