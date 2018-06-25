<template>
<div>
  <v-container fluid>
    <v-layout row wrap>
      <v-flex xs12 mt-0>
        <v-btn color="primary" v-on:click="refresh">
          Refresh
        </v-btn>
        <v-dialog v-model="dialog" max-width="500px">
          <v-btn color="primary" dark slot="activator" class="mb-2">Add user</v-btn>
          <v-card>
            <v-card-title>
              <span class="headline">{{ formTitle }}</span>
            </v-card-title>
            <v-card-text>
              <v-container grid-list-md>
                <v-layout wrap>
                  <v-flex xs12 sm6 md4>
                    <v-text-field label="Username" v-model="editedItem.username"></v-text-field>
                    <v-text-field label="Password" v-model="editedItem.password"></v-text-field>
                    <v-text-field label="Email" v-model="editedItem.email"></v-text-field>
                    <v-text-field label="First name" v-model="editedItem.firstname"></v-text-field>
                    <v-text-field label="Last name" v-model="editedItem.lastname"></v-text-field>
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
      <td>{{ props.item.firstname }}</td>
      <td>{{ props.item.lastname }}</td>
      <td>{{ props.item.username }}</td>
      <td>{{ props.item.email }}</td>
      <td>{{ formatResvisits(props.item.resvisits) }}</td>
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
export default {
  data () {
    return {
      alert: true,
      error: "",
      headers: [
        { text: 'Id', align: 'left', sortable: true, value: 'id' },
        { text: 'First name', align: 'left', sortable: true, value: 'firstname' },
        { text: 'Last name', align: 'left', sortable: true, value: 'lastname' },
        { text: 'User name', align: 'left', sortable: true, value: 'username' },
        { text: 'Email', align: 'left', sortable: true, value: 'email' },
        { text: 'Visits', align: 'left', sortable: true, value: 'resvisits' },
        { text: 'Actions', align: 'left', sortable: false },
      ],
      items: [],
      dialog: false,
      editedIndex: -1,
      editedItem: {
        username: '',
        id: '',
        password: '',
        email: '',
        firstname: '',
        lastname: '',
      },
      defaultItem: {
        username: '',
        id: '',
        password: '',
        email: '',
        firstname: '',
        lastname: '',
      }
    }
  },
  created: function() {
    this.getUsers()
  },
  computed: {
    formTitle () {
      return this.editedIndex === -1 ? 'Add user' : 'Edit user'
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
    getUsers() {
      console.log("Get users")
      let config = { headers: { "Token": this.$store.state.token } }
      axios.get('http://tocsig.ddns.net:3000/api/user', config)
        .then((resp) => {
          this.error = ""
          //console.log("resp", resp)
          var items = []
          for(var key in resp.data) {
            items.push(resp.data[key])
          }
          this.items = items
        })
        .catch((err) => {
          console.log("Error")
          console.log(err.response)
          this.error = "Can't load list of users: " + err.response.data
        });
    },
    refresh() {
      this.getUsers()
    },
    formatResvisits(visits) {
      if(visits == undefined) {
        return "0"
      }
      return Object.keys(visits).length
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
        this.updateUser(this.editedItem)
      } else {
        this.addUser(this.editedItem)
      }
      this.close()
    },
    deleteItem (item) {
      if(confirm('Are you sure you want to delete '+item.username+'?')) {
        this.deleteUser(item)
      }
    },
    editItem (item) {
      this.editedIndex = this.items.indexOf(item)
      this.editedItem = Object.assign({}, item)
      this.dialog = true
    },
    deleteUser(item) {
      console.log("Deleting: " + item.name)
      let config = { headers: { "Token": this.$store.state.token } }
      axios.get('http://tocsig.ddns.net:3000/api/user/delete/'+item.id, config)
        .then((resp) => {
          this.error = ""
          const index = this.items.indexOf(item)
          this.items.splice(index, 1)
        })
        .catch((err) => {
          console.log(err.response)
          this.error = "Can't delete user: " + err.response.data
        });
    },
    addUser(item) {
      console.log("Adding: " + item.name)
      var json = "{\"firstname\": \""+item.firstname.trim()+"\", \"lastname\": \""+item.lastname.trim()+"\", \"username\": \""+item.username.trim()+"\", \"password\": \""+item.password+"\", \"email\": \""+item.email.trim()+"\"}"
      let config = { headers: { "Token": this.$store.state.token } }
      axios.post('http://tocsig.ddns.net:3000/api/user', json, config)
        .then((resp) => {
          console.log("received:", resp.data)
          this.editedItem.id = resp.data.id
          this.error = ""
          this.items.push(this.editedItem)
        })
        .catch((err) => {
          console.log(err.response)
          this.error = "Can't add user: " + err.response.data
        });
    },
    updateUser(item) {
      // Copy since we want to _not_ send list of visits
      var myItem = JSON.parse(JSON.stringify(item))
      console.log("Updating: " + myItem.name)
      myItem.resvisits = {}
      var json = JSON.stringify(myItem)
      console.log("json: " + json)
      let config = { headers: { "Token": this.$store.state.token } }
      axios.post('http://tocsig.ddns.net:3000/api/user/update', json, config)
        .then((resp) => {
          console.log("received:", resp.data)
          this.error = ""
          Object.assign(this.items[this.editedIndex], this.editedItem)
        })
        .catch((err) => {
          console.log(err.response)
          this.error = "Can't add user: " + err.response.data
        });
    },
  }
}
</script>
  
<style>
</style>
  
