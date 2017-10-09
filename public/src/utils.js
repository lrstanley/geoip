export const utils = {
  data: function () {
    return {
      // TODO.
    }
  },
  methods: {
    $autofocus: function () {
      setTimeout(function () { $("[autofocus]").focus(); }, 350);
    },
    $lookup: function (address) {
      return new Promise((resolve, reject) => {
        let error;
        this.$http.get(`/api/${address}`).then(response => {
          if (response.body.error != undefined) {
            if (query == 'self') {
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
          console.log(response.body);
        }, response => {
          reject("An unknown exception occurred or service unavailable");
          return
        });
      });
    }
  }
}
