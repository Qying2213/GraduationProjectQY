<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useThemeStore } from '@/store/theme'

const themeStore = useThemeStore()

// 全局初始化主题
onMounted(() => {
  themeStore.init()
})
</script>

<style>
/* Global app styles handled in global.scss */
</style>
