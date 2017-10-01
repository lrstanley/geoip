<template>
  <div style="padding-bottom: 0;">
    <div class="ui form bulk-form" :class="{ loading: loading, error: error }">
      <div class="field">
        <textarea class="fluid" rows="10" v-model="query" placeholder="List of hostnames, IP's, or a text block that contains IP's; e.g. webserver logs" autofocus></textarea>
      </div>

      <button class="ui right floated submit primary button" @click="lookup"><i class="yellow database icon"></i> Bulk lookup</button>
    </div>

    <div v-if="error" class="ui negative message"><i class="icon warning sign"></i> {{ error }}</div>


    <div v-if="results.length > 0" class="result-box">
        <transition-group name="list" tag="div" class="ui middle relaxed divided list" appear>
          <div v-for="(item, index) in results" class="item" :index="index" :key="item.query">
            <div class="ui right floated content">
              <router-link class="ui mini button" :class="{ loading: !item.data, negative: item.error, disabled: item.error }" :to="{ name: 'lookup', query: { q: item.data.ip } }"><i class="icon search"></i> Lookup</router-link>
            </div>

            <i v-if="item.data.ip && !item.data.country_abbr" class="help circle icon image"></i>
            <i v-if="item.data.country_abbr" :class="[item.data.country_abbr.toLowerCase()]" class="flag image"></i>
            <i v-if="item.error" class="red ban icon image"></i>
            <i v-if="!item.data.ip && !item.error" class="spinner loading icon image"></i>

            <div class="content">
              <span class="header">
                Q: <strong>{{ item.query }}</strong>
                <span v-if="item.data.ip && item.query != item.data.ip">:: <router-link :to="{ name: 'lookup', query: { q: item.data.ip } }">{{ item.data.ip }}</router-link></span>
                <span v-if="item.data.host">:: <router-link :to="{ name: 'lookup', query: { q: item.data.host } }" data-inverted="" :data-tooltip="item.data.host">{{ ellipsis(item.data.host) }}</router-link></span>
              </span>

              <div class="description">
                <span v-if="item.data.summary">{{ item.data.summary }}</span>
                <span v-if="item.error" class="error">{{ item.error }}</span>
                <span v-if="!item.data.ip && !item.error">Loading...</span>
              </div>
            </div>
          </div>
        </transition-group>
    </div>

  </div>
</template>

<script>
export default {
  name: "bulklookup",
  data: function () {
    return {
      query: "",
      error: false,
      loading: false,
      results: [],
    }
  },
  methods: {
    ellipsis: (text) => { return text.length > 45 ? text.substring(0, 43) + '...' : text; },
    copyClipboard: (event) => {
      let clipboard = new Clipboard('.null');
      clipboard.onClick(event);
      toastr.success('Copied to clipboard', '', { timeOut: 1000, preventDuplicates: true });
      clipboard.destroy();
    },
    selectInput: () => {
      setTimeout(function() { $(".bulk-form textarea").focus(); }, 350);
    },
    lookup: function () {
      if (this.query.length == 0 || this.loading) { return; }

      this.error = false;
      this.loading = true;
      this.results = [];
      this.$Progress.start();

      let hostRegex = /^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$/ig;
      let ipRegex = /([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})/g;
      let match;

      let queries = new Set();
      while (match = ipRegex.exec(this.query)) {
        queries.add(match[0]);
      }

      if (queries.size == 0) {
        // Try hostnames.
        var lines = this.query.split("\n");
        for (let i = 0; i < lines.length; i++) {
          match = hostRegex.exec(lines[i]);
          if (match != null) {
            queries.add(match[0]);
          }
        }

        if (queries.size == 0) {
          this.loading = false;
          this.error = "No IP addresses or hostnames found in supplied input!";
          this.$Progress.fail();
          return;
        }
      }

      let promises = [];

      for (let query of queries) {
        let index = this.results.push({ query: query, data: {}, error: false });
        index--;

        promises.push(new Promise((resolve, reject) => {
          this.$http.get(`/api/${query}`).then(response => {
            if (response.body.error != undefined) {
              this.results[index].error = "Error: " + response.body.error.charAt(0).toUpperCase() + response.body.error.slice(1);
            } else {
              this.results[index].data = response.body;
            }

            this.$set(this.results, index, this.results[index]);
            resolve();
          }, response => {
            this.results[index].error = "An unknown exception occurred or service unavailable";
            this.results[index].data = {error: true};

            this.$set(this.results, index, this.results[index]);
            resolve();
          });
        }));
      }

      Promise.all(promises).then((values) => {
        this.loading = false;
        this.$Progress.finish();
        this.query = "";
        this.selectInput();
      });
    }
  },
  mounted: function () {
    this.selectInput();
  }
}
</script>

<style scoped>
.bulk-form { overflow: auto; }
.bulk-form textarea { min-height: 250px; }

.error { color: #db2828; }

.result-box {
  border-top: 1px lightgray solid;
  padding-top: 20px;
  margin-top: 15px;
  margin-left: -25px;
  margin-right: -25px;
  padding: 15px 15px 0 15px;
}

.result-box .image, .result-box .image {padding-top: 10px !important; }
.result-box .icon { padding-right: 11px !important; }

.list-enter, .list-leave-to { opacity: 0; }
.list-enter-active, .list-leave-active {
  animation-duration: .1s;
  animation-name: fadeInRight;
}

@keyframes fadeInRight {
   0% { opacity: 0; transform: translateX(100px); }
   100% { opacity: 1; transform: translateX(0); }
}
</style>
