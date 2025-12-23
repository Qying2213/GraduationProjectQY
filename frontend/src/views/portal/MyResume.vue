<template>
  <div class="my-resume">
    <div class="page-container">
      <div class="page-header">
        <h1>我的简历</h1>
        <el-button type="primary" @click="showUploadDialog = true">
          <el-icon><Upload /></el-icon> 上传简历
        </el-button>
      </div>

      <el-row :gutter="24">
        <!-- 在线简历 -->
        <el-col :xs="24" :lg="16">
          <div class="resume-card">
            <div class="card-header">
              <h2>在线简历</h2>
              <el-button type="primary" link @click="editResume">编辑</el-button>
            </div>

            <!-- 基本信息 -->
            <div class="resume-section">
              <div class="section-header">
                <el-icon><User /></el-icon>
                <h3>基本信息</h3>
              </div>
              <div class="info-grid">
                <div class="info-item">
                  <label>姓名</label>
                  <span>{{ resume.name }}</span>
                </div>
                <div class="info-item">
                  <label>性别</label>
                  <span>{{ resume.gender }}</span>
                </div>
                <div class="info-item">
                  <label>年龄</label>
                  <span>{{ resume.age }}岁</span>
                </div>
                <div class="info-item">
                  <label>手机</label>
                  <span>{{ resume.phone }}</span>
                </div>
                <div class="info-item">
                  <label>邮箱</label>
                  <span>{{ resume.email }}</span>
                </div>
                <div class="info-item">
                  <label>现居地</label>
                  <span>{{ resume.location }}</span>
                </div>
              </div>
            </div>

            <!-- 求职意向 -->
            <div class="resume-section">
              <div class="section-header">
                <el-icon><Aim /></el-icon>
                <h3>求职意向</h3>
              </div>
              <div class="info-grid">
                <div class="info-item">
                  <label>期望职位</label>
                  <span>{{ resume.expectPosition }}</span>
                </div>
                <div class="info-item">
                  <label>期望城市</label>
                  <span>{{ resume.expectCity }}</span>
                </div>
                <div class="info-item">
                  <label>期望薪资</label>
                  <span>{{ resume.expectSalary }}</span>
                </div>
                <div class="info-item">
                  <label>到岗时间</label>
                  <span>{{ resume.availableTime }}</span>
                </div>
              </div>
            </div>

            <!-- 工作经历 -->
            <div class="resume-section">
              <div class="section-header">
                <el-icon><Suitcase /></el-icon>
                <h3>工作经历</h3>
              </div>
              <div class="experience-item" v-for="exp in resume.workExperience" :key="exp.company">
                <div class="exp-header">
                  <div class="exp-title">
                    <h4>{{ exp.company }}</h4>
                    <span class="exp-position">{{ exp.position }}</span>
                  </div>
                  <span class="exp-time">{{ exp.startTime }} - {{ exp.endTime }}</span>
                </div>
                <p class="exp-desc">{{ exp.description }}</p>
              </div>
            </div>

            <!-- 教育经历 -->
            <div class="resume-section">
              <div class="section-header">
                <el-icon><School /></el-icon>
                <h3>教育经历</h3>
              </div>
              <div class="experience-item" v-for="edu in resume.education" :key="edu.school">
                <div class="exp-header">
                  <div class="exp-title">
                    <h4>{{ edu.school }}</h4>
                    <span class="exp-position">{{ edu.major }} · {{ edu.degree }}</span>
                  </div>
                  <span class="exp-time">{{ edu.startTime }} - {{ edu.endTime }}</span>
                </div>
              </div>
            </div>

            <!-- 技能特长 -->
            <div class="resume-section">
              <div class="section-header">
                <el-icon><Medal /></el-icon>
                <h3>技能特长</h3>
              </div>
              <div class="skills-list">
                <el-tag v-for="skill in resume.skills" :key="skill" size="large">{{ skill }}</el-tag>
              </div>
            </div>
          </div>
        </el-col>

        <!-- 附件简历 -->
        <el-col :xs="24" :lg="8">
          <div class="attachment-card">
            <h3>附件简历</h3>
            <div class="attachment-list">
              <div class="attachment-item" v-for="file in attachments" :key="file.id">
                <div class="file-icon">
                  <el-icon :size="24"><Document /></el-icon>
                </div>
                <div class="file-info">
                  <span class="file-name">{{ file.name }}</span>
                  <span class="file-meta">{{ file.size }} · {{ file.uploadTime }}</span>
                </div>
                <div class="file-actions">
                  <el-button link type="primary" size="small">预览</el-button>
                  <el-button link type="danger" size="small" @click="deleteAttachment(file.id)">删除</el-button>
                </div>
              </div>
              <el-empty v-if="attachments.length === 0" description="暂无附件简历" :image-size="80" />
            </div>
          </div>

          <div class="tips-card">
            <h3>简历优化建议</h3>
            <ul>
              <li><el-icon><CircleCheck /></el-icon>完善基本信息</li>
              <li><el-icon><CircleCheck /></el-icon>添加工作经历</li>
              <li><el-icon><Warning /></el-icon>建议上传附件简历</li>
              <li><el-icon><Warning /></el-icon>添加项目经验</li>
            </ul>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 上传弹窗 -->
    <el-dialog v-model="showUploadDialog" title="上传简历" width="500px">
      <el-upload drag :auto-upload="false" accept=".pdf,.doc,.docx" :limit="1">
        <el-icon class="el-icon--upload" :size="48"><UploadFilled /></el-icon>
        <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip">支持 PDF、DOC、DOCX 格式，文件大小不超过 10MB</div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="showUploadDialog = false">取消</el-button>
        <el-button type="primary">上传</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, User, Aim, Suitcase, School, Medal, Document, CircleCheck, Warning, UploadFilled } from '@element-plus/icons-vue'

const showUploadDialog = ref(false)

const resume = ref({
  name: '张三',
  gender: '男',
  age: 28,
  phone: '138****1234',
  email: 'zhangsan@email.com',
  location: '北京',
  expectPosition: '高级前端工程师',
  expectCity: '北京、上海',
  expectSalary: '30-45K',
  availableTime: '随时到岗',
  workExperience: [
    {
      company: '某科技有限公司',
      position: '前端开发工程师',
      startTime: '2021-06',
      endTime: '至今',
      description: '负责公司核心产品的前端开发工作，使用Vue3+TypeScript技术栈，参与架构设计和性能优化。'
    },
    {
      company: '某互联网公司',
      position: '前端开发',
      startTime: '2019-07',
      endTime: '2021-05',
      description: '参与多个项目的前端开发，熟悉React和Vue框架。'
    }
  ],
  education: [
    {
      school: '某大学',
      major: '计算机科学与技术',
      degree: '本科',
      startTime: '2015-09',
      endTime: '2019-06'
    }
  ],
  skills: ['Vue3', 'TypeScript', 'React', 'Node.js', 'Webpack', 'Git']
})

const attachments = ref([
  { id: 1, name: '张三_前端工程师_简历.pdf', size: '256KB', uploadTime: '2024-01-10' }
])

const editResume = () => {
  ElMessage.info('编辑功能开发中')
}

const deleteAttachment = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这份简历吗？', '删除确认', { type: 'warning' })
    attachments.value = attachments.value.filter(a => a.id !== id)
    ElMessage.success('已删除')
  } catch {}
}
</script>

<style scoped lang="scss">
.my-resume {
  padding: 24px;
  background: #f8fafc;
  min-height: calc(100vh - 160px);

  .page-container {
    max-width: 1200px;
    margin: 0 auto;
  }

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    h1 {
      font-size: 24px;
      font-weight: 700;
      color: #1e293b;
      margin: 0;
    }
  }

  .resume-card {
    background: white;
    border-radius: 12px;
    padding: 24px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 24px;
      padding-bottom: 16px;
      border-bottom: 1px solid #f1f5f9;

      h2 {
        font-size: 20px;
        font-weight: 600;
        margin: 0;
      }
    }

    .resume-section {
      margin-bottom: 32px;

      &:last-child { margin-bottom: 0; }

      .section-header {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 16px;

        .el-icon { color: #0ea5e9; }

        h3 {
          font-size: 16px;
          font-weight: 600;
          margin: 0;
        }
      }

      .info-grid {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 16px;

        .info-item {
          label {
            display: block;
            font-size: 12px;
            color: #94a3b8;
            margin-bottom: 4px;
          }

          span {
            color: #1e293b;
          }
        }
      }

      .experience-item {
        padding: 16px;
        background: #f8fafc;
        border-radius: 8px;
        margin-bottom: 12px;

        &:last-child { margin-bottom: 0; }

        .exp-header {
          display: flex;
          justify-content: space-between;
          margin-bottom: 8px;

          h4 {
            font-size: 15px;
            font-weight: 600;
            margin: 0 0 4px 0;
          }

          .exp-position {
            color: #64748b;
            font-size: 14px;
          }

          .exp-time {
            color: #94a3b8;
            font-size: 13px;
          }
        }

        .exp-desc {
          color: #475569;
          font-size: 14px;
          line-height: 1.6;
          margin: 0;
        }
      }

      .skills-list {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
      }
    }
  }

  .attachment-card, .tips-card {
    background: white;
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 16px;

    h3 {
      font-size: 16px;
      font-weight: 600;
      margin: 0 0 16px 0;
    }
  }

  .attachment-list {
    .attachment-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px;
      background: #f8fafc;
      border-radius: 8px;
      margin-bottom: 8px;

      .file-icon {
        width: 40px;
        height: 40px;
        background: #fee2e2;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #ef4444;
      }

      .file-info {
        flex: 1;

        .file-name {
          display: block;
          font-weight: 500;
          color: #1e293b;
        }

        .file-meta {
          font-size: 12px;
          color: #94a3b8;
        }
      }
    }
  }

  .tips-card {
    ul {
      list-style: none;
      padding: 0;
      margin: 0;

      li {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 8px 0;
        color: #475569;

        .el-icon {
          &:first-child { color: #10b981; }
        }
      }

      li:has(.el-icon:first-child[class*="Warning"]) .el-icon {
        color: #f59e0b;
      }
    }
  }
}

@media (max-width: 768px) {
  .my-resume .resume-card .resume-section .info-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
