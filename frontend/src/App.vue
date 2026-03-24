<template>
  <router-view />
  <router-link
    v-if="showDevLogsEntry"
    to="/dev-logs"
    class="dev-logs-entry"
  >
    <span class="entry-dot" />
    <span class="entry-text">
      <strong>后台运行终端日志</strong>
      <small>无需登录即可查看</small>
    </span>
  </router-link>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useThemeStore } from '@/store/theme'

const themeStore = useThemeStore()
const route = useRoute()
const showDevLogsEntry = computed(() => route.path !== '/dev-logs')

// 全局初始化主题
onMounted(() => {
  themeStore.init()
})
</script>

<style>
/* Global app styles handled in global.scss */
.dev-logs-entry {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 2500;
  display: inline-flex;
  align-items: center;
  gap: 12px;
  min-width: 188px;
  padding: 14px 18px;
  border-radius: 18px;
  background: linear-gradient(135deg, rgba(14, 101, 87, 0.96), rgba(7, 62, 53, 0.94));
  box-shadow: 0 16px 36px rgba(9, 59, 51, 0.28);
  color: #f8fbfa;
  text-decoration: none;
  transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.dev-logs-entry:hover {
  transform: translateY(-2px);
  box-shadow: 0 20px 42px rgba(9, 59, 51, 0.34);
}

.entry-dot {
  width: 12px;
  height: 12px;
  border-radius: 999px;
  background: #7ef0c5;
  box-shadow: 0 0 0 6px rgba(126, 240, 197, 0.18);
  flex-shrink: 0;
}

.entry-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  line-height: 1.2;
}

.entry-text strong {
  font-size: 14px;
  font-weight: 700;
}

.entry-text small {
  font-size: 12px;
  color: rgba(248, 251, 250, 0.78);
}

@media (max-width: 768px) {
  .dev-logs-entry {
    right: 14px;
    bottom: 14px;
    min-width: 0;
    padding: 12px 14px;
  }

  .entry-text strong {
    font-size: 13px;
  }

  .entry-text small {
    font-size: 11px;
  }
}
</style>
