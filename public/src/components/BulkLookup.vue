<template>
  <div style="padding-bottom: 0;">
    <div class="ui form bulk-form" :class="{ loading: loading, error: error }">
      <div class="field">
        <textarea class="fluid" rows="10" v-model="query" placeholder="List of hostnames, IP's, or a text block that contains IP's; e.g. webserver logs" autofocus></textarea>
      </div>

      <a v-if="stats.rate_remaining" data-position="top left" class="ui label" :data-tooltip="stats.rate_remaining + '/' + stats.rate_limit + ' calls remain; resets in ' + stats.rate_reset_seconds + ' sec'" data-inverted>remaining calls <div class="detail"><strong>{{ stats.rate_remaining.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",") }}</strong></div></a>
      <button class="ui right floated submit primary button" @click="lookup"><i class="yellow database icon"></i> Bulk lookup</button>
    </div>

    <div v-if="error" class="ui negative message"><i class="icon warning sign"></i> {{ error }}</div>

    <div v-if="Object.keys(countries).length > 0" class="result-box">
      <div class="ui list country-list">
        <div v-for="country in countrySort()" :key="country.name" class="item">
          <p>
            <span class="ui teal circular horizontal mini label">{{ country.count }}</span>
            <i v-if="country.name" :class="[country.name.toLowerCase()]" class="flag"></i>
            {{ country.name }}
          </p>

          <div class="ui tiny progress" data-auto-success="false" :data-value="country.count" :data-total="country.top"><div class="bar"></div></div>
        </div>
      </div>
    </div>
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
import { utils } from '../utils'

export default {
  name: "bulklookup",
  mixins: [utils],
  data: function () {
    return {
      query: "",
      error: false,
      loading: false,
      results: [],
      countries: {},
    }
  },
  methods: {
    ellipsis: (text) => { return text.length > 45 ? text.substring(0, 43) + '...' : text; },
    countrySort: function () {
      let out = [];
      let key;
      let top = 0;

      for (key in this.countries) {
        if (top < this.countries[key]) {
          top = this.countries[key];
        }
        out.push({ name: key, count: this.countries[key] });
      }

      out.sort(function(a, b) {
        return a.count == b.count ? 0 : +(a.count < b.count) || -1;
      });

      for (var i = 0; i < out.length; i++) {
        out.top = top;
      }

      $(".ui.list.country-list .ui.tiny.progress").each(function() {
        $(this).progress('update progress', $(this).attr("data-value"));
        $(this).progress('set total', top);
      });

      return out;
    },
    copyClipboard: (event) => {
      let clipboard = new Clipboard('.null');
      clipboard.onClick(event);
      toastr.success('Copied to clipboard', '', { timeOut: 1000, preventDuplicates: true });
      clipboard.destroy();
    },
    lookup: function () {
      if (this.query.length == 0 || this.loading) { return; }

      this.error = false;
      this.loading = true;
      this.results = [];
      this.countries = {};
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
          this.$lookup(query).then(data => {
            this.results[index].data = data;

            this.$set(this.results, index, this.results[index]);

            this.countries[data.country_abbr] = this.countries[data.country_abbr] ? this.countries[data.country_abbr] + 1 : 1;
            resolve();
          }, error => {
            this.results[index].error = error;
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
        this.$autofocus();
      });
    }
  },
  mounted: function () {
    this.$autofocus();
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

.result-box .ui.list.country-list .item {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;
}

.result-box .ui.list.country-list .item p { min-width: 100px; margin-bottom: 0; }
.result-box .ui.list.country-list .item p span { margin-right: 10px; }
.result-box .ui.list.country-list .item .ui.progress {
  flex-grow: 1;
  margin-top: 4px;
  margin-bottom: 2px;
}

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
