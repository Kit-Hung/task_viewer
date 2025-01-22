import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import TaskDetail from '../views/TaskDetail.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/task/:id',
      name: 'taskDetail',
      component: TaskDetail,
      props: true
    }
  ]
})

export default router 