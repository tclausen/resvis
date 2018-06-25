<template>
  <v-container fluid>
    <v-layout row wrap>
      <v-flex xs0 sm2></v-flex>
      <v-flex xs12 sm8 mt-5 ref="mycolumn" v-resize="resizecallback">
        <v-layout column>
          <v-flex xs12 class="text-xs-center">
            <v-card fluid>
              <v-card-media contain>
                <v-flex xs12 class="text-xs-center">
                  <iframe width="600" height="300" frameborder="0" style="border:0" :src="getResAddress(visit.restaurantid)" allowfullscreen>
                  </iframe>
                </v-flex>
              </v-card-media>
              <v-card-text>
                <v-layout column>
                  <v-flex xs12>
                    <div class="headline">{{getResName(visit.restaurantid)}}</div>
                  </v-flex>
                  <v-flex xs12>
                    {{formatTime(visit.time)}}
                  </v-flex>
                  <v-flex xs12 v-if="visit.with">
                    with {{visit.with}}
                  </v-flex>
                  <v-flex xs12 v-if="visit.comment" mt-3>
                    <b>{{visit.comment}}</b>
                  </v-flex>
                  <v-flex xs12>
                    <v-layout row wrap>
                      <v-flex xs12>
                        <v-btn flat color="orange" @click="editVisit()" :disabled="busy">Edit</v-btn>
                        <v-btn flat color="orange" @click="addImages()" :disabled="busy">Add images</v-btn>
                        <v-btn flat color="orange" @click="removeVisit()" :disabled="busy">Remove</v-btn>
                      </v-flex>
                      <v-flex xs12>
                        <v-layout row>
                          <v-flex xs6>
                            <v-switch :label="`Public`" v-model="localIsPublic" :hint="publicLink" persistent-hint></v-switch>
                            <textarea v-model="publicLink" ref="text"></textarea>
                          </v-flex>
                          <v-flex xs6>
                            <v-btn v-if="localIsPublic" flat color="orange" @click="copyToClipboard('#p1')">Copy link</v-btn>
                          </v-flex>
                        </v-layout>
                      </v-flex>
                      <label>
                        <input type="file" id="files" ref="files" multiple v-on:change="handleFileUploads()" accept="image/*" />
                      </label>
                    </v-layout>
                  </v-flex>
                </v-layout>
              </v-card-text>
            </v-card>
          </v-flex>
          <v-flex mt-5 v-if="imageData.length">
            <v-card>
              <v-card-text>
                <v-flex xs12 class="text-xs-center">
                  <div class="headline">
                    New images
                  </div>
                </v-flex>
              </v-card-text>
              <v-layout row wrap>
                <v-flex xs4 v-for="d in imageData" v-bind:key="d">
                  <img class="preview" :src="d">
                </v-flex>
              </v-layout>
              <v-layout>
                <v-flex class="text-xs-center" mt-3>
                  <v-btn color="primary" v-on:click="submit" :disabled="busy" >Save</v-btn>
                </v-flex>
              </v-layout>s
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
              <v-card-text style="position: relative">
                <v-flex xs12 class="text-xs-center">
                  {{p.comment}}
                </v-flex>
                <v-btn flat absolute dark fab style="right: 70px; bottom: -20px;"
                       color="pink"v-on:click="editComment(p)"><v-icon>edit</v-icon></v-btn>
                <v-btn flat absolute dark fab style="right: 10px; bottom: -20px;" right 
                       color="pink"v-on:click="removeImage(p.name)"><v-icon>delete</v-icon></v-btn>
              </v-card-text>
              <v-card-actions>
                <!--                <v-btn icon>
                  <v-icon>favorite</v-icon>
                </v-btn>
                <v-btn flat v-on:click="rotateImage(p.name)">
                  Rotate
                </v-btn>
-->
              </v-card-actions>
            </v-card>
          </v-flex>
        </v-layout>
      </v-flex>
      <v-flex xs2>
        <v-dialog v-model="dialog" max-width="500px">
          <v-card>
            <v-card-title>
              <span class="headline">Edit comment</span>
            </v-card-title>
            <v-card-media :src="editedPicture.fullPath" height="200"></v-card-media>
            <v-card-text>
              <v-container grid-list-md>
                <v-layout wrap>
                  <v-flex xs12 sm12 md12>
                    <v-text-field label="Comment" v-model="editedComment" textarea></v-text-field>
                  </v-flex>
                </v-layout>
              </v-container>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="blue darken-1" flat @click.native="closeCommentDialog">Cancel</v-btn>
              <v-btn color="blue darken-1" flat @click.native="saveComment">Save</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-flex>
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
        publicLink: "",
        files: [],
        imageData: [],
        filesUploaded: [],
        refreshMedia: 0,
        busy: false,
        dialog: false,
        editedComment: "unset",
        editedPicture: ""
      }
    },
    methods: {
      resizecallback() {
        console.log("Resize")
        this.refreshMedia = this.refreshMedia + 1
      },
      copyToClipboard(element) {
        console.log("Copy to clipboard")
        console.log(this.$refs)
        this.$refs.text.select();
        console.log("Copying: " + this.$refs.text)
        document.execCommand('copy');
      },
      formatTime(t) {
        var options = {
          weekday: 'long',
          year: 'numeric',
          month: 'long',
          day: 'numeric'
        };
        var date = new Date(t);
        return date.toLocaleDateString("dk-DK", options)
      },
      closeCommentDialog() {
        this.dialog = false;
        setTimeout(() => {
          this.editedComment = "";
        }, 300);
      },
      getResName(id) {
        if (id == undefined) {
          return "Undefined id"
        }
        console.log("getresname id = ", id)
        var rs = this.$store.state.user.restaurants
        console.log("ress = ", rs)
        if (rs != undefined && id in rs)
          return rs[id].name
        return "Unknown"
      },
      getResAddress(id) {
        if (id == undefined) {
          return "Undefined id"
        }
        console.log("getresaddress id = ", id)
        var rs = this.$store.state.user.restaurants
        console.log("ress = ", rs)
        if (rs != undefined && id in rs) {
          var r = rs[id]
          if(r.address != undefined) {
            var address = r.address.street + "," + r.address.city 
            console.log("address = ", address)
            return "https://www.google.com/maps/embed/v1/place?key=BIzaSyDW3b3gV4L1lOAjoMaZA16WlAJehGXWjtU&q="+address
            }
        }
        return "Unknown"
      },
      editVisit() {
        console.log("Edit visit")
        console.log(this.visit)
        this.error = "Edit visit not implemented yet"
      },
      removeVisit() {
        console.log("Remove visit")
        var visit = this.visit
        console.log(visit)
        if (confirm('Are you sure you want to delete visit to ' + this.getResName(visit.restaurantid) + '?')) {
          let config = {
            headers: {
              "Token": this.$store.state.token
            }
          }
          axios.get('http://tocsig.ddns.net:3000/api/resvis/delete/' + visit.id, config)
            .then((resp) => {
              //console.log("resp: " + resp)
              console.log("id " + visit.id)
              this.error = ""
              this.$store.commit("loadUser")
              this.$router.push("/")
            })
            .catch((err) => {
              console.log(err.response)
              this.error = "Can't delete restaurant visit: " + err.response
            });
        }
      },
      addImages() {
        this.error = ""
        this.$refs.files.click();
      },
      handleFileUploads() {
        console.log("Handle file uploads - number of files: " + this.$refs.files.files.length)
        for (var i = 0; i < this.$refs.files.files.length; i++) {
          var f = this.$refs.files.files[i]
          var found = false
          for (var ii = 0; ii < this.files.length; ii++) {
            console.log(f.name, this.files[ii].name)
            if (f.name == this.files[ii].name) {
              found = true
              break
            }
          }
          if (found) {
            continue
          }
          console.log("push file " + f.name)
          this.files.push(f)
          var reader = new FileReader();
          // Define a callback function to run, when FileReader finishes its job
          reader.onload = (e) => {
            // Note: arrow function used here, so that "this.imageData" refers to the imageData of Vue component
            // Read image as base64 and set to imageData
            this.imageData.push(e.target.result)
          }
          // Start the reader job - read file as a data url (base64 format)
          reader.readAsDataURL(f);
        }
        console.log(this.files)
      },
      submit() {
        console.log("Update resvis");
        if(this.busy) {
          this.error = "Busy uploading visit"
          return
        }
        this.busy = true
        Promise.all([this.submitFiles()])
          .then(() => {
            console.log("image upload done")
            this.submitVisit()
          })
      },
      submitFiles() {
        console.log("Submit files")
        let formData = new FormData();
        for (var i = 0; i < this.files.length; i++) {
          let file = this.files[i];
          console.log("Append file ", i)
          console.log(file)
          formData.append('files[' + i + ']', file);
        }
        let config = {
          headers: {
            Token: this.$store.state.token,
            'Content-Type': 'application/x-www-form-urlencoded',
          }
        }
        var self = this
        return axios.post("http://tocsig.ddns.net:3000/api/resvis/uploadimage", formData, config)
          .then(resp => {
            for (var i = 0; i < resp.data.length; i++) {
              self.filesUploaded.push(resp.data[i])
            }
          })
          .catch(function () {
            self.error = "Failed to upload images"
            console.log('FAILURE!!');
          });
      },
      submitVisit() {
        console.log("submit visit")
        console.log(this.visit)
        var visit = this.visit
        if (visit.pictures2 == undefined) {
          visit.pictures2 = []
        }
        for (var i = 0; i < this.filesUploaded.length; i++) {
          visit.pictures2.push({
            filename: this.filesUploaded[i]
          })
        }
        var json = JSON.stringify(visit);
        console.log("json: " + json);
        let config = {
          headers: {
            Token: this.$store.state.token
          }
        };
        // Reset state used for uploading files
        this.files = []
        this.filesUploaded = []
        this.imageData = []
        // Copy visit to force re-computation of this.visit
        this.visit = JSON.parse(json)
        visit = this.visit
        axios
          .post("http://tocsig.ddns.net:3000/api/resvis/update", json, config)
          .then(resp => {
            this.busy = false
            this.error = "";
            this.$store.commit("loadUser")
          })
          .catch(err => {
            if (!err.response) {
              this.busy = false
              this.error = "No response from server";
              console.log("No response from server");
              return;
            }
            this.error = "Unknown error";
          });
      },
      removeImage(name) {
        console.log("Remove image " + name)
        if (confirm('Are you sure you want to delete image ' + name)) {
          for (var i = 0; i < this.visit.pictures2.length; i++) {
            if (this.visit.pictures2[i].filename == name) {
              this.visit.pictures2.splice(i, 1)
              break
            }
          }
          this.submitVisit();
        }
      },
      editComment(p) {
        console.log("Edit comment: ", p)
        this.dialog = true
        this.editedPicture = p
        if(p.comment != undefined) {
          this.editedComment = p.comment
        } else {
          this.editedComment = ""
        }
      },
      saveComment() {
        console.log("Save comment: ", this.editedComment)
        var pictures2 = this.visit.pictures2
        for (var i = 0; i < pictures2.length; i++) {
          if(pictures2[i].filename == this.editedPicture.name) {
            pictures2[i].comment = this.editedComment
          }
        }
        this.submitVisit()
        this.closeCommentDialog()
      },
      imgHeight(p) {
        if (this.$refs.mycolumn != undefined) {
          var cWidth = this.$refs.mycolumn.clientWidth
          //console.log((cWidth * p.height) / p.width)
          return Math.min(1000, (cWidth * p.height) / p.width, p.height)
        } 
        //console.log("mycolumn is undefined")
        return 200
      },
      rotateImage(name) {
        console.log("Rotate image " + name)
      }
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
            comment: p.comment,
            fullPath: "/static/" + p.filename
          }
        })
        return ps
      },
      localIsPublic: {
        get() {
          if (this.visit.ispublic) {
            this.publicLink = "http://tocsig.ddns.net:8080/#/visit/" + this.$store.state.user.id + "/" + this.$route.params
              .id
          } else {
            this.publicLink = ""
          }
          return this.visit.ispublic
        },
        set(newValue) {
          this.visit.ispublic = newValue
          this.submitVisit()
        }
      },
      restaurantAddress: {
        get: function() {
          if (this.visit.restaurantid == 0) {
            return ""
          }
          var id = this.visit.restaurantid
          console.log("hasAddress id = ", id)
          var rs = this.$store.state.user.restaurants
          console.log("ress = ", rs)
          if (rs != undefined && id in rs) {
            var r = rs[id]
            if(r.address != undefined) {
              return r.address
            }
          }
          return ""
        }
      },
      visit: {
        get: function () {
          console.log("Compute visit")
          if (this.$route.params.id != undefined) {
            if (this.$store.state.user.resvisits != undefined) {
              return this.$store.state.user.resvisits[this.$route.params.id]
            }
          }
          return {
            time: 0,
            with: "",
            restaurantid: 0,
            pictures2: []
          }
        },
        set: function (newVisit) {
          console.log("Set visit")
          this.$store.state.user.resvisits[this.$route.params.id] = newVisit
        }
      }
    },
  }

</script>

<style scoped>
  input[type="file"] {
    position: absolute;
    top: -500px;
  }

  textarea {
    position: absolute;
    top: -500px;
  }

  img.preview {
    background-color: white;
    border: 1px solid #DDD;
    padding: 5px;
    max-width: 100%;
  }

  .rotate90 {
    -webkit-transform: rotate(90deg);
    -moz-transform: rotate(90deg);
    -o-transform: rotate(90deg);
    -ms-transform: rotate(90deg);
    transform: rotate(90deg);
  }

</style>
