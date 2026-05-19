<template>
  <div class="chat-input">
    <el-input
      v-model="inputContent"
      type="textarea"
      :rows="2"
      :disabled="disabled"
      :placeholder="disabled ? '无法发送消息' : '输入消息，按 Enter 发送...'"
      resize="none"
      @keydown.enter.exact.prevent="handleSend"
    />
    <el-button
      type="primary"
      :disabled="disabled || !inputContent.trim()"
      @click="handleSend"
    >
      <el-icon><Promotion /></el-icon>
      发送
    </el-button>
  </div>
</template>

<script setup lang="ts">
/**
 * ChatInput.vue - 聊天输入框组件
 *
 * 文本输入框和发送按钮组件，支持回车发送，并在发送时向父组件抛出事件。
 */
import { ref } from 'vue'
import { Promotion } from '@element-plus/icons-vue'

defineProps<{
  /** 是否禁用输入框 */
  disabled?: boolean
}>()

const emit = defineEmits<{
  /** 用户发送消息时触发 */
  (e: 'send', content: string): void
}>()

/** 当前输入内容 */
const inputContent = ref('')

/**
 * 处理发送动作：校验非空后触发 send 事件。
 */
const handleSend = () => {
  const content = inputContent.value.trim()
  if (!content) return
  
  emit('send', content)
  inputContent.value = ''
}

/**
 * Clear the input content
 * Can be called by parent component
 */
const clearInput = () => {
  inputContent.value = ''
}

/**
 * Focus the input
 * Can be called by parent component
 */
const focus = () => {
  // Focus will be handled by the textarea
}

// Expose methods to parent
defineExpose({
  clearInput,
  focus
})
</script>

<style scoped lang="scss">
.chat-input {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: #fff;
  border-top: 1px solid #e2e8f0;
  align-items: flex-end;

  :deep(.el-textarea__inner) {
    border-radius: 8px;
    padding: 10px 12px;
    font-size: 14px;
    line-height: 1.5;
    
    &:focus {
      border-color: #0ea5e9;
      box-shadow: 0 0 0 2px rgba(14, 165, 233, 0.1);
    }
  }

  .el-button {
    height: 40px;
    padding: 0 20px;
    border-radius: 8px;
    
    .el-icon {
      margin-right: 4px;
    }
  }
}
</style>
