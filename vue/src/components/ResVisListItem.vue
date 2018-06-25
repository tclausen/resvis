<template>
<v-card>
  <v-card-media v-if="hasPictures()"
     class="white--text" height="200px"
     :src="getPicture()"
     >
    <v-container fill-height fluid>
      <v-layout fill-height>
        <v-flex xs12 align-end flexbox>
          <span class="headline">{{getResName()}}</span>
        </v-flex>
      </v-layout>
    </v-container>
  </v-card-media>
  <v-card-title>
    <div>
      <span class="grey--text">{{getResName()}}</span>
      <span class="grey--text" v-if="hasPictures()">({{numberOfPictures()}} images)</span>
      <br>
      <span>{{formatTime(visit.time)}} with {{visit.with}}</span>
    </div>
  </v-card-title>
  <v-card-actions>
    <v-btn flat color="orange" v-on:click="gotoDetail(visit.id)">Details</v-btn>
  </v-card-actions>
</v-card>
</template>

<script>
export default {
  props: {
    visit: {
      type: Object,
      required: true
		}
	},
	methods: {
		formatTime(t) {
			var options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
			var date = new Date(t);
			return date.toLocaleDateString("dk-DK",options)
		},
    getResName() {
      var rs = this.$store.state.user.restaurants
      if(rs == undefined) {
        return "Unknown"
      }
      if(this.visit.restaurantid in rs)
        return rs[this.visit.restaurantid].name
      return "Unknown"
    },
    getPicture() {
      if(this.hasPictures()) {
        return "/static/"+this.visit.pictures2[0].filename
      }
      return "no pictures"
    },
    hasPictures() {
      if(this.visit.pictures2 !== undefined && this.visit.pictures2.length > 0) {
        return true
      }
      return false
    },
    numberOfPictures() {
      if(this.visit.pictures2 !== undefined) {
        return this.visit.pictures2.length
      }
      return 0
    },
    gotoDetail(id) {
      this.$router.push("/detail/"+id)
    }

	}
}
</script>
