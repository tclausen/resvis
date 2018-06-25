import Vue from 'vue'
import Router from 'vue-router'
import axios from 'axios'
import HelloWorld from '@/components/HelloWorld'
import Login from '@/components/Login'
import Landing from '@/components/Landing'
import ResVisList from '@/components/ResVisList'
import ResVisDetail from '@/components/ResVisDetail'
import ResVisPublic from '@/components/ResVisPublic'
import AddResVis from '@/components/AddResVis'
import Admin from '@/components/Admin'
import Profile from '@/components/Profile'
import Signup from '@/components/Signup'
import store from '@/store'

Vue.use(Router)

var router = new Router({
  routes: [
    {
      path: '/',
      name: 'Landing',
      component: Landing
    },
    {
      path: '/login',
      name: 'Login',
      component: Login
    },
    {
      path: '/detail/:id',
      name: 'ResVisDetail',
      component: ResVisDetail
    },
    {
      path: '/addresvis',
      name: 'AddResVis',
      component: AddResVis
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin
    },
    {
      path: '/profile',
      name: 'Profile',
      component: Profile
    },
    {
      path: '/visit/:userid/:visitid',
      name: 'ResVisPublic',
      component: ResVisPublic
    },
    {
      path: '/signup',
      name: 'Signup',
      component: Signup
    },
  ],
})

router.beforeEach((to, from, next) => {
  if(to.path == "/login") {
    store.state.isLoggedIn = false
    next()
    return
  }
  if(to.path == "/signup") {
    next()
    return
  }
  console.log("to path", to.path)
  // Do not force login screen for public visit pages
  if(to.path.startsWith("/visit")) {
    store.state.isLoggedIn = false
    next()
    return
  }
  var needsAuth = store.state.token == "" || store.state.token == undefined
  console.log("to", store.state.token, needsAuth)
  if (needsAuth) {
    store.state.isLoggedIn = false
    next('/login')
    return
  } 
  // Do we need to read user info?
  if (store.state.user == "" || store.state.user == undefined) {
    console.log("Router load user")
    let config = { headers: { "Token": store.state.token } }
    axios.get('http://tocsig.ddns.net:3000/api/user/fromToken', config)
      .then((resp) => {
        console.log("Router load user, setUser")
        store.commit("setUser", {"user": resp.data})
        console.log(resp)
      })
      .catch((err) => {
        console.log("Can't get user from token")
      });
  }
  next()
})

export default router
