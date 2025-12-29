import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
    // 后台管理登录
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/auth/Login.vue'),
        meta: { requiresAuth: false }
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('@/views/auth/Register.vue'),
        meta: { requiresAuth: false }
    },
    // 前台求职端
    {
        path: '/portal',
        component: () => import('@/components/layout/PortalLayout.vue'),
        meta: { requiresAuth: false },
        children: [
            {
                path: '',
                name: 'PortalHome',
                component: () => import('@/views/portal/PortalHome.vue'),
                meta: { title: '首页' }
            },
            {
                path: 'jobs',
                name: 'PortalJobs',
                component: () => import('@/views/portal/PortalJobList.vue'),
                meta: { title: '职位列表' }
            },
            {
                path: 'jobs/:id',
                name: 'PortalJobDetail',
                component: () => import('@/views/portal/PortalJobDetail.vue'),
                meta: { title: '职位详情' }
            },
            {
                path: 'companies',
                name: 'PortalCompanies',
                component: () => import('@/views/portal/PortalCompanies.vue'),
                meta: { title: '企业招聘' }
            },
            {
                path: 'login',
                name: 'PortalLogin',
                component: () => import('@/views/portal/PortalLogin.vue'),
                meta: { title: '求职者登录' }
            },
            {
                path: 'register',
                name: 'PortalRegister',
                component: () => import('@/views/portal/PortalRegister.vue'),
                meta: { title: '求职者注册' }
            },
            {
                path: 'my-applications',
                name: 'MyApplications',
                component: () => import('@/views/portal/MyApplications.vue'),
                meta: { title: '我的投递', requiresAuth: true }
            },
            {
                path: 'my-resume',
                name: 'MyResume',
                component: () => import('@/views/portal/MyResume.vue'),
                meta: { title: '我的简历', requiresAuth: true }
            }
        ]
    },
    // 数据大屏（全屏独立页面）
    {
        path: '/data-screen',
        name: 'DataScreen',
        component: () => import('@/views/dashboard/DataScreen.vue'),
        meta: { title: '数据大屏', requiresAuth: true }
    },
    // 后台管理系统
    {
        path: '/',
        component: () => import('@/components/layout/MainLayout.vue'),
        meta: { requiresAuth: true },
        children: [
            {
                path: '',
                redirect: '/dashboard'
            },
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('@/views/dashboard/Dashboard.vue'),
                meta: { title: '仪表板' }
            },
            {
                path: 'talents',
                name: 'Talents',
                component: () => import('@/views/talents/TalentList.vue'),
                meta: { title: '人才管理' }
            },
            {
                path: 'talents/:id',
                name: 'TalentDetail',
                component: () => import('@/views/talents/TalentDetail.vue'),
                meta: { title: '人才详情' }
            },
            {
                path: 'jobs',
                name: 'Jobs',
                component: () => import('@/views/jobs/JobList.vue'),
                meta: { title: '职位管理' }
            },
            {
                path: 'jobs/:id',
                name: 'JobDetail',
                component: () => import('@/views/jobs/JobDetail.vue'),
                meta: { title: '职位详情' }
            },
            {
                path: 'resumes',
                name: 'Resumes',
                component: () => import('@/views/resumes/ResumeList.vue'),
                meta: { title: '简历管理' }
            },
            {
                path: 'recommend',
                name: 'Recommend',
                component: () => import('@/views/recommend/RecommendPage.vue'),
                meta: { title: '智能推荐' }
            },
            {
                path: 'kanban',
                name: 'Kanban',
                component: () => import('@/views/kanban/RecruitmentKanban.vue'),
                meta: { title: '招聘看板' }
            },
            {
                path: 'calendar',
                name: 'Calendar',
                component: () => import('@/views/calendar/InterviewCalendar.vue'),
                meta: { title: '面试日历' }
            },
            {
                path: 'interviews/:id',
                name: 'InterviewDetail',
                component: () => import('@/views/interviews/InterviewDetail.vue'),
                meta: { title: '面试详情' }
            },
            {
                path: 'messages',
                name: 'Messages',
                component: () => import('@/views/messages/MessageCenter.vue'),
                meta: { title: '消息中心' }
            },
            {
                path: 'profile',
                name: 'Profile',
                component: () => import('@/views/profile/UserProfile.vue'),
                meta: { title: '个人中心' }
            },
            {
                path: 'roles',
                name: 'RoleManagement',
                component: () => import('@/views/system/RoleManagement.vue'),
                meta: { title: '权限管理', permission: 'role:view' }
            },
            {
                path: 'reports',
                name: 'Reports',
                component: () => import('@/views/reports/ReportsPage.vue'),
                meta: { title: '数据报表' }
            },
            {
                path: 'settings',
                name: 'Settings',
                component: () => import('@/views/system/SettingsPage.vue'),
                meta: { title: '系统设置' }
            },
            {
                path: 'logs',
                name: 'OperationLogs',
                component: () => import('@/views/system/OperationLogs.vue'),
                meta: { title: '操作日志', permission: 'log:view' }
            }
        ]
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/NotFound.vue'),
        meta: { requiresAuth: false }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
    const userStore = useUserStore()

    if (to.meta.requiresAuth && !userStore.isLoggedIn) {
        next('/login')
    } else if ((to.path === '/login' || to.path === '/register') && userStore.isLoggedIn) {
        next('/dashboard')
    } else {
        next()
    }
})

export default router
