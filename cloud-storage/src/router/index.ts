import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/login.vue'
import FileBrowser from '@/views/FileBrowser.vue'
import FileWindow from '@/views/FileWindow.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/files',
    },
    {
      path: '/files/:pathMatch(.*)*',
      name: 'files',
      component: FileBrowser,
      meta: {
        requiresAuth: true // This route requires authentication
      },
      children: [
        {
          path: '', // matches /files
          name: 'files-root',
          component: FileWindow, // your main file window component
        },
        {
          path: ':subPath(.*)*', // matches /files/anything
          name: 'files-sub',
          component: FileWindow, // or another component for subfolders
        }
      ]
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    }
  ],
})


router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const token = localStorage.getItem('jwt')
  if (requiresAuth && !token) {
      next({ name: 'login' }) // Change 'Login' to your login route name
  } else {
      next()
  }
})
export default router
