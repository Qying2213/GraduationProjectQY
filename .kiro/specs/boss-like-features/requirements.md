# Requirements Document

## Introduction

本文档定义了智能人才运营平台的核心招聘功能需求，类似于 Boss 直聘的基础功能模块。系统分为求职者端和企业端/HR端两大模块，并包含即时通讯功能实现双方的实时沟通。

## Glossary

- **System**: 智能人才运营平台系统
- **Candidate**: 求职者，使用平台寻找工作机会的用户
- **HR**: 企业招聘人员，负责发布职位和筛选候选人
- **Job**: 职位信息，包含职位名称、要求、薪资等
- **Application**: 求职申请，求职者投递简历到特定职位的记录
- **Resume**: 简历，求职者的个人信息和工作经历文档
- **Interview**: 面试，HR安排的候选人面试活动
- **Conversation**: 会话，求职者与HR之间的聊天会话
- **Chat_Message**: 聊天消息，会话中的单条消息

## Requirements

### Requirement 1: 职位浏览与搜索

**User Story:** As a Candidate, I want to browse and search job listings, so that I can find suitable job opportunities.

#### Acceptance Criteria

1. WHEN a Candidate visits the job list page, THE System SHALL display a paginated list of open jobs with title, company, location, salary, and skills
2. WHEN a Candidate enters a keyword in the search box, THE System SHALL filter jobs by title or description containing the keyword
3. WHEN a Candidate selects location filter, THE System SHALL display only jobs matching the selected location
4. WHEN a Candidate selects experience level filter, THE System SHALL display only jobs matching the experience requirement
5. WHEN a Candidate clicks on a job item, THE System SHALL navigate to the job detail page showing full job information
6. THE System SHALL display the number of applicants for each job in the list

### Requirement 2: 职位投递

**User Story:** As a Candidate, I want to apply for jobs with my resume, so that I can be considered for positions.

#### Acceptance Criteria

1. WHEN a Candidate clicks the apply button on a job detail page, THE System SHALL display an application dialog
2. WHEN a Candidate submits an application with a selected resume, THE System SHALL create an application record with status "pending"
3. WHEN a Candidate submits an application, THE System SHALL send a notification to the HR who created the job
4. IF a Candidate has already applied to a job, THEN THE System SHALL prevent duplicate applications and display a message
5. WHEN a Candidate provides a cover letter, THE System SHALL store it with the application record
6. THE System SHALL validate that the Candidate has at least one resume before allowing application submission

### Requirement 3: 申请状态追踪

**User Story:** As a Candidate, I want to view my application status and progress, so that I can track my job search activities.

#### Acceptance Criteria

1. WHEN a Candidate visits the my-applications page, THE System SHALL display all applications with job title, company, status, and application date
2. WHEN an application status changes, THE System SHALL update the display in real-time
3. THE System SHALL display application status as one of: pending, viewed, interview, offer, rejected
4. WHEN a Candidate filters by status, THE System SHALL display only applications matching the selected status
5. WHEN a Candidate clicks withdraw on an application, THE System SHALL delete the application record and confirm the action
6. THE System SHALL display a timeline showing status change history for each application

### Requirement 4: 在线简历管理

**User Story:** As a Candidate, I want to create and manage my online resume, so that I can present my qualifications to employers.

#### Acceptance Criteria

1. WHEN a Candidate visits the my-resume page, THE System SHALL display the current resume information including basic info, work experience, education, and skills
2. WHEN a Candidate clicks edit, THE System SHALL display an editable form for resume information
3. WHEN a Candidate saves resume changes, THE System SHALL persist the updated information to the database
4. WHEN a Candidate uploads a resume file, THE System SHALL store the file and add it to the attachments list
5. WHEN a Candidate deletes an attachment, THE System SHALL remove the file and update the attachments list
6. THE System SHALL validate required fields (name, phone, email) before saving resume changes
7. THE System SHALL support resume file formats: PDF, DOC, DOCX with maximum size of 10MB

### Requirement 5: 职位发布与管理

**User Story:** As an HR, I want to create and manage job postings, so that I can attract qualified candidates.

#### Acceptance Criteria

1. WHEN an HR clicks create job, THE System SHALL display a job creation form with fields for title, description, requirements, salary, location, type, and skills
2. WHEN an HR submits a valid job form, THE System SHALL create a new job record with status "open"
3. WHEN an HR edits an existing job, THE System SHALL update the job record with the new information
4. WHEN an HR changes job status to "closed", THE System SHALL stop displaying the job in candidate search results
5. WHEN an HR deletes a job, THE System SHALL soft-delete the job record and associated data
6. THE System SHALL validate required fields (title, description, location) before saving job

### Requirement 6: 候选人筛选与管理

**User Story:** As an HR, I want to view and filter candidates who applied to my jobs, so that I can identify suitable candidates.

#### Acceptance Criteria

1. WHEN an HR views a job's applications, THE System SHALL display a list of candidates with name, resume summary, application date, and status
2. WHEN an HR filters by application status, THE System SHALL display only applications matching the selected status
3. WHEN an HR clicks on a candidate, THE System SHALL display the candidate's full resume and application details
4. WHEN an HR updates an application status, THE System SHALL persist the change and notify the Candidate
5. WHEN an HR searches candidates by keyword, THE System SHALL filter by candidate name or skills
6. THE System SHALL display match score for each candidate if AI evaluation is available

### Requirement 7: 面试安排

**User Story:** As an HR, I want to schedule interviews with candidates, so that I can evaluate them for positions.

#### Acceptance Criteria

1. WHEN an HR clicks schedule interview for a candidate, THE System SHALL display an interview scheduling form
2. WHEN an HR submits interview details, THE System SHALL create an interview record with status "scheduled"
3. WHEN an interview is scheduled, THE System SHALL send a notification to the Candidate with interview details
4. WHEN an HR updates interview status, THE System SHALL persist the change and notify the Candidate
5. THE System SHALL validate that interview date is in the future before scheduling
6. WHEN an HR provides interview feedback, THE System SHALL store the feedback with rating and comments

### Requirement 8: 即时聊天通讯

**User Story:** As a Candidate or HR, I want to chat in real-time, so that I can communicate about job opportunities efficiently.

#### Acceptance Criteria

1. WHEN a user opens a conversation, THE System SHALL establish a WebSocket connection for real-time messaging
2. WHEN a user sends a message, THE System SHALL deliver it to the recipient in real-time via WebSocket
3. WHEN a user receives a message, THE System SHALL display it immediately in the conversation view
4. IF the recipient is offline, THEN THE System SHALL store the message and deliver it when the recipient connects
5. THE System SHALL persist all chat messages to the database for history retrieval
6. WHEN a user scrolls up in a conversation, THE System SHALL load older messages with pagination

### Requirement 9: 会话管理

**User Story:** As a user, I want to manage my chat conversations, so that I can organize my communications.

#### Acceptance Criteria

1. WHEN a user visits the messages page, THE System SHALL display a list of conversations sorted by last message time
2. THE System SHALL display the last message preview and timestamp for each conversation
3. THE System SHALL display unread message count badge for conversations with unread messages
4. WHEN a user opens a conversation, THE System SHALL mark all messages in that conversation as read
5. WHEN a user searches conversations, THE System SHALL filter by participant name or message content
6. THE System SHALL display the online status indicator for conversation participants

### Requirement 10: 消息通知

**User Story:** As a user, I want to receive notifications for new messages, so that I can respond promptly.

#### Acceptance Criteria

1. WHEN a new message arrives, THE System SHALL display a notification badge in the header
2. THE System SHALL update the unread count in real-time without page refresh
3. WHEN a user clicks the notification badge, THE System SHALL navigate to the messages page
4. IF the user is on a different page, THEN THE System SHALL show a toast notification for new messages
5. THE System SHALL play a notification sound for new messages when the browser tab is not focused
6. WHEN a user marks messages as read, THE System SHALL update the unread count immediately
