<template>
  <div>
    <div class="ui fluid icon input" :class="{ loading: loading, noresults: expanded }">
      <input type="text" v-model="address" @keyup.enter="lookup" placeholder="Lookup IP address, e.g. '8.8.8.8' or host, e.g. 'google.com'" :class="{ error: error }" autofocus>
      <i class="circular search link icon" @click="lookup"></i>
    </div>
  </div>
</template>

<script>
export default {
  name: "lookup",
  data: function () {
    return {
      address: "",
      expanded: true,
      error: false,
      loading: false,
    }
  },
  methods: {
    lookup: function () {
      if (this.address.length == 0 || this.loading) { return }

      this.error = false
      this.loading = true

      this.$http.get(`/api/${this.address}`).then(response => {
        console.log(response)
        this.loading = false
      }, response => {
        console.log(response)
        this.error = true
        this.loading = false
      });
    },
  },
  mounted: function () {
    // TODO: check localstorage here and see if there are any entries.
    // If not, we should fetch the users IP itself.
  }
}
</script>

<style>
.noresults {
  margin-top: 100px;
  margin-bottom: 100px;
}
</style>
