<template>
<v-card>
  <v-card-media v-if="user.profilepicture">
     class="white--text"
     height="100px"
     src="/static/toc.jpg"
     >
  </v-card-media>
  <v-card-title>
    <div>
      <span class="grey--text">Logged in as</span><br>
      <span v-if="user.firstname">{{user.firstname}} {{user.lastname}}</span>
      <span v-else>{{user.username]}</span><br>
      <span v-if="user.email">email: {{user.email}}</span>
    </div>
    </v-card-title>
  <v-card-actions>
    <v-btn flat to="/profile" color="orange">Settings</v-btn>
  </v-card-actions>
</v-card>
</template>

<script>
export default {
  props: {
	},
	methods: {
		formatTime(t) {
			var options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
			var date = new Date(t);
			return date.toLocaleDateString("dk-DK",options)
		}
	},
  computed: {
    user() {
      return this.$store.state.user
    }
  }
}
</script>
