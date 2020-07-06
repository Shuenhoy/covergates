import Vue from 'vue';
import VueRouter, { RouteConfig } from 'vue-router';
import { fetchReportSource, fetchCurrentRepository } from './fetchers';
import store, { Actions, State } from '@/store';

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: '/',
    name: 'Main',
    component: () => import('@/views/Main.vue'),
    children: [
      {
        path: '/',
        name: 'Home',
        component: () => import('@/views/Home.vue')
      },
      {
        path: '/repo',
        name: 'Repo',
        component: () => import('@/views/Repo.vue')
      }, {
        path: '/report/:scm/:namespace/:name',
        name: 'Report',
        component: () => import('@/views/Report.vue'),
        beforeEnter: fetchCurrentRepository(store),
        children: [
          {
            path: '',
            name: 'report-overview',
            component: () => import('@/components/ReportOverview.vue')
          },
          {
            path: 'code',
            name: 'report-code',
            component: () => import('@/components/ReportCode.vue')
          },
          {
            path: 'code/:path+',
            name: 'report-source',
            component: () => import('@/components/ReportSource.vue'),
            beforeEnter: fetchReportSource(store)
          }
        ]
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  }
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.NODE_ENV === 'production' ? VUE_BASE : process.env.BASE_URL,
  routes
});

export default router;
