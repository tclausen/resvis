<template>
  <v-container fluid>
    <v-layout row wrap>
      <v-flex xs0 sm2></v-flex>
      <v-flex xs12 sm5 mt-5>
        <ResVisList></ResVisList>
      </v-flex>
      <v-flex xs12 sm3 mt-5 class="hidden-xs-only">
        <v-container fluid style="min-height: 0;" grid-list-lg>
          <v-layout row wrap>
            <v-flex xs12>
              <UserInfo></UserInfo>
            </v-flex>
          </v-layout>
        </v-container>
      </v-flex>
      <v-flex xs2></v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import ResVisList from './ResVisList.vue'
import UserInfo from './UserInfo.vue'
export default {
  components: {
    ResVisList,
    UserInfo,
  }
}
</script>
