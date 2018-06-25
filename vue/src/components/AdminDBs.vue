<template>
<div>
  <v-container fluid>
    <v-layout row wrap>
      <v-flex xs12 mt-0>
        <v-btn color="primary" v-on:click="refresh">
          Refresh
        </v-btn>
        <v-dialog v-model="dialog" max-width="500px">
          <v-btn color="primary" dark slot="activator" class="mb-2">Add database</v-btn>
          <v-card>
            <v-card-title>
              <span class="headline">Add database</span>
            </v-card-title>
            <v-card-text>
              <v-container grid-list-md>
                <v-layout wrap>
                  <v-flex xs12 sm6 md4>
                    <v-text-field label="Name" v-model="editedItem.name"></v-text-field>
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
      <td>{{ props.item.name }}</td>
      <td>{{ props.item.active }}</td>
      <td class="justify-left layout px-0">
        <v-btn v-if="!isActive(props.item)" icon class="mx-0" @click="deleteItem(props.item)">
          <v-icon color="pink">delete</v-icon>
        </v-btn>
        <v-btn v-if="!isActive(props.item)" class="mx-0" @click="activateDB(props.item.name)">
          Activate
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
      error: "",
      headers: [
        { text: 'Name', align: 'left', sortable: true, value: 'name' },
        { text: 'Active', value: 'active' },
        { text: 'Actions', align: 'left', sortable: false },
      ],
      items: [],
      dialog: false,
      editedIndex: -1,
      editedItem: {
        name: '',
        address: '',
      },
      defaultItem: {
        name: '',
        address: '',
      },
      dbActive: ''
    }
  },
  created: function() {
    this.getDBs()
  },
  computed: {
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
    getDBs() {
      console.log("Get DBs")
      let config = { headers: { "Token": this.$store.state.token } }
      axios.get('http://tocsig.ddns.net:3000/api/db', config)
        .then((resp) => {
          this.error = ""
          this.dbActive = resp.data[0]
          var dbList = resp.data.slice(1)
          var items = []
          for(var i = 0; i < dbList.length; i++) {
            if(dbList[i] == undefined) {
              continue
            }
            var item = {name: dbList[i], active: ""}
            if(item.name == this.dbActive) {
              item.active = "X"
            }
            items.push(item)
          }
          this.items = items
        })
        .catch((err) => {
          console.log(err.response)
          this.error = "Can't load database list: " + err.response.data
        });
    },
    isActive(item) {
      return item.name == this.dbActive
    },
    refresh() {
      this.getDBs()
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
        Object.assign(this.items[this.editedIndex], this.editedItem)
      } else {
        this.addDB(this.editedItem)
      }
      this.close()
    },
    deleteItem (item) {
      if(confirm('Are you sure you want to delete '+item.name+'?')) {
        this.deleteDB(item)
      }
    },
    addDB(db) {
      console.log("Adding: " + db.name)
      let config = { headers: { "Token": this.$store.state.token } }
      axios.get('http://tocsig.ddns.net:3000/api/db/create/'+db.name+'/DBMem', config)
        .then((resp) => {
          this.error = ""
          this.items.push(db)
        })
        .catch((err) => {
          console.log(err.response)
          this.error = "Can't create database: " + err.response.data
        });
    },
    activateDB(db) {
      console.log("Activating: " + db)
      let config = { headers: { "Token": this.$store.state.token } }
      axios.get('http://tocsig.ddns.net:3000/api/db/'+db, config)
        .then((resp) => {
          this.error = ""
          this.$store.commit("logout", {})
          this.$router.push("/login")
        })
        .catch((err) => {
          console.log(err.response)
          this.error = "Can't activate database: " + err.response.data
        });
      
    },
    deleteDB(item) {
      console.log("Deleting: " + item.name)
      let config = { headers: { "Token": this.$store.state.token } }
      axios.get('http://tocsig.ddns.net:3000/api/db/delete/'+item.name, config)
        .then((resp) => {
          this.error = ""
          const index = this.items.indexOf(item)
          this.items.splice(index, 1)
        })
        .catch((err) => {
          console.log(err.response)
          this.error = "Can't delete database: " + err.response.data
        });
    },
  }
}
</script>
  
<style>
</style>
  
