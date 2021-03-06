<template>
<div>
  <v-container fluid>
    <v-layout row wrap>
      <v-flex xs12 mt-0>
        <v-btn color="primary" v-on:click="refresh">
          Refresh
        </v-btn>
        <v-dialog v-model="dialog" max-width="500px">
          <v-btn color="primary" dark slot="activator" class="mb-2">Add restaurant</v-btn>
          <v-card>
            <v-card-title>
              <span class="headline">{{ formTitle }}</span>
            </v-card-title>
            <v-card-text>
              <v-container grid-list-md>
                <v-layout wrap>
                  <v-flex xs12 sm12 md12>
                    <v-text-field label="Name" v-model="editedItem.name"></v-text-field>
                    <v-text-field label="Street" v-model="editedItem.address.street"></v-text-field>
                    <v-text-field label="City" v-model="editedItem.address.city"></v-text-field>
                    <v-text-field label="Country" v-model="editedItem.address.country"></v-text-field>
                  </v-flex>
                </v-layout>
              </v-container>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="blue darken-1" flat @click.native="close">Cancel</v-btn>
              <v-btn color="blue darken-1" flat @click.native="save">Save</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <!-- <AddRestaurantDialog v-bind:isVisible="dialog" v-bind:editedItem="editedItem"></AddRestaurantDialog>-->
      </v-flex>
    </v-layout>
  </v-container>
  <v-flex class="text-xs-center" mt-3>
    <v-alert color="error" :value="error" dismissible v-model="hasError" transition="scale-transition">
      {{error}}
    </v-alert>
  </v-flex>
  <v-data-table
     :headers="headers"
     :items="items"
     hide-actions
     class="elevation-1"
     >
    <template slot="items" slot-scope="props">
      <td>{{ props.item.id }}</td>
      <td>{{ props.item.name }}</td>
      <td>{{ props.item.address.street }}</td>
      <td>{{ props.item.address.city }}</td>
      <td>{{ props.item.address.country }}</td>
      <td class="justify-center layout px-0">
        <v-btn icon class="mx-0" @click="editItem(props.item)">
          <v-icon color="teal">edit</v-icon>
        </v-btn>
        <v-btn icon class="mx-0" @click="deleteItem(props.item)">
          <v-icon color="pink">delete</v-icon>
        </v-btn>
      </td>
    </template>
  </v-data-table>
</div>
</template>

<script>
import axios from 'axios'
import AddRestaurantDialog from './AddRestaurantDialog.vue'
export default {
  components: {
    AddRestaurantDialog
  },
  name: 'AdminRestaurants',
  data () {
    return {
      error: "",
      headers: [
        { text: 'Id', align: 'left', sortable: true, value: 'id' },
        { text: 'Name', align: 'left', sortable: true, value: 'name' },
        { text: 'Street', align: 'left', sortable: true, value: 'Street' },
        { text: 'City', align: 'left', sortable: true, value: 'City' },
        { text: 'Country', align: 'left', sortable: true, value: 'country' },
        { text: 'Actions', align: 'left', sortable: false },
      ],
      items: [],
      dialog: false,
      editedIndex: -1,
      editedItem: {
        name: '',
        address: {},
      },
      defaultItem: {
        name: '',
        address: {},
      }
    }
  },
  created: function() {
    this.getRestaurants()
  },
  computed: {
    formTitle () {
      return this.editedIndex === -1 ? 'Add restaurant' : 'Edit restaurant'
    },
    hasError: {
      // getter
      get: function () {
        return this.error != ""
      },
      // setter
      set: function (newValue) {
        this.error = ""
      }
    }
  },
  methods: {
    getRestaurants() {
      let config = { headers: { "Token": this.$store.state.token } }
      axios.get('http://tocsig.ddns.net:3000/api/restaurant', config)
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
          this.error = "Can't load list of restaurants: " + err.response.data
        });
    },
    refresh() {
      this.getRestaurants()
    },
    close () {
      this.dialog = false
      setTimeout(() => {
        this.editedItem = Object.assign({}, this.defaultItem)
        this.editedIndex = -1
      }, 300)
    },
    save () {
      if (this.editedIndex > -1) {
        this.updateRestaurant(this.editedItem)
      } else {
        this.addRestaurant(this.editedItem)
      }
      this.close()
    },
    deleteItem (item) {
      const index = this.items.indexOf(item)
      if(confirm('Are you sure you want to delete '+item.name+'?')) {
        this.deleteRestaurant(item)
      }
    },
    editItem (item) {
      this.editedIndex = this.items.indexOf(item)
      this.editedItem = Object.assign({}, item)
      this.dialog = true
    },
    deleteRestaurant(item) {
      console.log("Deleting: " + item.name)
      let config = { headers: { "Token": this.$store.state.token } }
      axios.get('http://tocsig.ddns.net:3000/api/restaurant/delete/'+item.id, config)
        .then((resp) => {
          this.error = ""
          const index = this.items.indexOf(item)
          this.items.splice(index, 1)
        })
        .catch((err) => {
          console.log(err.response)
          this.error = "Can't delete restaurant: " + err.response.data
        });
    },
    addRestaurant(item) {
      console.log("Adding: " + item.name)
      var json = JSON.stringify(item)
      console.log("json:" + json)
      let config = { headers: { "Token": this.$store.state.token } }
      axios.post('http://tocsig.ddns.net:3000/api/restaurant', json, config)
        .then((resp) => {
          console.log("received:", resp.data)
          this.editedItem.id = resp.data.id
          this.error = ""
          this.items.push(this.editedItem)
        })
        .catch((err) => {
          console.log(err.response)
          this.error = "Can't add restaurant: " + err.response.data
        });
    },
    updateRestaurant(item) {
      console.log("Adding: " + item.name)
      var json = JSON.stringify(item)
      console.log("json:" + json)
      let config = { headers: { "Token": this.$store.state.token } }
      axios.post('http://tocsig.ddns.net:3000/api/restaurant/update', json, config)
        .then((resp) => {
          console.log("received:", resp.data)
          Object.assign(this.items[this.editedIndex], this.editedItem)
          this.error = ""
        })
        .catch((err) => {
          console.log(err.response)
          this.error = "Can't add restaurant: " + err.response.data
        });
    },
  }
}
</script>
  
<style>
</style>
  
