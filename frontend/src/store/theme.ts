import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'

export const useThemeStore = defineStore('theme', () => {
  // 当前主题模式
  const mode = ref<ThemeMode>('system')
  // 实际应用的主题（light 或 dark）
  const actualTheme = ref<'light' | 'dark'>('light')

  // 获取系统主题偏好
  const getSystemTheme = (): 'light' | 'dark' => {
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      return 'dark'
    }
    return 'light'
  }

  // 应用主题到 DOM
  const applyTheme = (theme: 'light' | 'dark') => {
    const html = document.documentElement

    // 添加过渡类
    html.classList.add('theme-transition')

    if (theme === 'dark') {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }

    actualTheme.value = theme

    // 移除过渡类（延迟以确保过渡完成）
    setTimeout(() => {
      html.classList.remove('theme-transition')
    }, 300)
  }

  // 更新主题
  const updateTheme = () => {
    if (mode.value === 'system') {
      applyTheme(getSystemTheme())
    } else {
      applyTheme(mode.value)
    }
  }

  // 设置主题模式
  const setMode = (newMode: ThemeMode) => {
    mode.value = newMode
    localStorage.setItem('theme-mode', newMode)
    updateTheme()
  }

  // 切换主题（在 light 和 dark 之间切换）
  const toggleTheme = () => {
    if (actualTheme.value === 'light') {
      setMode('dark')
    } else {
      setMode('light')
    }
  }

  // 初始化
  const init = () => {
    // 从 localStorage 读取
    const savedMode = localStorage.getItem('theme-mode') as ThemeMode | null
    if (savedMode && ['light', 'dark', 'system'].includes(savedMode)) {
      mode.value = savedMode
    }

    // 应用主题
    updateTheme()

    // 监听系统主题变化
    if (window.matchMedia) {
      window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
        if (mode.value === 'system') {
          applyTheme(e.matches ? 'dark' : 'light')
        }
      })
    }
  }

  // 监听 mode 变化
  watch(mode, updateTheme)

  return {
    mode,
    actualTheme,
    setMode,
    toggleTheme,
    init
  }
})
