<template>
  <v-container fluid>
    <v-layout row wrap>
      <v-flex xs0 sm2></v-flex>
      <v-flex xs12 sm8 mt-5 v-if="!visIsPublic">
        <v-layout column>
          <v-flex xs12 class="text-xs-center">
            <div class="headline">Visit is not public :-(</div>
          </v-flex>
          <v-flex class="text-xs-center" mt-3>
            <v-alert color="error" :value="error">
              {{error}}
            </v-alert>
          </v-flex>
        </v-layout>
      </v-flex>
      <v-flex xs12 sm8 mt-5 ref="mycolumn" v-resize="resizecallback" v-if="visIsPublic">
        <v-layout column>
          <v-flex xs12 class="text-xs-center">
            <v-card>
              <v-card-media contain>
                <v-flex xs12 class="text-xs-center">
                  <iframe width="600" height="300" frameborder="0" style="border:0" :src="getResAddressMap()" allowfullscreen>
                  </iframe>
                </v-flex>
              </v-card-media>
              <v-card-text>
                <v-layout column>
                  <v-flex xs12>
                    <div class="headline">{{getResName()}}</div>
                  </v-flex>
                  <v-flex xs12>
                    {{getResAddress()}}
                  </v-flex>
                  <v-flex xs12>
                    {{formatTime(visit.time)}}
                  </v-flex>
                  <v-flex xs12 v-if="visit.with">
                    {{name}} with {{visit.with}}
                  </v-flex>
                  <v-flex xs12 v-if="!visit.with">
                    {{name}}
                  </v-flex>
                  <v-flex xs12 v-if="visit.comment" mt-3>
                    <b>{{visit.comment}}</b>
                  </v-flex>
                </v-layout>
              </v-card-text>
            </v-card>
          </v-flex>
          <v-flex class="text-xs-center" mt-3>
            <v-alert color="error" :value="error">
              {{error}}
            </v-alert>
          </v-flex>
        </v-layout>
        <v-layout row wrap>
          <v-flex xs12 mt-5 v-for="p in picturesWithPath" :key="p.name">
            <v-card :key="refreshMedia">
              <v-card-media :src="p.fullPath" :height="imgHeight(p)" contain></v-card-media>
              <v-card-text>
                {{p.comment}}
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-flex>
      <v-flex xs2></v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import axios from 'axios'
  export default {
    components: {},
    data() {
      return {
        error: "",
        visit: {
          time: 0,
          with: "",
          restaurantid: 0,
          pictures2: []
        },
        restaurant: "",
        name: "",
        refreshMedia: 0,
        visIsPublic: true
      }
    },
    methods: {
      resizecallback() {
        console.log("Resize")
        this.refreshMedia = this.refreshMedia + 1
      },
      formatTime(t) {
        if (t == 0)
          return ""
        var options = {
          weekday: 'long',
          year: 'numeric',
          month: 'long',
          day: 'numeric'
        };
        var date = new Date(t);
        return date.toLocaleDateString("dk-DK", options)
      },
      getResName() {
        if (this.restaurant != "")
          return this.restaurant.name
        return ""
      },
      getResAddress() {
        if (this.restaurant != "" && this.restaurant.address != undefined)
          return this.restaurant.address.street + ", " + this.restaurant.address.city;
        return ""
      },
      getResAddressMap() {
        if (this.restaurant != "" && this.restaurant.address != undefined) {
          var address = this.restaurant.address.street + "," + this.restaurant.address.city 
          console.log("address = ", address)
          return "https://www.google.com/maps/embed/v1/place?key=AIzaSyDW3b3gV4L1lOAjoMaZA16WlAJehGXWjtU&q="+address
        }
        return ""
      },
      imgHeight(p) {
        if (this.$refs.mycolumn != undefined) {
          var cWidth = this.$refs.mycolumn.clientWidth
          return Math.min(1000, (cWidth * p.height) / p.width)
        }
        return 200
      },
      loadVisit() {
        var userid = this.$route.params.userid
        var visid = this.$route.params.visitid
        console.log("Load visit")
        axios.get('http://tocsig.ddns.net:3000/api/resvis/public/' + userid + '/' + visid)
          .then((resp) => {
            console.log("resp:")
            console.log(resp)
            if (resp.data.visit.ispublic == false) {
              this.visIsPublic = false
              console.log("Visit it not public")
              return
            }
            this.visIsPublic = true
            this.visit = resp.data.visit
            this.restaurant = resp.data.restaurant
            this.name = resp.data.name
          })
          .catch((err) => {
            this.error = "Can't get public res vis"
            console.log("Can't get public res vis")
          });
      }
    },
    created: function () {
      console.log("Created")
      this.loadVisit()
    },
    watch: {
      '$route.params.visitid': function (id) {
        console.log("Route changed")
        this.loadVisit()
      },
    },
    computed: {
      picturesWithPath() {
        //console.log("Calculating pictures with path")
        var pictures2 = this.visit.pictures2
        if (pictures2 == undefined) {
          return []
        }
        var ps = pictures2.map(function (p) {
          return {
            name: p.filename,
            width: p.width,
            height: p.height,
            fullPath: "/static/" + p.filename
          }
        })
        return ps
      },
    },
  }

</script>

<style scoped>


</style>
