export const utils = {
  data: function () {
    return {
      stats: {
        db_build: null,
        db_type: null,
        rate_limit: null,
        rate_remaining: null,
        rate_reset_seconds: null,
      }
    }
  },
  methods: {
    $autofocus: function () {
      setTimeout(function () { $("[autofocus]").focus(); }, 350);
    },
    $updateStats: function (httpResponse) {
      this.stats.db_build = httpResponse.headers.map["x-maxmind-build"][0];
      this.stats.db_type = httpResponse.headers.map["x-maxmind-type"][0];
      this.stats.rate_limit = httpResponse.headers.map["x-ratelimit-limit"][0];
      this.stats.rate_remaining = httpResponse.headers.map["x-ratelimit-remaining"][0];
      this.stats.rate_reset_seconds = httpResponse.headers.map["x-ratelimit-reset"][0];
    },
    $ping: function () {
      this.$http.get(`/api/ping`).then(response => { this.$updateStats(response); });
    },
    $lookup: function (address) {
      return new Promise((resolve, reject) => {
        this.$http.get(`/api/${address}`).then(response => {
          this.$updateStats(response);

          if (response.body.error != undefined) {
            if (address == 'self') {
              // Don't show a nasty error if we can't even look up their own
              // IP address on page load.
              reject(null);
              return;
            }

            reject("Error: " + response.body.error.charAt(0).toUpperCase() + response.body.error.slice(1));
            return;
          }

          // Add our query into the result, so when it gets saved to history,
          // we can use it later.
          response.body.query = address;

          resolve(response.body);
        }, () => {
          reject("An unknown exception occurred or service unavailable");
          return
        });
      });
    }
  }
}
