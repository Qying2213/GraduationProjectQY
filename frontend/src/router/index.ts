import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
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
            }
        ]
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
