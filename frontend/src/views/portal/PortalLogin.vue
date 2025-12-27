<template>
  <div class="portal-login">
    <div class="login-container">
      <div class="login-card">
        <div class="card-header">
          <h2>求职者登录</h2>
          <p>登录后可投递简历、查看申请状态</p>
        </div>

        <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleLogin">
          <el-form-item prop="email">
            <el-input v-model="form.email" placeholder="邮箱" size="large" :prefix-icon="Message" />
          </el-form-item>
          <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="密码" size="large"
                      :prefix-icon="Lock" show-password @keyup.enter="handleLogin" />
          </el-form-item>
          <div class="form-options">
            <el-checkbox v-model="form.remember">记住我</el-checkbox>
            <a href="#">忘记密码？</a>
          </div>
          <el-button type="primary" size="large" :loading="loading" @click="handleLogin" style="width: 100%">
            登录
          </el-button>
        </el-form>

        <div class="divider"><span>或</span></div>

        <div class="social-login">
          <el-button size="large" style="flex: 1"><el-icon><ChatDotRound /></el-icon>微信登录</el-button>
        </div>

        <div class="footer-link">
          还没有账号？<router-link to="/portal/register">立即注册</router-link>
        </div>

        <div class="hr-login">
          <router-link to="/login">HR/企业登录入口 →</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { Message, Lock, ChatDotRound } from '@element-plus/icons-vue'
import { authApi } from '@/api/auth'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  email: '',
  password: '',
  remember: false
})

const rules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        // 调用真实登录 API
        const res = await authApi.login({
          username: loginForm.username,
          password: loginForm.password
        })
        
        if (res.data?.code === 0 && res.data.data?.token) {
          // 保存 token
          localStorage.setItem('token', res.data.data.token)
          localStorage.setItem('user', JSON.stringify(res.data.data.user))
          ElMessage.success('登录成功')
          router.push('/portal')
        } else {
          ElMessage.error(res.data?.message || '登录失败')
        }
      } catch (e: any) {
        ElMessage.error(e.response?.data?.message || '登录失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped lang="scss">
.portal-login {
  min-height: calc(100vh - 160px);
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f8fafc 0%, #e0f2fe 100%);
  padding: 40px 20px;

  .login-container {
    width: 100%;
    max-width: 420px;
  }

  .login-card {
    background: white;
    border-radius: 16px;
    padding: 40px;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);

    .card-header {
      text-align: center;
      margin-bottom: 32px;

      h2 {
        font-size: 24px;
        font-weight: 700;
        color: #1e293b;
        margin: 0 0 8px 0;
      }

      p {
        color: #64748b;
        margin: 0;
      }
    }

    .form-options {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;

      a {
        color: #0ea5e9;
        text-decoration: none;
        font-size: 14px;
      }
    }

    .divider {
      display: flex;
      align-items: center;
      margin: 24px 0;

      &::before, &::after {
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
      gap: 12px;
    }

    .footer-link {
      text-align: center;
      margin-top: 24px;
      color: #64748b;

      a {
        color: #0ea5e9;
        text-decoration: none;
        font-weight: 500;
      }
    }

    .hr-login {
      text-align: center;
      margin-top: 16px;
      padding-top: 16px;
      border-top: 1px solid #f1f5f9;

      a {
        color: #94a3b8;
        text-decoration: none;
        font-size: 14px;

        &:hover { color: #0ea5e9; }
      }
    }
  }
}
</style>
