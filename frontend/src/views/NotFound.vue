<template>
  <div class="not-found-page">
    <div class="content">
      <div class="error-illustration">
        <div class="number">4</div>
        <div class="circle">
          <div class="face">
            <div class="eye left"></div>
            <div class="eye right"></div>
            <div class="mouth"></div>
          </div>
        </div>
        <div class="number">4</div>
      </div>

      <h1 class="title">页面走丢了</h1>
      <p class="description">抱歉，您访问的页面不存在或已被移除</p>

      <div class="actions">
        <el-button type="primary" size="large" @click="goHome">
          <el-icon><HomeFilled /></el-icon>
          返回首页
        </el-button>
        <el-button size="large" @click="goBack">
          <el-icon><Back /></el-icon>
          返回上页
        </el-button>
      </div>

      <div class="suggestions">
        <h3>您可以尝试</h3>
        <div class="suggestion-list">
          <div class="suggestion-item" @click="$router.push('/dashboard')">
            <el-icon><DataAnalysis /></el-icon>
            <span>查看仪表板</span>
          </div>
          <div class="suggestion-item" @click="$router.push('/talents')">
            <el-icon><User /></el-icon>
            <span>浏览人才库</span>
          </div>
          <div class="suggestion-item" @click="$router.push('/jobs')">
            <el-icon><Suitcase /></el-icon>
            <span>查看职位</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { HomeFilled, Back, DataAnalysis, User, Suitcase } from '@element-plus/icons-vue'

const router = useRouter()

const goHome = () => {
  router.push('/dashboard')
}

const goBack = () => {
  router.back()
}
</script>

<style scoped lang="scss">
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8f0 100%);
  position: relative;
  overflow: hidden;
}

.content {
  text-align: center;
  z-index: 1;
  padding: 40px;
}

.error-illustration {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 40px;

  .number {
    font-size: 150px;
    font-weight: 800;
    color: #667eea;
    line-height: 1;
    text-shadow: 4px 4px 0 rgba(102, 126, 234, 0.2);
  }

  .circle {
    width: 140px;
    height: 140px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    animation: float 3s ease-in-out infinite;
    box-shadow: 0 20px 40px rgba(102, 126, 234, 0.3);

    .face {
      position: relative;
      width: 80px;
      height: 60px;

      .eye {
        position: absolute;
        width: 16px;
        height: 16px;
        background: white;
        border-radius: 50%;
        top: 10px;
        animation: blink 3s ease-in-out infinite;

        &.left { left: 10px; }
        &.right { right: 10px; }

        &::after {
          content: '';
          position: absolute;
          width: 8px;
          height: 8px;
          background: #1a1a2e;
          border-radius: 50%;
          top: 4px;
          left: 4px;
        }
      }

      .mouth {
        position: absolute;
        width: 30px;
        height: 15px;
        border: 4px solid white;
        border-top: none;
        border-radius: 0 0 30px 30px;
        bottom: 5px;
        left: 50%;
        transform: translateX(-50%);
      }
    }
  }
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-15px); }
}

@keyframes blink {
  0%, 90%, 100% { transform: scaleY(1); }
  95% { transform: scaleY(0.1); }
}

.title {
  font-size: 36px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 16px 0;
}

.description {
  font-size: 18px;
  color: #6b7280;
  margin: 0 0 40px 0;
}

.actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  margin-bottom: 60px;

  .el-button {
    padding: 12px 32px;
    border-radius: 12px;
    font-size: 16px;

    &--primary {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border: none;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
      }
    }
  }
}

.suggestions {
  h3 {
    font-size: 16px;
    color: #6b7280;
    margin: 0 0 20px 0;
    font-weight: 500;
  }

  .suggestion-list {
    display: flex;
    gap: 16px;
    justify-content: center;

    .suggestion-item {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 12px 24px;
      background: white;
      border-radius: 12px;
      cursor: pointer;
      transition: all 0.3s ease;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);

      &:hover {
        transform: translateY(-3px);
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
        color: #667eea;
      }

      .el-icon {
        font-size: 20px;
      }

      span {
        font-size: 14px;
        font-weight: 500;
      }
    }
  }
}

// 背景装饰
.bg-decoration {
  position: absolute;
  inset: 0;
  pointer-events: none;

  .shape {
    position: absolute;
    border-radius: 50%;
    opacity: 0.1;

    &.shape-1 {
      width: 400px;
      height: 400px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      top: -100px;
      right: -100px;
    }

    &.shape-2 {
      width: 300px;
      height: 300px;
      background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
      bottom: -50px;
      left: -50px;
    }

    &.shape-3 {
      width: 200px;
      height: 200px;
      background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
      top: 50%;
      left: 10%;
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .error-illustration {
    .number {
      font-size: 80px;
    }

    .circle {
      width: 80px;
      height: 80px;

      .face {
        width: 50px;
        height: 40px;

        .eye {
          width: 10px;
          height: 10px;

          &::after {
            width: 5px;
            height: 5px;
            top: 2px;
            left: 2px;
          }
        }

        .mouth {
          width: 20px;
          height: 10px;
          border-width: 3px;
        }
      }
    }
  }

  .title {
    font-size: 28px;
  }

  .description {
    font-size: 16px;
  }

  .actions {
    flex-direction: column;

    .el-button {
      width: 100%;
    }
  }

  .suggestions {
    .suggestion-list {
      flex-direction: column;
    }
  }
}
</style>
