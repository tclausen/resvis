<template>
  <v-container fluid>
    <v-layout row wrap>
      <v-flex xs12 sm10 offset-sm1 mt-2>
        <v-tabs centered v-model="activeTab" grow>
          <v-tab>Databases</v-tab>
          <v-tab>Users</v-tab>
          <v-tab>Restaurants</v-tab>
          <v-tab>Event log</v-tab>
          <v-tab-item>
            <AdminDBs></AdminDBs>
          </v-tab-item>
          <v-tab-item>
            <AdminUsers></AdminUsers>
          </v-tab-item>
          <v-tab-item>
            <AdminRestaurants></AdminRestaurants>
          </v-tab-item>
          <v-tab-item>
            <AdminEvents></AdminEvents>
          </v-tab-item>
        </v-tabs>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import axios from 'axios'
import AdminRestaurants from './AdminRestaurants.vue'
import AdminUsers from './AdminUsers.vue'
import AdminDBs from './AdminDBs.vue'
import AdminEvents from './AdminEvents.vue'
export default {
  components: {
    AdminRestaurants,
    AdminUsers,
    AdminDBs,
    AdminEvents
  },
  data () {
    return {
      activeTab: null
    }
  },
  created: function() {
  },
  methods: {
  }
}
</script>
  
<style>
</style>
  
