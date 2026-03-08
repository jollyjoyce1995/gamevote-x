import { createRouter, createWebHistory } from 'vue-router'
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

router.beforeEach((to) => {
  const isAuthenticated = !!Cookies.get('username')
  if (to.meta.requiresAuth !== false && !isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && isAuthenticated) {
    return { name: 'home' }
  }
})

export default router
