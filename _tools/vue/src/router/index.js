import Vue from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import(/* webpackChunkName: "home" */ '../views/Home.vue'),
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue'),
  },
  {
    path: '/viewer',
    name: 'Viewer',
    component: () => import(/* webpackChunkName: "icve" */ '../views/Viewer.vue'),
  },
  {
    path: '/icve',
    name: 'Icve',
    component: () => import(/* webpackChunkName: "icve" */ '../views/Viewer.vue'),
  },
  {
    path: '/mirador',
    name: 'Mirador',
    component: () => import(/* webpackChunkName: "mirador" */ '../views/Viewer.vue'),
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import(/* webpackChunkName: "search" */ '../views/Search.vue'),
  },
];

const router = new VueRouter({
  mode: 'history',
  // base: process.env.BASE_URL,
  base: '/jikei/',
  routes,
});

export default router;
