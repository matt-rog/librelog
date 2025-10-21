// AI-assisted code
import { createRouter, createWebHashHistory } from 'vue-router'
import { api } from './api.js'
import Landing from './views/Landing.vue'
import App from './views/App.vue'
import Docs from './views/Docs.vue'

const routes = [
  { path: '/', component: Landing },
  { path: '/app', component: App, meta: { auth: true } },
  { path: '/docs', component: Docs },
  { path: '/docs/:page', component: Docs },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach((to) => {
  if (to.meta.auth && !api.getToken()) {
    return '/'
  }
})

export default router
