<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-left">
        <div class="gradient-overlay"></div>
        <div class="content">
          <h1 class="title">智能人才运营平台</h1>
          <p class="subtitle">高效管理人才，智能匹配职位</p>
          <div class="features">
            <div class="feature-item">
              <el-icon><Check /></el-icon>
              <span>智能推荐算法</span>
            </div>
            <div class="feature-item">
              <el-icon><Check /></el-icon>
              <span>简历智能解析</span>
            </div>
            <div class="feature-item">
              <el-icon><Check /></el-icon>
              <span>数据可视化分析</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="login-right">
        <div class="login-form-wrapper">
          <h2>欢迎登录</h2>
          <p class="desc">请输入您的账号和密码</p>
          
          <el-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="rules"
            class="login-form"
          >
            <el-form-item prop="username">
              <el-input
                v-model="loginForm.username"
                placeholder="用户名 / 邮箱"
                size="large"
                :prefix-icon="User"
              />
            </el-form-item>
            
            <el-form-item prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="密码"
                size="large"
                :prefix-icon="Lock"
                show-password
                @keyup.enter="handleLogin"
              />
            </el-form-item>
            
            <el-button
              type="primary"
              size="large"
              class="login-btn"
              :loading="loading"
              @click="handleLogin"
            >
              登录
            </el-button>
          </el-form>
          
          <div class="footer-links">
            <span>还没有账号？</span>
            <router-link to="/register">立即注册</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { User, Lock, Check } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const loginFormRef = ref<FormInstance>()

const loginForm = reactive({
  username: '',
  password: ''
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await userStore.login(loginForm.username, loginForm.password)
        ElMessage.success('登录成功')
        router.push('/dashboard')
      } catch (error: any) {
        ElMessage.error(error.message || '登录失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-container {
  display: flex;
  width: 900px;
  height: 550px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

.login-left {
  flex: 1;
  position: relative;
  background: url('https://images.unsplash.com/photo-1522071820081-009f0129c71c?auto=format&fit=crop&w=800') center/cover;
  
  .gradient-overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.9) 0%, rgba(118, 75, 162, 0.9) 100%);
  }
  
  .content {
    position: relative;
    z-index: 1;
    padding: 60px 40px;
    color: white;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    
    .title {
      font-size: 36px;
      font-weight: 700;
      margin-bottom: 16px;
    }
    
    .subtitle {
      font-size: 18px;
      opacity: 0.9;
      margin-bottom: 48px;
    }
    
    .features {
      display: flex;
      flex-direction: column;
      gap: 16px;
      
      .feature-item {
        display: flex;
        align-items: center;
        gap: 12px;
        font-size: 16px;
        
        .el-icon {
          font-size: 20px;
        }
      }
    }
  }
}

.login-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.login-form-wrapper {
  width: 100%;
  max-width: 360px;
  
  h2 {
    font-size: 28px;
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: 8px;
  }
  
  .desc {
    color: var(--text-secondary);
    margin-bottom: 32px;
  }
  
  .login-form {
    .login-btn {
      width: 100%;
      margin-top: 8px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border: none;
      
      &:hover {
        opacity: 0.9;
      }
    }
  }
  
  .footer-links {
    margin-top: 24px;
    text-align: center;
    color: var(--text-secondary);
    
    a {
      color: var(--primary-color);
      text-decoration: none;
      margin-left: 8px;
      font-weight: 500;
      
      &:hover {
        text-decoration: underline;
      }
    }
  }
}
</style>
