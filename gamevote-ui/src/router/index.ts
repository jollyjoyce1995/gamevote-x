import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Cookies from 'js-cookie'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/DashboardView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/parties/new',
      name: 'new-party',
      component: () => import('@/views/NewPartyView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/parties/:code',
      name: 'party',
      component: () => import('@/views/PartyView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/parties/:code/drinks',
      name: 'drinks',
      component: () => import('@/views/DrinksView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

router.beforeEach(async (to) => {
  // Skip validation for login page
  if (to.name === 'login') {
    return
  }

  // For routes that require auth, validate the session
  if (to.meta.requiresAuth !== false) {
    const authStore = useAuthStore()
    const isValidSession = await authStore.validateSession()

    if (!isValidSession) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
  }

  // If already authenticated and trying to go to login, redirect to home
  const hasCookie = !!Cookies.get('username')
  if (to.name === 'login' && hasCookie) {
    return { name: 'home' }
  }
})

export default router
