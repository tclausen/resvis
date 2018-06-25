userUtils_GetUserToStore() {
  let config = { headers: { "Token": this.$store.state.token } }
  console.log("userUtils set user")
  axios.get('http://tocsig.ddns.net:3000/api/user/fromToken', config)
    .then((resp) => {
      this.error = ""
      this.$store.commit("setUser", {"user": resp.data})
      console.log(resp)
      console.log("name:", resp.data.firstname)
    })
    .catch((err) => {
      console.log("Can't get user from token")
    });
}
