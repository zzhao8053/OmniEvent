import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/desktop/LoginPage.vue')
  },
  {
    path: '/signup',
    name: 'Signup',
    component: () => import('@/views/desktop/SignupPage.vue')
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/desktop/ProfilePage.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else if ((to.path === '/login' || to.path === '/signup') && userStore.isLoggedIn) {
    next('/profile')
  } else {
    next()
  }
})

import { useUserStore } from '@/stores/user'
