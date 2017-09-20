<template>
  <div style="padding-bottom: 0;">
    <div class="ui fluid icon input input-box" :class="{ loading: loading, error: error, noresults: history.length == 0 }">
      <input id="addr_box" type="text" v-model="address" @keyup.enter="lookup" placeholder="IP address (e.g 1.2.3.4) or host (e.g google.com)" autofocus>
      <i class="circular search link icon" @click="lookup"></i>
    </div>

    <div v-if="error" class="ui negative message"><i class="icon warning sign"></i> {{ error }}</div>

    <div v-if="history.length != 0" class="result-box">
      <button class="ui right floated mini button" @click="clearHistory"><i class="ui remove circle icon"></i> Clear history</button>

      <transition-group name="list" appear>
        <div v-for="(item, index) in history" class="result" :key="item.ip">
          <div class="fluid ui raised card">
            <div class="content">
              <i v-if="item.country_abbr.length != 0" :class="[item.country_abbr.toLowerCase()]" class="right floated flag"></i>

              <div class="header">Q: {{ item.query }} <span v-if="item.query != item.ip">(<a data-tooltip="Click to copy to your clipboard" data-inverted="" @click="copyClipboard" :data-clipboard-text="item.ip">{{ item.ip }}</a>)</span></div>
              <div class="meta">
                <span class="category" v-if="item.hosts.length != 0">(<a data-tooltip="Click to copy to your clipboard" data-inverted="" @click="copyClipboard" :data-clipboard-text="item.hosts.join(', ')">{{ item.hosts.join(', ') }}</a>)</span>
              </div>
              <div class="description" v-if="item.longitude != 0 && item.latitude != 0">
                <v-map style="height: 150px" :zoom=3 v-bind:center="[item.latitude, item.longitude]">
                  <v-tilelayer attribution="&copy; <a href='http://osm.org/copyright'>OpenStreetMap</a> contributors" url="http://{s}.tile.osm.org/{z}/{x}/{y}.png"></v-tilelayer>
                  <v-marker v-bind:lat-lng="[item.latitude, item.longitude]"></v-marker>
                </v-map>
              </div>
            </div>
            <div class="extra content">
              <div v-if="item.timezone" class="ui label"><i class="wait icon"></i> {{item.timezone}}</div>
              <span class="ui right floated">
                <div v-if="item.summary" class="ui label"><i class="map icon"></i> <a data-tooltip="Click to copy to your clipboard" data-inverted="" data-position="bottom center" @click="copyClipboard" :data-clipboard-text="item.summary">{{item.summary}}</a></div>
              </span>
            </div>
          </div>
        </div>
      </transition-group>
    </div>
  </div>
</template>

<script>
import Vue2Leaflet from 'vue2-leaflet'

export default {
  name: "lookup",
  components: {
    'v-map': Vue2Leaflet.Map,
    'v-tilelayer' :Vue2Leaflet.TileLayer,
    'v-marker': Vue2Leaflet.Marker
  },
  data: function () {
    return {
      address: "",
      error: false,
      loading: false,
      history: [],
    }
  },
  methods: {
    copyClipboard: function (event) {
      var clipboard = new Clipboard('.null');
      clipboard.onClick(event)
      toastr.success('Copied to clipboard', '', {timeOut: 1000, preventDuplicates: true})
      clipboard.destroy()
    },
    selectInput: function () {
      // Select the address input box if it's not already selected.
      setTimeout(function() { $("#addr_box").focus(); }, 500)
    },
    lookup: function (lookupSelf) {
      let query = lookupSelf === true ? 'self' : this.address

      if (query.length == 0 || this.loading) { return }

      this.error = false
      this.loading = true
      this.$Progress.start()

      this.$http.get(`/api/${query}`).then(response => {
        this.loading = false

        if (response.body.error != undefined) {
          this.error = "Error: " + response.body.error.charAt(0).toUpperCase() + response.body.error.slice(1);
          this.$Progress.fail()
          return
        }

        // Add our query into the result, so when it gets saved to history,
        // we can use it later.
        response.body.query = query;

        this.$Progress.finish()
        this.address = ""
        this.addHistory(response.body)
        this.selectInput()
      }, response => {
        this.$Progress.fail()
        this.error = "An unknown exception occurred or service unavailable"
        this.loading = false
        this.selectInput()
      });
    },
    addHistory: function (result) {
        for (var i = 0; i < this.history.length; i++) {
          if (this.history[i].query == result.query) {
            this.history.splice(i, 1)
            break
          }
        }
        // Add the result to lookup history.
        this.history.unshift(result)

        // Make sure we're only storing the last ~10 items.
        if (this.history.length > 10) {
          // this.history = this.history.split(0, 10)
          this.history.length = 10
        }

        // And save it to localstorage.
        this.$ls.set("history", JSON.stringify(this.history))
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

    this.selectInput()
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

.ui.label a, .ui.label a:link, .ui.label a:hover, .ui.label a:active, .ui.label a:visited {
  color: rgba(0, 0, 0, 1) !important;
}

@keyframes fadeInRight {
   0% { opacity: 0; transform: translateX(100px); }
   100% { opacity: 1; transform: translateX(0); }
}
</style>
