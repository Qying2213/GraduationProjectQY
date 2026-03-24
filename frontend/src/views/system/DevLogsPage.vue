<template>
  <div class="dev-logs-page">
    <div class="page-shell">
      <section class="page-hero">
        <div class="hero-copy">
          <p class="eyebrow">Development Console</p>
          <h1>后台运行终端日志</h1>
          <p class="subtitle">
            这个页面不需要登录，直接聚合查看各个后台服务的端口、运行状态和终端日志。
          </p>
        </div>
        <div class="hero-actions">
          <el-button type="primary" @click="reloadFrame">刷新日志面板</el-button>
          <el-button @click="openInNewTab">新标签页打开</el-button>
          <el-button @click="router.back()">返回上一页</el-button>
        </div>
      </section>

      <section class="hint-card">
        <div>
          <strong>日志面板地址</strong>
          <span>{{ dashboardUrl }}</span>
        </div>
        <div>
          <strong>启动方式</strong>
          <span>先执行项目根目录下的 `./start-all.sh`</span>
        </div>
      </section>

      <section class="iframe-shell">
        <iframe
          :key="frameKey"
          :src="dashboardUrl"
          title="后台运行终端日志"
          class="dashboard-frame"
        />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const dashboardUrl = '/__dev_dashboard/'
const directDashboardUrl = 'http://localhost:8091'
const frameKey = ref(0)

const reloadFrame = () => {
  frameKey.value += 1
}

const openInNewTab = () => {
  window.open(directDashboardUrl, '_blank', 'noopener,noreferrer')
}
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
.hint-card,
.iframe-shell {
  border: 1px solid rgba(67, 53, 34, 0.12);
  border-radius: 24px;
  background: rgba(255, 252, 246, 0.92);
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
.subtitle,
p {
  margin: 0;
}

h1 {
  font-size: clamp(28px, 4vw, 42px);
  line-height: 1.05;
  color: #1d2731;
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

.hint-card {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  margin-top: 18px;
  padding: 18px 22px;
}

.hint-card div {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.hint-card strong {
  color: #1d2731;
  font-size: 14px;
}

.hint-card span {
  color: #60707b;
  line-height: 1.6;
  word-break: break-all;
}

.iframe-shell {
  margin-top: 18px;
  padding: 14px;
}

.dashboard-frame {
  width: 100%;
  min-height: calc(100vh - 230px);
  border: none;
  border-radius: 18px;
  background: #ffffff;
}

@media (max-width: 960px) {
  .page-hero,
  .hint-card {
    grid-template-columns: 1fr;
    display: grid;
  }

  .hero-actions {
    justify-content: flex-start;
  }

  .dashboard-frame {
    min-height: calc(100vh - 280px);
  }
}

@media (max-width: 640px) {
  .dev-logs-page {
    padding: 18px 12px 28px;
  }

  .page-hero,
  .hint-card,
  .iframe-shell {
    border-radius: 18px;
  }

  .page-hero {
    padding: 20px;
  }

  .hint-card {
    padding: 16px;
  }

  .iframe-shell {
    padding: 10px;
  }
}
</style>
