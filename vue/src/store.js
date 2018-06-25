import Vue from 'vue'
import Vuex from 'vuex'
import VuexPersist from 'vuex-persist'
import axios from 'axios' 

Vue.use(Vuex)

const vuexLocalStorage = new VuexPersist({
  key: 'vuexResVis', // The key to store the state on in the storage provider.
  storage: window.localStorage, // or window.sessionStorage or localForage
  // Function that passes the state and returns the state with only the objects you want to store.
  // reducer: state => state,
  reducer: state => ({
    token: state.token,
  })
  // Function that passes a mutation and lets you decide if it should update the state in localStorage.
  // filter: mutation => (true)
})


export default new Vuex.Store({
	state: {
		token: "",
		user: "",
		username: "",
    isLoggedIn: false
	},
	mutations: {
		setToken: (state, payload) => {
			console.log("setToken", payload.token)
			state.token = payload.token
			state.username = payload.username
		},
		logout: (state, payload) => {
			console.log("store logout")
      state.isLoggedIn = false
			state.token = undefined
			state.username = ""
			state.user = ""
		},
		setUser: (state, payload) => {
			console.log("setUser")
			console.log(payload.user)
			state.user = payload.user
      state.isLoggedIn = true
		},
		loadUser: (state) => {
      console.log("Re-loading user")
      console.log("state: " + state)
      let config = { headers: { "Token": state.token } }
      axios.get('http://tocsig.ddns.net:3000/api/user/fromToken', config)
        .then((resp) => {
          console.log("Got user")
          console.log("resp:")
          console.log(resp)
			    state.user = resp.data
          state.isLoggedIn = true
        })
        .catch((err) => {
          console.log("Can't get user from token")
        });
    }
	},
  plugins: [vuexLocalStorage.plugin]
})
