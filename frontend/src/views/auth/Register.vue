<template>
  <div class="register-page">
    <!-- 背景动画粒子 -->
    <div class="bg-particles">
      <div v-for="i in 20" :key="i" class="particle" :style="getParticleStyle()"></div>
    </div>

    <div class="register-container" :class="{ 'animate-in': isLoaded }">
      <div class="register-left">
        <div class="gradient-overlay"></div>
        <div class="content">
          <div class="logo-icon">
            <el-icon :size="48"><UserFilled /></el-icon>
          </div>
          <h1 class="title">加入我们</h1>
          <p class="subtitle">开启智能人才管理之旅</p>

          <div class="benefits">
            <div class="benefit-item" v-for="(benefit, index) in benefits" :key="index"
                 :style="{ animationDelay: `${index * 0.15}s` }">
              <div class="benefit-icon">
                <el-icon><component :is="benefit.icon" /></el-icon>
              </div>
              <div class="benefit-text">
                <h4>{{ benefit.title }}</h4>
                <p>{{ benefit.desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="register-right">
        <div class="register-form-wrapper">
          <div class="form-header">
            <h2>创建账号</h2>
            <p class="desc">请填写以下信息完成注册</p>
          </div>

          <!-- 步骤指示器 -->
          <div class="step-indicator">
            <div class="step" :class="{ active: currentStep >= 1, completed: currentStep > 1 }">
              <span class="step-number">1</span>
              <span class="step-label">基本信息</span>
            </div>
            <div class="step-line" :class="{ active: currentStep > 1 }"></div>
            <div class="step" :class="{ active: currentStep >= 2 }">
              <span class="step-number">2</span>
              <span class="step-label">角色选择</span>
            </div>
          </div>

          <el-form
            ref="registerFormRef"
            :model="registerForm"
            :rules="rules"
            class="register-form"
          >
            <!-- 步骤 1: 基本信息 -->
            <div v-show="currentStep === 1" class="step-content">
              <el-form-item prop="username">
                <el-input
                  v-model="registerForm.username"
                  placeholder="用户名"
                  size="large"
                  :prefix-icon="User"
                  class="custom-input"
                />
              </el-form-item>

              <el-form-item prop="email">
                <el-input
                  v-model="registerForm.email"
                  placeholder="邮箱地址"
                  size="large"
                  :prefix-icon="Message"
                  class="custom-input"
                />
              </el-form-item>

              <el-form-item prop="password">
                <el-input
                  v-model="registerForm.password"
                  type="password"
                  placeholder="设置密码"
                  size="large"
                  :prefix-icon="Lock"
                  show-password
                  class="custom-input"
                />
              </el-form-item>

              <el-form-item prop="confirmPassword">
                <el-input
                  v-model="registerForm.confirmPassword"
                  type="password"
                  placeholder="确认密码"
                  size="large"
                  :prefix-icon="Lock"
                  show-password
                  class="custom-input"
                />
              </el-form-item>

              <el-button type="primary" size="large" class="next-btn" @click="nextStep">
                下一步
                <el-icon class="el-icon--right"><ArrowRight /></el-icon>
              </el-button>
            </div>

            <!-- 步骤 2: 角色选择 -->
            <div v-show="currentStep === 2" class="step-content">
              <div class="role-selection">
                <div
                  class="role-card"
                  :class="{ selected: registerForm.role === 'hr' }"
                  @click="registerForm.role = 'hr'"
                >
                  <div class="role-icon hr">
                    <el-icon :size="32"><OfficeBuilding /></el-icon>
                  </div>
                  <h3>HR / 招聘方</h3>
                  <p>发布职位、管理候选人、智能匹配人才</p>
                  <div class="check-mark" v-if="registerForm.role === 'hr'">
                    <el-icon><Check /></el-icon>
                  </div>
                </div>

                <div
                  class="role-card"
                  :class="{ selected: registerForm.role === 'candidate' }"
                  @click="registerForm.role = 'candidate'"
                >
                  <div class="role-icon candidate">
                    <el-icon :size="32"><User /></el-icon>
                  </div>
                  <h3>求职者</h3>
                  <p>浏览职位、投递简历、获取推荐机会</p>
                  <div class="check-mark" v-if="registerForm.role === 'candidate'">
                    <el-icon><Check /></el-icon>
                  </div>
                </div>
              </div>

              <el-form-item prop="phone" style="margin-top: 20px">
                <el-input
                  v-model="registerForm.phone"
                  placeholder="手机号码（选填）"
                  size="large"
                  :prefix-icon="Phone"
                  class="custom-input"
                />
              </el-form-item>

              <el-form-item>
                <el-checkbox v-model="registerForm.agreement">
                  我已阅读并同意 <a href="#" class="link">用户协议</a> 和 <a href="#" class="link">隐私政策</a>
                </el-checkbox>
              </el-form-item>

              <div class="btn-group">
                <el-button size="large" class="back-btn" @click="currentStep = 1">
                  <el-icon class="el-icon--left"><ArrowLeft /></el-icon>
                  上一步
                </el-button>
                <el-button
                  type="primary"
                  size="large"
                  class="submit-btn"
                  :loading="loading"
                  :disabled="!registerForm.agreement"
                  @click="handleRegister"
                >
                  <span v-if="!loading">完成注册</span>
                  <span v-else>注册中...</span>
                </el-button>
              </div>
            </div>
          </el-form>

          <div class="footer-links">
            <span>已有账号？</span>
            <router-link to="/login">立即登录</router-link>
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
import {
  User, Lock, Message, UserFilled, OfficeBuilding, Check,
  ArrowRight, ArrowLeft, Phone, Briefcase, TrendCharts
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const isLoaded = ref(false)
const currentStep = ref(1)
const registerFormRef = ref<FormInstance>()

const benefits = [
  { icon: Briefcase, title: '海量职位', desc: '覆盖各行业优质岗位' },
  { icon: TrendCharts, title: '智能匹配', desc: 'AI算法精准推荐' },
  { icon: Lock, title: '安全可靠', desc: '信息加密保护隐私' }
]

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  role: 'candidate',
  phone: '',
  agreement: false
})

const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在3-50个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 生成随机粒子样式
const getParticleStyle = () => {
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

const nextStep = async () => {
  if (!registerFormRef.value) return

  // 验证第一步的字段
  try {
    await registerFormRef.value.validateField(['username', 'email', 'password', 'confirmPassword'])
    currentStep.value = 2
  } catch {
    // 验证失败
  }
}

const handleRegister = async () => {
  if (!registerFormRef.value) return

  if (!registerForm.agreement) {
    ElMessage.warning('请先同意用户协议和隐私政策')
    return
  }

  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await userStore.register({
          username: registerForm.username,
          email: registerForm.email,
          password: registerForm.password,
          role: registerForm.role as 'hr' | 'candidate',
          phone: registerForm.phone
        })
        ElMessage.success('注册成功！即将跳转登录页面')
        setTimeout(() => {
          router.push('/login')
        }, 1500)
      } catch (error: any) {
        ElMessage.error(error.message || '注册失败，请稍后重试')
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  setTimeout(() => {
    isLoaded.value = true
  }, 100)
})
</script>

<style scoped lang="scss">
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
  padding: 20px;
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

.register-container {
  display: flex;
  width: 1000px;
  min-height: 650px;
  background: white;
  border-radius: 24px;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.3);
  overflow: hidden;
  opacity: 0;
  transform: translateY(30px) scale(0.95);
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);

  &.animate-in {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.register-left {
  flex: 1;
  position: relative;
  background: linear-gradient(135deg, #764ba2 0%, #667eea 100%);

  .gradient-overlay {
    position: absolute;
    inset: 0;
    background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  }

  .content {
    position: relative;
    z-index: 1;
    padding: 50px 40px;
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

    .benefits {
      display: flex;
      flex-direction: column;
      gap: 24px;

      .benefit-item {
        display: flex;
        align-items: flex-start;
        gap: 16px;
        opacity: 0;
        animation: slideInBenefit 0.5s ease forwards;

        .benefit-icon {
          width: 48px;
          height: 48px;
          background: rgba(255, 255, 255, 0.2);
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 24px;
          flex-shrink: 0;
        }

        .benefit-text {
          h4 {
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 4px;
          }

          p {
            font-size: 14px;
            opacity: 0.8;
            margin: 0;
          }
        }
      }
    }
  }
}

@keyframes slideInBenefit {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.register-right {
  flex: 1.1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 50px;
  background: #fafbfc;
}

.register-form-wrapper {
  width: 100%;
  max-width: 400px;

  .form-header {
    margin-bottom: 24px;

    h2 {
      font-size: 28px;
      font-weight: 700;
      color: #1a1a2e;
      margin-bottom: 8px;
    }

    .desc {
      color: #6b7280;
      font-size: 15px;
    }
  }

  .step-indicator {
    display: flex;
    align-items: center;
    margin-bottom: 28px;

    .step {
      display: flex;
      align-items: center;
      gap: 8px;

      .step-number {
        width: 28px;
        height: 28px;
        border-radius: 50%;
        background: #e5e7eb;
        color: #9ca3af;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        font-weight: 600;
        transition: all 0.3s;
      }

      .step-label {
        font-size: 14px;
        color: #9ca3af;
        transition: all 0.3s;
      }

      &.active {
        .step-number {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: white;
        }
        .step-label {
          color: #374151;
        }
      }

      &.completed {
        .step-number {
          background: #10b981;
        }
      }
    }

    .step-line {
      flex: 1;
      height: 2px;
      background: #e5e7eb;
      margin: 0 12px;
      transition: background 0.3s;

      &.active {
        background: linear-gradient(90deg, #667eea, #764ba2);
      }
    }
  }

  .register-form {
    .custom-input {
      :deep(.el-input__wrapper) {
        border-radius: 12px;
        padding: 4px 16px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
        border: 2px solid transparent;
        transition: all 0.3s;

        &:hover {
          box-shadow: 0 4px 12px rgba(102, 126, 234, 0.1);
        }

        &.is-focus {
          border-color: #667eea;
          box-shadow: 0 4px 16px rgba(102, 126, 234, 0.2);
        }
      }
    }

    .role-selection {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 16px;
      margin-bottom: 8px;

      .role-card {
        position: relative;
        padding: 24px 16px;
        border: 2px solid #e5e7eb;
        border-radius: 16px;
        text-align: center;
        cursor: pointer;
        transition: all 0.3s;

        &:hover {
          border-color: #c7d2fe;
          transform: translateY(-2px);
        }

        &.selected {
          border-color: #667eea;
          background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
        }

        .role-icon {
          width: 64px;
          height: 64px;
          border-radius: 16px;
          display: flex;
          align-items: center;
          justify-content: center;
          margin: 0 auto 12px;

          &.hr {
            background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
            color: white;
          }

          &.candidate {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
            color: white;
          }
        }

        h3 {
          font-size: 16px;
          font-weight: 600;
          color: #374151;
          margin-bottom: 6px;
        }

        p {
          font-size: 12px;
          color: #9ca3af;
          margin: 0;
          line-height: 1.4;
        }

        .check-mark {
          position: absolute;
          top: 12px;
          right: 12px;
          width: 24px;
          height: 24px;
          border-radius: 50%;
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: white;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 14px;
        }
      }
    }

    .link {
      color: #667eea;
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }

    .next-btn {
      width: 100%;
      height: 50px;
      border-radius: 12px;
      font-size: 16px;
      font-weight: 600;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border: none;
      margin-top: 8px;
      transition: all 0.3s;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
      }
    }

    .btn-group {
      display: flex;
      gap: 12px;
      margin-top: 16px;

      .back-btn {
        flex: 0.4;
        height: 50px;
        border-radius: 12px;
        font-size: 15px;
      }

      .submit-btn {
        flex: 0.6;
        height: 50px;
        border-radius: 12px;
        font-size: 16px;
        font-weight: 600;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        border: none;
        transition: all 0.3s;

        &:hover:not(:disabled) {
          transform: translateY(-2px);
          box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
        }

        &:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }
      }
    }
  }

  .footer-links {
    text-align: center;
    color: #6b7280;
    font-size: 15px;
    margin-top: 24px;

    a {
      color: #667eea;
      text-decoration: none;
      font-weight: 600;
      margin-left: 6px;
      transition: color 0.3s;

      &:hover {
        color: #764ba2;
      }
    }
  }
}

// 响应式适配
@media (max-width: 900px) {
  .register-container {
    flex-direction: column;
    width: 95%;
    max-width: 450px;
    min-height: auto;
  }

  .register-left {
    .content {
      padding: 30px;

      .title {
        font-size: 28px;
      }

      .benefits {
        display: none;
      }
    }
  }

  .register-right {
    padding: 30px;
  }

  .register-form-wrapper {
    .register-form {
      .role-selection {
        grid-template-columns: 1fr;
      }
    }
  }
}
</style>
