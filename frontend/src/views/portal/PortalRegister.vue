<template>
  <div class="portal-register">
    <div class="register-container">
      <div class="register-card">
        <div class="card-header">
          <h2>求职者注册</h2>
          <p>创建账号，开启求职之旅</p>
        </div>

        <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleRegister">
          <el-form-item prop="name">
            <el-input v-model="form.name" placeholder="姓名" size="large" :prefix-icon="User" />
          </el-form-item>
          <el-form-item prop="email">
            <el-input v-model="form.email" placeholder="邮箱" size="large" :prefix-icon="Message" />
          </el-form-item>
          <el-form-item prop="phone">
            <el-input v-model="form.phone" placeholder="手机号" size="large" :prefix-icon="Phone" />
          </el-form-item>
          <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="密码" size="large"
                      :prefix-icon="Lock" show-password />
          </el-form-item>
          <el-form-item prop="confirmPassword">
            <el-input v-model="form.confirmPassword" type="password" placeholder="确认密码" size="large"
                      :prefix-icon="Lock" show-password />
          </el-form-item>
          <el-form-item prop="agree">
            <el-checkbox v-model="form.agree">
              我已阅读并同意 <a href="#">用户协议</a> 和 <a href="#">隐私政策</a>
            </el-checkbox>
          </el-form-item>
          <el-button type="primary" size="large" :loading="loading" @click="handleRegister" style="width: 100%">
            注册
          </el-button>
        </el-form>

        <div class="footer-link">
          已有账号？<router-link to="/portal/login">立即登录</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { User, Message, Phone, Lock } from '@element-plus/icons-vue'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  name: '',
  email: '',
  phone: '',
  password: '',
  confirmPassword: '',
  agree: false
})

const validateConfirmPassword = (rule: any, value: string, callback: any) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  agree: [
    { validator: (rule, value, callback) => value ? callback() : callback(new Error('请同意用户协议')), trigger: 'change' }
  ]
}

const handleRegister = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await new Promise(r => setTimeout(r, 1000))
        ElMessage.success('注册成功，请登录')
        router.push('/portal/login')
      } catch (e: any) {
        ElMessage.error(e.message || '注册失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped lang="scss">
.portal-register {
  min-height: calc(100vh - 160px);
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f8fafc 0%, #e0f2fe 100%);
  padding: 40px 20px;

  .register-container {
    width: 100%;
    max-width: 420px;
  }

  .register-card {
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

    :deep(.el-checkbox__label) {
      a {
        color: #0ea5e9;
        text-decoration: none;
      }
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
  }
}
</style>
