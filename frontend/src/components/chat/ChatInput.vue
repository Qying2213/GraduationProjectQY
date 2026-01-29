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
 * Requirements: 8.2 (Send messages), 9.1 (Chat input)
 * 
 * Text input with send button. Supports Enter key to send.
 * Emits event when message is sent.
 */
import { ref } from 'vue'
import { Promotion } from '@element-plus/icons-vue'

defineProps<{
  /** Whether the input is disabled */
  disabled?: boolean
}>()

const emit = defineEmits<{
  /** Emitted when user sends a message */
  (e: 'send', content: string): void
}>()

/** Current input content */
const inputContent = ref('')

/**
 * Handle send action
 * Validates content and emits send event
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
