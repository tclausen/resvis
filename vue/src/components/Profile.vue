<template>
  <v-container fluid>
    <v-layout row wrap>
      <v-flex xs2></v-flex>
      <v-flex xs8 mt-5>
        Profile
      </v-flex>
      <v-flex xs2></v-flex>
    </v-layout>
  </v-container>
</template>

<script>
export default {
  components: {
  },
  data () {
    return {
    }
  },
  methods: {
  },
  created: function() {
    console.log("Profile: Created called")
  },
}
</script>

<style scoped>
</style>
