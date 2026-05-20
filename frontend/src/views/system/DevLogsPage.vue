<template>
  <div class="dev-logs-page">
    <div class="page-shell">
      <section class="page-hero">
        <div class="hero-copy">
          <p class="eyebrow">Development Console</p>
          <h1>后台运行终端日志</h1>
          <p class="subtitle">
            这个页面直接读取本地开发日志，不再依赖 8091 日志面板服务。
          </p>
        </div>
        <div class="hero-actions">
          <el-button type="primary" :loading="loading" @click="refreshAll">
            刷新日志
          </el-button>
          <el-button @click="router.back()">返回上一页</el-button>
        </div>
      </section>

      <section class="summary-grid">
        <div class="summary-card">
          <strong>{{ runningCount }}/{{ services.length }}</strong>
          <span>正在运行服务</span>
        </div>
        <div class="summary-card">
          <strong>{{ selectedService?.name || '-' }}</strong>
          <span>当前查看服务</span>
        </div>
        <div class="summary-card">
          <strong>{{ lastUpdatedText }}</strong>
          <span>最后刷新时间</span>
        </div>
      </section>

      <el-alert
        v-if="errorMessage"
        class="error-alert"
        type="error"
        :closable="false"
        :title="errorMessage"
      />

      <section class="logs-layout">
        <aside class="service-list">
          <button
            v-for="service in services"
            :key="service.name"
            class="service-item"
            :class="{ active: service.name === selectedName }"
            type="button"
            @click="selectService(service.name)"
          >
            <div class="service-main">
              <span class="status-dot" :class="{ running: service.running }" />
              <strong>{{ service.name }}</strong>
            </div>
            <div class="service-meta">
              <span>端口 {{ service.port }}</span>
              <span>{{ formatSize(service.logSize) }}</span>
            </div>
            <div class="service-meta">
              <span>PID {{ service.pid || '-' }}</span>
              <span>{{ service.running ? '运行中' : '未运行' }}</span>
            </div>
          </button>
        </aside>

        <main class="log-panel">
          <div class="log-toolbar">
            <div>
              <h2>{{ selectedService?.name || '请选择服务' }}</h2>
              <p v-if="selectedService">
                {{ selectedService.logFile }}
              </p>
            </div>
            <div class="toolbar-actions">
              <el-select v-model="lineCount" size="small" class="line-select" @change="loadLog">
                <el-option :value="100" label="最近 100 行" />
                <el-option :value="300" label="最近 300 行" />
                <el-option :value="800" label="最近 800 行" />
                <el-option :value="1500" label="最近 1500 行" />
              </el-select>
              <el-button size="small" :loading="logLoading" @click="loadLog">
                刷新当前日志
              </el-button>
            </div>
          </div>

          <pre class="log-content">{{ logContent || '暂无日志内容' }}</pre>
        </main>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'

interface DevService {
  name: string
  port: number
  url: string
  pid: number | null
  running: boolean
  logFile: string
  pidFile: string
  logSize: number
  updatedAt: string | null
}

interface ServicesResponse {
  generatedAt: string
  services: DevService[]
}

interface LogResponse {
  service: string
  content: string
}

const router = useRouter()
const services = ref<DevService[]>([])
const selectedName = ref('')
const logContent = ref('')
const loading = ref(false)
const logLoading = ref(false)
const errorMessage = ref('')
const lineCount = ref(300)
const lastUpdatedAt = ref<Date | null>(null)
let refreshTimer: number | undefined

const selectedService = computed(() =>
  services.value.find((service) => service.name === selectedName.value),
)

const runningCount = computed(() =>
  services.value.filter((service) => service.running).length,
)

const lastUpdatedText = computed(() => {
  if (!lastUpdatedAt.value) {
    return '-'
  }
  return lastUpdatedAt.value.toLocaleTimeString()
})

const refreshAll = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const response = await fetch('/__dev_logs/api/services')
    if (!response.ok) {
      throw new Error(`服务列表加载失败: ${response.status}`)
    }

    const data = (await response.json()) as ServicesResponse
    services.value = data.services
    lastUpdatedAt.value = new Date(data.generatedAt)

    if (!selectedName.value && data.services.length > 0) {
      const firstRunning = data.services.find((service) => service.running)
      selectedName.value = firstRunning?.name || data.services[0].name
    }

    await loadLog()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '日志加载失败'
  } finally {
    loading.value = false
  }
}

const selectService = async (name: string) => {
  selectedName.value = name
  await loadLog()
}

const loadLog = async () => {
  if (!selectedName.value) {
    logContent.value = ''
    return
  }

  logLoading.value = true
  try {
    const params = new URLSearchParams({
      name: selectedName.value,
      lines: String(lineCount.value),
    })
    const response = await fetch(`/__dev_logs/api/log?${params.toString()}`)
    if (!response.ok) {
      throw new Error(`日志读取失败: ${response.status}`)
    }

    const data = (await response.json()) as LogResponse
    logContent.value = data.content
  } catch (error) {
    logContent.value = error instanceof Error ? error.message : '日志读取失败'
  } finally {
    logLoading.value = false
  }
}

const formatSize = (size: number) => {
  if (size < 1024) {
    return `${size} B`
  }
  if (size < 1024 * 1024) {
    return `${(size / 1024).toFixed(1)} KB`
  }
  return `${(size / 1024 / 1024).toFixed(1)} MB`
}

onMounted(() => {
  refreshAll()
  refreshTimer = window.setInterval(refreshAll, 10000)
})

onUnmounted(() => {
  if (refreshTimer) {
    window.clearInterval(refreshTimer)
  }
})
</script>

<style scoped lang="scss">
.dev-logs-page {
  min-height: 100vh;
  padding: 28px 20px 40px;
  background:
    radial-gradient(circle at top left, rgba(20, 108, 92, 0.14), transparent 30%),
    radial-gradient(circle at top right, rgba(213, 153, 58, 0.12), transparent 26%),
    linear-gradient(180deg, #faf6ee 0%, #f2ede3 100%);
}

.page-shell {
  width: min(1440px, 100%);
  margin: 0 auto;
}

.page-hero,
.summary-card,
.logs-layout {
  border: 1px solid rgba(67, 53, 34, 0.12);
  border-radius: 24px;
  background: rgba(255, 252, 246, 0.94);
  box-shadow: 0 20px 50px rgba(92, 70, 35, 0.12);
  backdrop-filter: blur(12px);
}

.page-hero {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  align-items: flex-start;
  padding: 28px;
}

.eyebrow {
  margin: 0 0 8px;
  color: #0f6b5a;
  font-size: 12px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  font-weight: 700;
}

h1,
h2,
.subtitle,
p {
  margin: 0;
}

h1 {
  font-size: clamp(28px, 4vw, 42px);
  line-height: 1.05;
  color: #1d2731;
}

h2 {
  color: #1d2731;
  font-size: 20px;
}

.subtitle {
  max-width: 720px;
  margin-top: 12px;
  color: #60707b;
  line-height: 1.7;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 12px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin-top: 18px;
}

.summary-card {
  padding: 18px 22px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.summary-card strong {
  color: #0f6b5a;
  font-size: 24px;
}

.summary-card span {
  color: #60707b;
}

.error-alert {
  margin-top: 18px;
}

.logs-layout {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  gap: 16px;
  margin-top: 18px;
  padding: 14px;
}

.service-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: calc(100vh - 300px);
  overflow: auto;
}

.service-item {
  border: 1px solid rgba(67, 53, 34, 0.12);
  border-radius: 16px;
  padding: 14px;
  background: #fffdf8;
  color: inherit;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
}

.service-item:hover,
.service-item.active {
  border-color: rgba(15, 107, 90, 0.42);
  box-shadow: 0 10px 26px rgba(15, 107, 90, 0.12);
}

.service-main,
.service-meta,
.log-toolbar,
.toolbar-actions {
  display: flex;
  align-items: center;
}

.service-main {
  gap: 8px;
  margin-bottom: 8px;
}

.service-main strong {
  color: #1d2731;
}

.service-meta {
  justify-content: space-between;
  gap: 10px;
  color: #60707b;
  font-size: 12px;
  line-height: 1.7;
}

.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #c8c1b5;
  box-shadow: 0 0 0 3px rgba(200, 193, 181, 0.16);
}

.status-dot.running {
  background: #17a673;
  box-shadow: 0 0 0 3px rgba(23, 166, 115, 0.16);
}

.log-panel {
  min-width: 0;
  border-radius: 18px;
  background: #111827;
  overflow: hidden;
}

.log-toolbar {
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  background: #1f2937;
}

.log-toolbar h2 {
  color: #f8fafc;
}

.log-toolbar p {
  margin-top: 6px;
  color: #9ca3af;
  font-size: 12px;
  word-break: break-all;
}

.toolbar-actions {
  gap: 10px;
  flex-shrink: 0;
}

.line-select {
  width: 140px;
}

.log-content {
  min-height: calc(100vh - 390px);
  max-height: calc(100vh - 290px);
  margin: 0;
  padding: 18px;
  overflow: auto;
  color: #d1d5db;
  background: #111827;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 13px;
  line-height: 1.65;
  white-space: pre-wrap;
}

@media (max-width: 960px) {
  .page-hero,
  .summary-grid,
  .logs-layout {
    grid-template-columns: 1fr;
    display: grid;
  }

  .hero-actions,
  .log-toolbar,
  .toolbar-actions {
    align-items: flex-start;
    justify-content: flex-start;
    flex-wrap: wrap;
  }

  .service-list {
    max-height: none;
  }
}

@media (max-width: 640px) {
  .dev-logs-page {
    padding: 18px 12px 28px;
  }

  .page-hero,
  .summary-card,
  .logs-layout {
    border-radius: 18px;
  }

  .page-hero {
    padding: 20px;
  }
}
</style>
