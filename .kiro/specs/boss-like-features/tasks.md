# Implementation Plan: Boss-Like Features

## Overview

本实现计划基于现有项目代码进行增量开发，确保每个功能完整可用。任务按照依赖关系排序，每个任务都包含完整的实现和验证步骤。

## Tasks

- [x] 1. 数据库扩展 - 聊天功能表结构
  - [x] 1.1 创建 conversations 和 chat_messages 表
    - 在 backend/database/schema.sql 中添加表结构
    - 添加必要的索引
    - _Requirements: 8.5, 9.1_
  
  - [x] 1.2 创建数据库迁移脚本
    - 创建 migrations 目录和迁移文件
    - 确保可以在现有数据库上安全执行
    - _Requirements: 8.5_

- [x] 2. Checkpoint - 验证数据库结构
  - 执行迁移脚本，确保表创建成功
  - 验证索引正确创建

- [x] 3. 后端 - Application Service 增强
  - [x] 3.1 增强申请创建逻辑
    - 在 job-service 中添加重复申请检查
    - 申请创建后自动发送通知给 HR
    - 验证求职者是否有简历
    - _Requirements: 2.2, 2.3, 2.4, 2.6_
  
  - [x] 3.2 增强申请状态更新逻辑
    - 状态更新时发送通知给求职者
    - 记录状态变更历史
    - _Requirements: 3.2, 6.4_
  
  - [x] 3.3 添加获取职位申请列表接口
    - GET /jobs/:id/applications 接口
    - 支持状态筛选和分页
    - _Requirements: 6.1, 6.2_
  
  - [ ]* 3.4 编写 Application Service 单元测试
    - 测试重复申请检查
    - 测试状态更新通知
    - **Property 3: Duplicate Application Prevention**
    - **Validates: Requirements 2.4**

- [x] 4. Checkpoint - 验证申请功能
  - 测试申请创建、状态更新、通知发送
  - 确保所有接口正常工作

- [x] 5. 后端 - Chat Service 实现
  - [x] 5.1 创建聊天数据模型
    - 在 message-service/models 中添加 Conversation 和 ChatMessage 模型
    - 添加 GORM 关联关系
    - _Requirements: 8.5, 9.1_
  
  - [x] 5.2 实现会话管理接口
    - GET /conversations - 获取会话列表
    - POST /conversations - 创建/获取会话
    - 包含未读数统计和最后消息
    - _Requirements: 9.1, 9.2, 9.3_
  
  - [x] 5.3 实现聊天消息接口
    - GET /conversations/:id/messages - 获取消息列表（分页）
    - POST /conversations/:id/messages - 发送消息
    - PUT /conversations/:id/read - 标记已读
    - _Requirements: 8.5, 8.6, 9.4_
  
  - [x] 5.4 扩展 WebSocket Hub 支持聊天
    - 添加聊天消息类型处理
    - 实现实时消息推送
    - 实现在线状态管理
    - _Requirements: 8.2, 8.4, 9.6_
  
  - [ ]* 5.5 编写 Chat Service 单元测试
    - 测试消息发送和接收
    - 测试会话排序
    - **Property 15: Chat Message Persistence Round-Trip**
    - **Property 17: Conversation Sorting Order**
    - **Validates: Requirements 8.5, 9.1**

- [x] 6. Checkpoint - 验证聊天后端功能
  - 使用 Postman/curl 测试所有聊天接口
  - 测试 WebSocket 消息推送

- [x] 7. 后端 - Resume Service 在线编辑
  - [x] 7.1 添加在线简历数据模型
    - 扩展 talents 表或创建 online_resumes 表
    - 定义简历数据结构（基本信息、工作经历、教育经历、技能）
    - _Requirements: 4.1, 4.3_
  
  - [x] 7.2 实现在线简历接口
    - GET /resumes/online - 获取当前用户在线简历
    - PUT /resumes/online - 保存在线简历
    - 包含字段验证
    - _Requirements: 4.3, 4.6_
  
  - [ ]* 7.3 编写 Resume Service 单元测试
    - 测试简历保存和读取
    - 测试字段验证
    - **Property 7: Resume Data Persistence Round-Trip**
    - **Validates: Requirements 4.3**

- [x] 8. Checkpoint - 验证简历编辑功能
  - 测试简历保存和读取
  - 验证数据一致性

- [x] 9. 前端 - 申请功能增强
  - [x] 9.1 完善 PortalJobList 投递功能
    - 修复投递弹窗逻辑
    - 添加简历选择（从用户简历列表）
    - 添加投递前验证（是否有简历、是否已投递）
    - _Requirements: 2.1, 2.4, 2.6_
  
  - [x] 9.2 完善 PortalJobDetail 投递功能
    - 添加投递按钮和弹窗
    - 显示已投递状态
    - 添加求职信输入
    - _Requirements: 2.1, 2.5_
  
  - [x] 9.3 完善 MyApplications 状态追踪
    - 从后端获取真实申请数据
    - 显示状态时间线
    - 实现状态筛选
    - 实现撤回功能
    - _Requirements: 3.1, 3.3, 3.4, 3.5, 3.6_
  
  - [ ]* 9.4 编写申请功能前端测试
    - 测试投递流程
    - 测试状态筛选
    - **Property 5: Application Status Filter Correctness**
    - **Validates: Requirements 3.4**

- [x] 10. Checkpoint - 验证前端申请功能
  - 完整测试投递流程
  - 测试状态追踪和筛选

- [x] 11. 前端 - 在线简历编辑
  - [x] 11.1 创建简历编辑表单组件
    - 基本信息表单（姓名、电话、邮箱、地址）
    - 工作经历表单（可添加多条）
    - 教育经历表单（可添加多条）
    - 技能标签输入
    - _Requirements: 4.1, 4.2_
  
  - [x] 11.2 完善 MyResume 页面
    - 集成编辑表单组件
    - 实现保存功能
    - 添加表单验证
    - 显示保存成功/失败提示
    - _Requirements: 4.2, 4.3, 4.6_
  
  - [ ]* 11.3 编写简历编辑前端测试
    - 测试表单验证
    - 测试保存功能

- [x] 12. Checkpoint - 验证简历编辑功能
  - 完整测试简历编辑流程
  - 验证数据保存正确

- [x] 13. 前端 - 聊天功能实现
  - [x] 13.1 创建聊天 API 模块
    - 创建 frontend/src/api/chat.ts
    - 实现所有聊天相关 API 调用
    - _Requirements: 8.2, 9.1_
  
  - [x] 13.2 创建聊天组件
    - ConversationList.vue - 会话列表组件
    - ChatWindow.vue - 聊天窗口组件
    - MessageItem.vue - 消息项组件
    - ChatInput.vue - 输入框组件
    - _Requirements: 8.3, 9.1, 9.2_
  
  - [x] 13.3 创建求职者端聊天页面
    - 创建 PortalChat.vue 页面
    - 添加路由配置
    - 集成 WebSocket 实时消息
    - _Requirements: 8.1, 8.2, 8.3_
  
  - [x] 13.4 创建 HR 端聊天页面
    - 创建 ChatCenter.vue 页面
    - 添加路由配置
    - 集成 WebSocket 实时消息
    - _Requirements: 8.1, 8.2, 8.3_
  
  - [x] 13.5 实现未读消息提醒
    - 在 Header 组件添加未读数 badge
    - 实现实时更新
    - 添加消息通知
    - _Requirements: 10.1, 10.2, 10.6_
  
  - [ ]* 13.6 编写聊天功能前端测试
    - 测试消息发送和接收
    - 测试未读数更新
    - **Property 18: Unread Count Accuracy**
    - **Validates: Requirements 9.3**

- [x] 14. Checkpoint - 验证聊天功能
  - 完整测试聊天流程
  - 测试实时消息推送
  - 测试未读数更新

- [x] 15. HR 端功能增强
  - [x] 15.1 增强职位详情页申请管理
    - 在 JobDetail 页面添加申请列表 Tab
    - 显示申请人列表和状态
    - 实现状态更新功能
    - _Requirements: 6.1, 6.4_
  
  - [x] 15.2 增强候选人筛选功能
    - 在 TalentList 添加关键词搜索
    - 添加技能筛选
    - 显示匹配分数
    - _Requirements: 6.5, 6.6_
  
  - [x] 15.3 完善面试安排功能
    - 从申请创建面试
    - 面试安排后发送通知
    - 日期验证（必须是未来时间）
    - _Requirements: 7.1, 7.2, 7.3, 7.5_
  
  - [ ]* 15.4 编写 HR 功能测试
    - 测试申请状态更新
    - 测试面试安排
    - **Property 14: Interview Date Future Validation**
    - **Validates: Requirements 7.5**

- [x] 16. Checkpoint - 验证 HR 功能
  - 测试申请管理流程
  - 测试面试安排流程

- [x] 17. 集成测试和修复
  - [x] 17.1 端到端流程测试
    - 测试完整求职流程：浏览职位 → 投递 → 查看状态 → 聊天
    - 测试完整招聘流程：发布职位 → 查看申请 → 更新状态 → 安排面试 → 聊天
    - _Requirements: All_
  
  - [x] 17.2 修复发现的问题
    - 修复测试中发现的 bug
    - 优化用户体验
    - _Requirements: All_

- [x] 18. Final Checkpoint - 完整功能验证
  - 确保所有功能正常工作
  - 确保没有明显 bug
  - 确保用户体验流畅

## Notes

- 任务按依赖关系排序，后端功能先于前端功能
- 每个 Checkpoint 用于验证阶段性成果，确保功能完整可用
- 标记 `*` 的测试任务为可选，但建议执行以确保代码质量
- 所有功能都基于现有代码增量开发，复用现有组件和样式
