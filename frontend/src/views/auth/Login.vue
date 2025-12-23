<template>
  <div class="login-page">
    <!-- 背景动画粒子 -->
    <div class="bg-particles">
      <div v-for="i in 20" :key="i" class="particle" :style="getParticleStyle(i)"></div>
    </div>

    <div class="login-container" :class="{ 'animate-in': isLoaded }">
      <div class="login-left">
        <div class="gradient-overlay"></div>
        <div class="content">
          <div class="logo-icon">
            <el-icon :size="48"><Connection /></el-icon>
          </div>
          <h1 class="title">智能人才运营平台</h1>
          <p class="subtitle">高效管理人才，智能匹配职位</p>
          <div class="features">
            <div v-for="(feature, index) in features" :key="index"
                 class="feature-item"
                 :style="{ animationDelay: `${index * 0.2}s` }">
              <el-icon><Check /></el-icon>
              <span>{{ feature }}</span>
            </div>
          </div>
          <div class="stats">
            <div class="stat-item">
              <span class="stat-number">10K+</span>
              <span class="stat-label">人才库</span>
            </div>
            <div class="stat-item">
              <span class="stat-number">500+</span>
              <span class="stat-label">企业入驻</span>
            </div>
            <div class="stat-item">
              <span class="stat-number">98%</span>
              <span class="stat-label">匹配成功</span>
            </div>
          </div>
        </div>
      </div>

      <div class="login-right">
        <div class="login-form-wrapper">
          <div class="form-header">
            <h2>欢迎回来</h2>
            <p class="desc">请输入您的账号和密码登录系统</p>
          </div>

          <el-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="rules"
            class="login-form"
            @submit.prevent="handleLogin"
          >
            <el-form-item prop="username">
              <el-input
                v-model="loginForm.username"
                placeholder="用户名 / 邮箱"
                size="large"
                :prefix-icon="User"
                class="custom-input"
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
                class="custom-input"
                @keyup.enter="handleLogin"
              />
            </el-form-item>

            <div class="form-options">
              <el-checkbox v-model="loginForm.remember">记住密码</el-checkbox>
              <a href="#" class="forgot-link">忘记密码？</a>
            </div>

            <el-button
              type="primary"
              size="large"
              class="login-btn"
              :loading="loading"
              @click="handleLogin"
            >
              <span v-if="!loading">登录</span>
              <span v-else>登录中...</span>
            </el-button>
          </el-form>

          <div class="divider">
            <span>或</span>
          </div>

          <div class="social-login">
            <div class="social-btn wechat" title="微信登录">
              <el-icon><ChatDotRound /></el-icon>
              <span class="social-label">微信</span>
            </div>
            <div class="social-btn dingtalk" title="钉钉登录">
              <el-icon><Message /></el-icon>
              <span class="social-label">钉钉</span>
            </div>
            <div class="social-btn wecom" title="企业微信">
              <el-icon><OfficeBuilding /></el-icon>
              <span class="social-label">企微</span>
            </div>
          </div>

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
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { User, Lock, Check, Connection, ChatDotRound, Message, OfficeBuilding } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const isLoaded = ref(false)
const loginFormRef = ref<FormInstance>()

const features = ['智能推荐算法', '简历智能解析', '数据可视化分析', '多角色权限管理']

const loginForm = reactive({
  username: '',
  password: '',
  remember: false
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

// 生成随机粒子样式
const getParticleStyle = (_index: number) => {
  const size = Math.random() * 10 + 5
  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${Math.random() * 100}%`,
    top: `${Math.random() * 100}%`,
    animationDelay: `${Math.random() * 5}s`,
    animationDuration: `${Math.random() * 10 + 10}s`
  }
}

// 加载记住的密码
const loadRememberedCredentials = () => {
  const remembered = localStorage.getItem('rememberedLogin')
  if (remembered) {
    try {
      const data = JSON.parse(remembered)
      loginForm.username = data.username || ''
      loginForm.password = data.password || ''
      loginForm.remember = true
    } catch (e) {
      localStorage.removeItem('rememberedLogin')
    }
  }
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await userStore.login(loginForm.username, loginForm.password)

        // 处理记住密码
        if (loginForm.remember) {
          localStorage.setItem('rememberedLogin', JSON.stringify({
            username: loginForm.username,
            password: loginForm.password
          }))
        } else {
          localStorage.removeItem('rememberedLogin')
        }

        ElMessage.success('登录成功，欢迎回来！')
        router.push('/dashboard')
      } catch (error: any) {
        ElMessage.error(error.message || '登录失败，请检查账号密码')
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  loadRememberedCredentials()
  // 触发入场动画
  setTimeout(() => {
    isLoaded.value = true
  }, 100)
})
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%);
  position: relative;
  overflow: hidden;
}

// 背景粒子动画
.bg-particles {
  position: absolute;
  inset: 0;
  overflow: hidden;

  .particle {
    position: absolute;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 50%;
    animation: float linear infinite;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  90% {
    opacity: 1;
  }
  100% {
    transform: translateY(-100vh) rotate(720deg);
    opacity: 0;
  }
}

.login-container {
  display: flex;
  width: 1000px;
  min-height: 600px;
  background: white;
  border-radius: 24px;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.15);
  overflow: hidden;
  opacity: 0;
  transform: translateY(30px) scale(0.95);
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);

  &.animate-in {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.login-left {
  flex: 1.1;
  position: relative;
  background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%);

  .gradient-overlay {
    position: absolute;
    inset: 0;
    background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  }

  .content {
    position: relative;
    z-index: 1;
    padding: 60px 50px;
    color: white;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;

    .logo-icon {
      width: 80px;
      height: 80px;
      background: rgba(255, 255, 255, 0.2);
      border-radius: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-bottom: 24px;
      backdrop-filter: blur(10px);
    }

    .title {
      font-size: 36px;
      font-weight: 700;
      margin-bottom: 12px;
      letter-spacing: 2px;
    }

    .subtitle {
      font-size: 18px;
      opacity: 0.9;
      margin-bottom: 40px;
    }

    .features {
      display: flex;
      flex-direction: column;
      gap: 16px;
      margin-bottom: 40px;

      .feature-item {
        display: flex;
        align-items: center;
        gap: 12px;
        font-size: 16px;
        opacity: 0;
        animation: slideInFeature 0.5s ease forwards;

        .el-icon {
          font-size: 20px;
          background: rgba(255, 255, 255, 0.2);
          padding: 6px;
          border-radius: 50%;
        }
      }
    }

    .stats {
      display: flex;
      gap: 32px;
      padding-top: 24px;
      border-top: 1px solid rgba(255, 255, 255, 0.2);

      .stat-item {
        display: flex;
        flex-direction: column;

        .stat-number {
          font-size: 28px;
          font-weight: 700;
        }

        .stat-label {
          font-size: 14px;
          opacity: 0.8;
        }
      }
    }
  }
}

@keyframes slideInFeature {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.login-right {
  flex: 0.9;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 50px;
  background: #ffffff;
}

.login-form-wrapper {
  width: 100%;
  max-width: 380px;

  .form-header {
    margin-bottom: 32px;

    h2 {
      font-size: 32px;
      font-weight: 700;
      color: #1e293b;
      margin-bottom: 8px;
    }

    .desc {
      color: #64748b;
      font-size: 15px;
    }
  }

  .login-form {
    .custom-input {
      :deep(.el-input__wrapper) {
        border-radius: 12px;
        padding: 4px 16px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
        border: 2px solid transparent;
        transition: all 0.3s;

        &:hover {
          box-shadow: 0 4px 12px rgba(14, 165, 233, 0.1);
        }

        &.is-focus {
          border-color: #0ea5e9;
          box-shadow: 0 4px 16px rgba(14, 165, 233, 0.15);
        }
      }
    }

    .form-options {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 24px;

      .forgot-link {
        color: #0ea5e9;
        font-size: 14px;
        text-decoration: none;
        transition: color 0.3s;

        &:hover {
          color: #0284c7;
        }
      }
    }

    .login-btn {
      width: 100%;
      height: 52px;
      border-radius: 12px;
      font-size: 16px;
      font-weight: 600;
      background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%);
      border: none;
      transition: opacity 0.3s ease;

      &:hover {
        opacity: 0.9;
      }

      &:active {
        opacity: 0.8;
      }
    }
  }

  .divider {
    display: flex;
    align-items: center;
    margin: 28px 0;

    &::before,
    &::after {
      content: '';
      flex: 1;
      height: 1px;
      background: #e2e8f0;
    }

    span {
      padding: 0 16px;
      color: #94a3b8;
      font-size: 14px;
    }
  }

  .social-login {
    display: flex;
    justify-content: center;
    gap: 20px;
    margin-bottom: 28px;

    .social-btn {
      width: 80px;
      height: 60px;
      border-radius: 12px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 4px;
      cursor: pointer;
      transition: opacity 0.3s ease;
      font-size: 22px;

      .social-label {
        font-size: 12px;
        font-weight: 500;
      }

      &:hover {
        opacity: 0.8;
      }

      &.wechat {
        background: #e0f2fe;
        color: #0ea5e9;
      }

      &.dingtalk {
        background: #ecfeff;
        color: #06b6d4;
      }

      &.wecom {
        background: #f0fdfa;
        color: #14b8a6;
      }
    }
  }

  .footer-links {
    text-align: center;
    color: #64748b;
    font-size: 15px;

    a {
      color: #0ea5e9;
      text-decoration: none;
      font-weight: 600;
      margin-left: 6px;
      transition: color 0.3s;

      &:hover {
        color: #0284c7;
      }
    }
  }
}

// 响应式适配
@media (max-width: 900px) {
  .login-container {
    flex-direction: column;
    width: 95%;
    max-width: 450px;
    min-height: auto;
  }

  .login-left {
    padding: 40px 30px;

    .content {
      padding: 30px;

      .title {
        font-size: 28px;
      }

      .stats {
        display: none;
      }
    }
  }

  .login-right {
    padding: 30px;
  }
}
</style>
