import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

import AuthPage from '../pages/AuthPage.vue'
import HomePage from '../pages/HomePage.vue'
import ProblemDetailPage from '../pages/ProblemDetailPage.vue'
import ProblemSetDetailPage from '../pages/ProblemSetDetailPage.vue'
import ProblemSetsPage from '../pages/ProblemSetsPage.vue'
import ProblemsPage from '../pages/ProblemsPage.vue'
import ProgressPage from '../pages/ProgressPage.vue'
import SubmissionsPage from '../pages/SubmissionsPage.vue'
import AppShell from '../shared/ui/AppShell.vue'

const routes: RouteRecordRaw[] = [
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
  return createRouter({
    history: createWebHistory(),
    routes,
    scrollBehavior() {
      return { top: 0 }
    },
  })
}

export const router = createAppRouter()
