<template>
  <div class="job-list">
    <div class="page-header">
      <h1>职位管理</h1>
      <el-button type="primary">
        <el-icon><Plus /></el-icon>
        发布职位
      </el-button>
    </div>
    
    <div class="card">
      <el-table :data="[]" style="width: 100%">
        <el-table-column prop="title" label="职位名称" />
        <el-table-column prop="department" label="部门" />
        <el-table-column prop="location" label="地点" />
        <el-table-column prop="type" label="类型" />
        <el-table-column prop="status" label="状态" />
        <el-table-column label="操作">
          <template #default>
            <el-button link type="primary">查看</el-button>
            <el-button link type="primary">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
// Job list page
</script>

<style scoped lang="scss">
.job-list {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    
    h1 {
      font-size: 24px;
      font-weight: 700;
      margin: 0;
    }
  }
}
</style>
