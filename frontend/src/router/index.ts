import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

import AdminDashboardPage from '../pages/AdminDashboardPage.vue'
import AdminLoginPage from '../pages/AdminLoginPage.vue'
import AuthPage from '../pages/AuthPage.vue'
import HomePage from '../pages/HomePage.vue'
import ProblemDetailPage from '../pages/ProblemDetailPage.vue'
import ProblemSetDetailPage from '../pages/ProblemSetDetailPage.vue'
import ProblemSetsPage from '../pages/ProblemSetsPage.vue'
import ProblemsPage from '../pages/ProblemsPage.vue'
import ProgressPage from '../pages/ProgressPage.vue'
import SubmissionsPage from '../pages/SubmissionsPage.vue'
import AppShell from '../shared/ui/AppShell.vue'
import { useAuthStore } from '../features/auth/store'
import { useAdminAuthStore } from '../features/admin-auth/store'

const routes: RouteRecordRaw[] = [
  { path: '/admin/login', component: AdminLoginPage },
  { path: '/admin', component: AdminDashboardPage },
  {
    path: '/',
    component: AppShell,
    children: [
      { path: '', component: HomePage },
      { path: 'login', component: AuthPage },
      { path: 'register', component: AuthPage },
      { path: 'problem-sets', component: ProblemSetsPage },
      { path: 'problem-sets/:slug', component: ProblemSetDetailPage },
      { path: 'problems', component: ProblemsPage },
      { path: 'problems/:slug', component: ProblemDetailPage },
      { path: 'submissions', component: SubmissionsPage },
      { path: 'progress', component: ProgressPage },
    ],
  },
]

export function createAppRouter() {
  const router = createRouter({
    history: createWebHistory(),
    routes,
    scrollBehavior() {
      return { top: 0 }
    },
  })

  // 导航守卫：检查登录状态
  router.beforeEach((to, _from, next) => {
    const authStore = useAuthStore()
    const adminAuthStore = useAdminAuthStore()
    const isProtected = to.path.startsWith('/problem-sets') ||
                        to.path.startsWith('/problems') ||
                        to.path.startsWith('/submissions') ||
                        to.path.startsWith('/progress')

    if (to.path.startsWith('/admin') && to.path !== '/admin/login' && !adminAuthStore.isAuthenticated) {
      next('/admin/login')
    } else if (to.path === '/admin/login' && adminAuthStore.isAuthenticated) {
      next('/admin')
    } else if (isProtected && !authStore.isAuthenticated) {
      next('/login')
    } else if ((to.path === '/login' || to.path === '/register') && authStore.isAuthenticated) {
      next('/problem-sets')
    } else {
      next()
    }
  })

  return router
}

export const router = createAppRouter()
