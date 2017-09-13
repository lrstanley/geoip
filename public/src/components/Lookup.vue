<template>
  <div style="padding-bottom: 0;">
    <div class="ui fluid icon input input-box" :class="{ loading: loading, error: error, noresults: history.length == 0 }">
      <input type="text" v-model="address" @keyup.enter="lookup" placeholder="Lookup IP address, e.g. '8.8.8.8' or host, e.g. 'google.com'" autofocus>
      <i class="circular search link icon" @click="lookup"></i>
    </div>

    <div v-if="error" class="ui negative message">{{ error }}</div>

    <div v-if="history.length != 0" class="result-box">
      <button class="ui right floated mini button" @click="clearHistory"><i class="ui remove circle icon"></i> Clear history</button>

      <transition-group name="list" appear>
        <div v-for="item in history" class="result" :key="item.ip">
          <div class="fluid ui raised card">
            <div class="content">
              <i v-if="item.country_abbr.length != 0" :class="[item.country_abbr.toLowerCase()]" class="right floated flag"></i>

              <div class="header">Q: {{ item.query }} <span v-if="item.query != item.ip">({{ item.ip }})</span></div>
              <div class="meta">
                <span class="category" v-if="item.hosts.length != 0">({{ item.hosts.join(', ') }})</span>
              </div>
              <div class="description">
                <iframe class="map" width="100%" height="110px" frameborder="0" scrolling="no" marginwidth="0" marginheight="0" :src="item.map_embed_url"></iframe>
              </div>
            </div>
            <div class="extra content">
              <div v-if="item.timezone" class="ui label"><i class="wait icon"></i> {{item.timezone}}</div>
              <!-- <button class="ui button">Join Project</button> -->

              <span class="ui right floated">
                <div v-if="item.summary" class="ui label"><i class="map icon"></i> {{item.summary}}</div>
              </span>
            </div>
          </div>
        </div>
      </transition-group>
    </div>
  </div>
</template>

<script>
export default {
  name: "lookup",
  data: function () {
    return {
      address: "",
      error: false,
      loading: false,
      history: [],
    }
  },
  methods: {
    lookup: function (lookupSelf) {
      let query = lookupSelf === true ? 'self' : this.address

      if (query.length == 0 || this.loading) { return }

      this.error = false
      this.loading = true

      this.$http.get(`/api/${query}`).then(response => {
        console.log(response)
        this.loading = false
        this.address = ""

        // Add our query into the result, so when it gets saved to history,
        // we can use it later.
        response.body.query = query;

        // Add the result to lookup history.
        this.history.unshift(response.body)

        // Make sure we're only storing the last ~10 items.
        if (this.history.length > 10) {
          this.history = this.history.split(0, 10)
        }

        // And save it to localstorage.
        this.$ls.set("history", JSON.stringify(this.history))
      }, response => {
        console.log(response)
        this.error = true
        this.loading = false
      });
    },
    clearHistory: function () {
      this.history = []
      this.$ls.set("history", JSON.stringify([]))
    }
  },
  mounted: function () {
    // On load, try looking in localstorage to see if they have any previous
    // results.
    this.history = JSON.parse(this.$ls.get("history", []))

    // If they don't have any results, lookup their own IP.
    if (this.history.length == 0) {
      this.lookup(true)
    }
  }
}
</script>

<style>
.input-box {
  transition: 0.15s padding ease-out, 0.15s margin ease-out, 0.15s border ease-out;
}

.input-box.noresults:not(.error) {
  margin-top: 100px;
  margin-bottom: 100px;
}

.input-box.noresults.error { margin-top: 100px; }
.input-box:not(.noresults):not(.error) input { margin-bottom: 20px; }

.result-box {
  border-top: 1px lightgray solid;
  padding-top: 20px;
  margin-left: -25px;
  margin-right: -25px;
  padding: 15px 15px 0 15px;
}

.result-box .result:first-child { margin-top: 40px; }
.result-box .result:not(:first-child) { margin-top: 15px; }

.list-enter, .list-leave-to { opacity: 0; }
.list-enter-active, .list-leave-active {
  animation-duration: .3s;
  animation-name: fadeInRight;
}

@keyframes fadeInRight {
   0% { opacity: 0; transform: translateX(100px); }
   100% { opacity: 1; transform: translateX(0); }
}
</style>
