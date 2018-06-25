<template>
  <v-app>
     <v-navigation-drawer v-model="sidebar" app>
      <v-list>
        <v-list-tile to="/" v-if="showHome">
          <v-list-tile-action><v-icon>home</v-icon></v-list-tile-action>
          <v-list-tile-content>Home</v-list-tile-content>
        </v-list-tile>
        <v-list-tile to="/addResVis" v-if="showAddResVis">
          <v-list-tile-action><v-icon></v-icon></v-list-tile-action>
          <v-list-tile-content>Add visit</v-list-tile-content>
        </v-list-tile>
        <v-list-tile to="/admin" v-if="showAdmin">
          <v-list-tile-action><v-icon></v-icon></v-list-tile-action>
          <v-list-tile-content>Admin</v-list-tile-content>
        </v-list-tile>
        <v-list-tile to="/login" v-if="showLogin">
          <v-list-tile-action><v-icon></v-icon></v-list-tile-action>
          <v-list-tile-content>Login</v-list-tile-content>
        </v-list-tile>
        <v-list-tile @click="logout" v-if="showLogout">
          <v-list-tile-action><v-icon></v-icon></v-list-tile-action>
          <v-list-tile-content>Logout</v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar app v-bind:color="toolbarColor">
      <span class="hidden-sm-and-up">
        <v-toolbar-side-icon @click="sidebar = !sidebar">
        </v-toolbar-side-icon>
      </span>
      <img src="/static/ResVisLogo.png" height="40px">
      <v-toolbar-title>
        Restaurant Visits
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-toolbar-items class="hidden-xs-only">
        <v-btn flat v-if="showAddResVis" to="/addResVis">Add visit</v-btn>
        <v-btn flat v-if="showHome" to="/">Home</v-btn>
        <v-btn flat v-if="showAdmin" to="/admin">Admin</v-btn>
        <v-btn flat v-if="showLogout" @click="logout">Logout</v-btn>
        <v-btn flat v-if="showLogin" to="/login">Login</v-btn>
      </v-toolbar-items>
    </v-toolbar>
    <v-content>
      <router-view/>
    </v-content>
  </v-app>
</template>

<script>
  export default {
    data() {
      return {
        hasAdmin: true,
        sidebar: false
      }
    },
    methods: {
      logout() {
        this.$store.commit("logout", {})
        this.$router.push("/login")
        console.log("logout")
      },
      admin() {
        this.$router.push("/admin")
      }
    },
    name: 'App',
    computed: {
      showAdmin() {
        var user = this.$store.state.user
        var userIsAdmin = (user != undefined && user.settings != undefined && user.settings.Admin == "true")
        console.log("Is admin:", userIsAdmin, user)
        return userIsAdmin && this.$route.path != "/login"
      },
      isLoginPage() {
        return this.$route.path == "/login"
      },
      showAddResVis() {
        return this.$store.state.isLoggedIn
      },
      showHome() {
        return this.$store.state.isLoggedIn
      },
      showLogout() {
        return this.$store.state.isLoggedIn && !this.isLoginPage
      },
      showLogin() {
        return !this.$store.state.isLoggedIn && !this.isLoginPage
      },
      toolbarColor() {
        if (this.$store.state.user == "" || this.$store.state.user.settings == undefined) {
          return "primary"
        }
        if (this.$store.state.user.settings["DB"] == "test") {
          return "green"
        }
        return "primary"
      }
    },
    created: function () {
      console.log("App created")
      //this.$store.commit("loadUser")
    }

  }

</script>

<style>
  a:link {
    color: #ff0000;
  }

  a:visited {
    color: #0000ff;
  }

  a:hover {
    color: #ffcc00;
  }

</style>
