<template>
<v-dialog v-model="isVisible" max-width="500px">
  <v-btn color="primary" dark slot="activator" class="mb-2">Add restaurant new</v-btn>
  <v-card>
    <v-card-title>
      <span class="headline">{{ formTitle }} new</span>
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
</template>

<script>
export default {
  props: {
    isVisible: {
      type: Boolean,
      required: true
		},
    editedItem: {
      type: Object,
      required: false,
      default: function() {
        return {
          name: '',
          address: {}
		    }
      }
    }
	},
  data () {
    return {
      editedIndex: -1,
      defaultItem: {
        name: '',
        address: {},
      }
    }
  },
  computed: {
    formTitle () {
      return this.editedIndex === -1 ? 'Add restaurant' : 'Edit restaurant'
    },
  },
  methods: {
    close () {
      this.isVisible = false
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
  
