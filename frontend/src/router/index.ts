import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import { useUserStore } from "@/store/user";

// 路由表按“后台管理端”和“求职者前台端”两条线组织。
// requiresAuth 控制是否需要登录，permission 字段则给权限菜单/页面访问控制预留。
const routes: RouteRecordRaw[] = [
  // 后台管理登录
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/auth/Login.vue"),
    meta: { requiresAuth: false },
  },
  {
    path: "/register",
    name: "Register",
    component: () => import("@/views/auth/Register.vue"),
    meta: { requiresAuth: false },
  },
  {
    path: "/dev-logs",
    name: "DevLogs",
    component: () => import("@/views/system/DevLogsPage.vue"),
    meta: {
      title: "后台运行终端日志",
      requiresAuth: false,
      hideDevLogEntry: true,
    },
  },
  // 前台求职端
  {
    path: "/portal",
    component: () => import("@/components/layout/PortalLayout.vue"),
    meta: { requiresAuth: false },
    children: [
      {
        path: "",
        name: "PortalHome",
        component: () => import("@/views/portal/PortalHome.vue"),
        meta: { title: "首页" },
      },
      {
        path: "jobs",
        name: "PortalJobs",
        component: () => import("@/views/portal/PortalJobList.vue"),
        meta: { title: "职位列表" },
      },
      {
        path: "jobs/:id",
        name: "PortalJobDetail",
        component: () => import("@/views/portal/PortalJobDetail.vue"),
        meta: { title: "职位详情" },
      },
      {
        path: "companies",
        name: "PortalCompanies",
        component: () => import("@/views/portal/PortalCompanies.vue"),
        meta: { title: "企业招聘" },
      },
      {
        path: "login",
        name: "PortalLogin",
        component: () => import("@/views/portal/PortalLogin.vue"),
        meta: { title: "求职者登录" },
      },
      {
        path: "register",
        name: "PortalRegister",
        component: () => import("@/views/portal/PortalRegister.vue"),
        meta: { title: "求职者注册" },
      },
      {
        path: "my-applications",
        name: "MyApplications",
        component: () => import("@/views/portal/MyApplications.vue"),
        meta: { title: "我的投递", requiresAuth: true },
      },
      {
        path: "my-resume",
        name: "MyResume",
        component: () => import("@/views/portal/MyResume.vue"),
        meta: { title: "我的简历", requiresAuth: true },
      },
      {
        path: "chat",
        name: "PortalChat",
        component: () => import("@/views/portal/PortalChat.vue"),
        meta: { title: "消息中心", requiresAuth: true },
      },
    ],
  },
  // 数据大屏（全屏独立页面）
  {
    path: "/data-screen",
    name: "DataScreen",
    component: () => import("@/views/dashboard/DataScreen.vue"),
    meta: { title: "数据大屏", requiresAuth: true },
  },
  // 后台管理系统
  {
    path: "/",
    component: () => import("@/components/layout/MainLayout.vue"),
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        redirect: "/dashboard",
      },
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("@/views/dashboard/Dashboard.vue"),
        meta: { title: "仪表板" },
      },
      {
        path: "talents",
        name: "Talents",
        component: () => import("@/views/talents/TalentList.vue"),
        meta: { title: "人才管理" },
      },
      {
        path: "talents/:id",
        name: "TalentDetail",
        component: () => import("@/views/talents/TalentDetail.vue"),
        meta: { title: "人才详情" },
      },
      {
        path: "jobs",
        name: "Jobs",
        component: () => import("@/views/jobs/JobList.vue"),
        meta: { title: "职位管理" },
      },
      {
        path: "jobs/:id",
        name: "JobDetail",
        component: () => import("@/views/jobs/JobDetail.vue"),
        meta: { title: "职位详情" },
      },
      {
        path: "resumes",
        name: "Resumes",
        component: () => import("@/views/resumes/ResumeList.vue"),
        meta: { title: "简历管理" },
      },
      {
        path: "recommend",
        name: "Recommend",
        component: () => import("@/views/recommend/RecommendPage.vue"),
        meta: { title: "智能推荐" },
      },
      {
        path: "kanban",
        name: "Kanban",
        component: () => import("@/views/kanban/RecruitmentKanban.vue"),
        meta: { title: "招聘看板" },
      },
      {
        path: "calendar",
        name: "Calendar",
        component: () => import("@/views/calendar/InterviewCalendar.vue"),
        meta: { title: "面试日历" },
      },
      {
        path: "interviews/:id",
        name: "InterviewDetail",
        component: () => import("@/views/interviews/InterviewDetail.vue"),
        meta: { title: "面试详情" },
      },
      {
        path: "messages",
        name: "Messages",
        component: () => import("@/views/messages/MessageCenter.vue"),
        meta: { title: "消息中心" },
      },
      {
        path: "chat",
        name: "ChatCenter",
        component: () => import("@/views/messages/ChatCenter.vue"),
        meta: { title: "即时聊天" },
      },
      {
        path: "profile",
        name: "Profile",
        component: () => import("@/views/profile/UserProfile.vue"),
        meta: { title: "个人中心" },
      },
      {
        path: "roles",
        name: "RoleManagement",
        component: () => import("@/views/system/RoleManagement.vue"),
        meta: { title: "权限管理", permission: "role:view" },
      },
      {
        path: "reports",
        name: "Reports",
        component: () => import("@/views/reports/ReportsPage.vue"),
        meta: { title: "数据报表" },
      },
      {
        path: "settings",
        name: "Settings",
        component: () => import("@/views/system/SettingsPage.vue"),
        meta: { title: "系统设置" },
      },
      {
        path: "logs",
        name: "OperationLogs",
        component: () => import("@/views/system/OperationLogs.vue"),
        meta: { title: "操作日志", permission: "log:view" },
      },
      {
        path: "ai-evaluate",
        name: "AIEvaluate",
        component: () => import("@/views/ai/AIEvaluate.vue"),
        meta: { title: "AI智能评估" },
      },
      {
        path: "ai-process",
        name: "AIProcessFlow",
        component: () => import("@/views/ai/AIProcessFlow.vue"),
        meta: { title: "AI处理流程" },
      },
      {
        path: "evaluation-results",
        name: "EvaluationResults",
        component: () => import("@/views/ai/EvaluationResults.vue"),
        meta: { title: "评估结果" },
      },
      {
        path: "notices/:id/edit",
        name: "NoticeEdit",
        component: () => import("@/views/system/NoticeEdit.vue"),
        meta: { title: "编辑公告" },
      },
      {
        path: "notices",
        name: "Notices",
        component: () => import("@/views/system/NoticePage.vue"),
        meta: { title: "公告管理" },
      },
      {
        path: "notices/:id",
        name: "NoticeDetail",
        component: () => import("@/views/system/NoticeDetail.vue"),
        meta: { title: "公告详情" },
      },
    ],
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: () => import("@/views/NotFound.vue"),
    meta: { requiresAuth: false },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const userStore = useUserStore();

  // 需要登录但未登录时拦截到后台登录页。
  // 注意：求职端受保护页面当前也会先进入统一登录判断，具体登录页由业务页面和请求拦截器处理。
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next("/login");
    return;
  }

  // 已登录用户再次访问登录/注册页时，按角色分流到对应首页。
  if (
    (to.path === "/login" || to.path === "/register") &&
    userStore.isLoggedIn
  ) {
    // 求职者跳转到求职端，其他角色跳转到后台。
    if (userStore.role === "candidate") {
      next("/portal");
    } else {
      next("/dashboard");
    }
    return;
  }

  // 求职者不能访问后台管理页面，避免候选人角色进入 HR 工作台。
  if (userStore.isLoggedIn && userStore.role === "candidate") {
    const isPortalRoute = to.path.startsWith("/portal");
    const isPublicRoute = to.meta.requiresAuth === false;

    if (!isPortalRoute && !isPublicRoute && to.path !== "/") {
      next("/portal");
      return;
    }
  }

  next();
});

export default router;
