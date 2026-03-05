import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import DocumentsView from '../views/DocumentsView.vue'
import ContactsView from '../views/ContactsView.vue'
import LocationsView from '../views/LocationsView.vue'
import DigitalAccessView from '../views/DigitalAccessView.vue'
import InsurancePoliciesView from '../views/InsurancePoliciesView.vue'
import ServiceAccountsView from '../views/ServiceAccountsView.vue'
import SurvivorLetterView from '../views/SurvivorLetterView.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
    meta: { requiresAuth: true },
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
  },
  {
    path: '/register',
    name: 'register',
    component: RegisterView,
  },
  {
    path: '/documents',
    name: 'documents',
    component: DocumentsView,
    meta: { requiresAuth: true },
  },
  {
    path: '/contacts',
    name: 'contacts',
    component: ContactsView,
    meta: { requiresAuth: true },
  },
  {
    path: '/locations',
    name: 'locations',
    component: LocationsView,
    meta: { requiresAuth: true },
  },
  {
    path: '/digital-access',
    name: 'digital-access',
    component: DigitalAccessView,
    meta: { requiresAuth: true },
  },
  {
    path: '/insurance-policies',
    name: 'insurance-policies',
    component: InsurancePoliciesView,
    meta: { requiresAuth: true },
  },
  {
    path: '/service-accounts',
    name: 'service-accounts',
    component: ServiceAccountsView,
    meta: { requiresAuth: true },
  },
  {
    path: '/cover-letter',
    name: 'cover-letter',
    component: SurvivorLetterView,
    meta: { requiresAuth: true },
  },
]

const router = createRouter({
  history: createWebHistory('/my/'),
  routes,
})

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    next({ name: 'login' })
  } else {
    next()
  }
})

export default router
