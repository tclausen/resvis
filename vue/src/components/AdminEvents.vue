<template>
<div>
  <v-container fluid>
    <v-layout row wrap>
      <v-flex xs12 mt-0>
        <v-btn color="primary" v-on:click="refresh">
          Refresh
        </v-btn>
      </v-flex>
    </v-layout>
  </v-container>
  <v-flex class="text-xs-center" mt-3>
    <v-alert color="error" :value="error">
      {{error}}
    </v-alert>
  </v-flex>
  <v-data-table
     :headers="headers"
     :items="items"
     class="elevation-1"
     v-bind:pagination.sync="pagination"
     >
    <template slot="items" slot-scope="props">
      <td>{{ formatTime(props.item.time) }}</td>
      <td>{{ props.item.category }}</td>
      <td>{{ props.item.text }}</td>
    </template>
  </v-data-table>
</div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'AdminRestaurants',
  data () {
    return {
      error: "",
      headers: [
        { text: 'Time', align: 'left', sortable: true, value: 'time' },
        { text: 'Category', align: 'left', sortable: true, value: 'category' },
        { text: 'Text', align: 'left', sortable: true, value: 'text' },
      ],
      items: [],
      pagination: {'sortBy': 'time', 'descending': true, 'rowsPerPage': 25}
    }
  },
  created: function() {
    this.getEvents()
  },
  methods: {
    getEvents() {
      let config = { headers: { "Token": this.$store.state.token } }
      axios.get('http://tocsig.ddns.net:3000/api/events', config)
        .then((resp) => {
          this.error = ""
          console.log("data:", resp.data)
          var items = []
          for(var key in resp.data) {
            items.push(resp.data[key])
          }
          this.items = items
        })
        .catch((err) => {
          console.log(err.response)
          this.error = "Can't load list of events: " + err.response.data
        });
    },
		formatTime(t) {
			var d = new Date(t);
      var datestring = d.getFullYear() + "-" + pad((d.getMonth()+1)) + "-" + pad(d.getDate()) + ", " + pad(d.getHours()) + ":" + pad(d.getMinutes()) + ":" + pad(d.getSeconds());
      return datestring
		},
    refresh() {
      this.getEvents()
    },
  }
}
function pad(num){ return ('00' + num).substr(-2); }
</script>
  
<style>
</style>
  
