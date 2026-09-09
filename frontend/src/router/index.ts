import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/RegisterView.vue'),
    },
    {
      path: '/user',
      name: 'user',
      component: () => import('@/views/UserView.vue'),
    },
    {
      path: '/token',
      name: 'token',
      component: () => import('@/views/TokenView.vue'),
    },
    {
      path: '/port',
      name: 'port',
      component: () => import('@/views/PortView.vue'),
    },
    {
      path: '/client',
      name: 'client',
      component: () => import('@/views/ClientView.vue'),
    },
    {
      path: '/proxy',
      name: 'proxy',
      component: () => import('@/views/ProxyView.vue'),
    },
    {
      path: '/config',
      name: 'config',
      component: () => import('@/views/ConfigView.vue'),
    },
    {
      path: '/admin/user',
      name: 'admin user',
      component: () => import('@/views/AdminUserView.vue'),
    },
    {
      path: '/admin/token',
      name: 'admin token',
      component: () => import('@/views/AdminTokenView.vue')
    },
    {
      path: '/admin/port',
      name: 'admin port',
      component: () => import('@/views/AdminPortView.vue'),
    },
    {
      path: '/admin/client',
      name: 'admin client',
      component: () => import('@/views/AdminClientView.vue'),
    },
    {
      path: '/admin/proxy',
      name: 'admin proxy',
      component: () => import('@/views/AdminProxyView.vue'),
    },
    {
      path: '/404',
      component: () => import('@/views/Error404.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/404'
    }
  ],
})

export default router
