<template>
  <v-container fluid style="min-height: 0;" grid-list-lg>
    <v-layout row wrap>
      <v-flex xs12 v-for="v in visitsSorted()" :key="v.id">
        <ResVisListItem v-bind:visit="v">
        </ResVisListItem>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import ResVisListItem from './ResVisListItem.vue'
  import UserInfo from './UserInfo.vue'
  export default {
    components: {
      ResVisListItem,
      UserInfo
    },
    data() {
      return {}
    },
    methods: {
      visits: function () {
        var visits = this.$store.state.user.resvisits
        return visits
      },
      visitsSorted: function () {
        var visits = this.$store.state.user.resvisits
        if (visits == undefined) {
          return []
        }
        var items = Object.keys(visits).map(function (key) {
          return [key, visits[key].time];
        });
        // Sort the array based on the second element
        items.sort(function (first, second) {
          return second[1] - first[1];
        });
        var sortedItems = items.map(function (item) {
          return visits[item[0]];
        });
        return sortedItems
      }
    }
  }

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h1,
  h2 {
    font-weight: normal;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    display: inline-block;
    margin: 0 10px;
  }

  a {
    color: #42b983;
  }

</style>
