<template>
  <div>
    <h2 class="ui header">
      <i class="settings icon" />
      <div class="content">
        API Documentation
        <div class="sub header">Query the GeoIP API through <strong>REST</strong> and <strong>JSON</strong></div>
      </div>
    </h2>

    <div class="ui divider" />

    <h3 class="ui dividing header">
      <a id="api-endpoint" href="#api-endpoint">#</a> API Endpoint
      <div class="sub header">Where/how do I query the API?</div>
    </h3>
    <p>
      The HTTP API takes <code>GET</code> requests in the following scheme:
      <pre class="block"><code>http[s]://geoip.pw/api{/endpoint}</code></pre>
    </p>

    <h3 class="ui dividing header">
      <a id="lookup-query" href="#lookup-query">#</a> Lookup an IP address or hostname
      <div class="sub header">How do I obtain the geographical info for a hostname or IP?</div>
    </h3>
    <p>
      The GeoIP API supports both IP (IPv4 and IPv6) lookups, as well as
      hostname/domain lookups. To obtain JSON information about a given IP
      address or hostname, simply use the following format (replacing
      <code class="inline">{query}</code> with the IP or hostname):
      <pre class="block"><code>https://geoip.pw/api/{ip}</code></pre>
      <br>

      For example:
      <pre class="block"><code>$ curl https://geoip.pw/api/8.8.8.8
{"ip":"8.8.8.8","summary":"United States, NA","city":"","subdivision":"","country":"United States","country_abbr":"US","continent":"North America","continent_abbr":"NA","latitude":37.751,"longitude":-97.822,"timezone":"America/Chicago","postal_code":"","proxy":false,"host":"dns.google"}</code></pre>
      <br>

      We can take that one step further, and prettify the JSON:
      <pre class="block"><code>$ curl https://geoip.pw/api/8.8.8.8?pretty=1
{
  "ip": "8.8.8.8",
  "summary": "United States, NA",
  "city": "",
  "subdivision": "",
  "country": "United States",
  "country_abbr": "US",
  "continent": "North America",
  "continent_abbr": "NA",
  "latitude": 37.751,
  "longitude": -97.822,
  "timezone": "America/Chicago",
  "postal_code": "",
  "proxy": false,
  "host": "dns.google"
}
</code></pre>
      <br>
      GeoIP also supports providing a <code class="inline">lang</code> query
      parameter (or <code class="inline">Accept-Language</code> header) to
      specify the language of the response:
      <pre class="block"><code>$ curl "https://geoip.pw/api/8.8.8.8?pretty=1&lang=zh-CN"
{
  "ip": "8.8.8.8",
  "summary": "美国, NA",
  "city": "",
  "subdivision": "",
  "country": "美国",
  "country_abbr": "US",
  "continent": "北美洲",
  "continent_abbr": "NA",
  "latitude": 37.751,
  "longitude": -97.822,
  "timezone": "America/Chicago",
  "postal_code": "",
  "proxy": false,
  "host": "dns.google"
}
</code></pre>
    </p>

    <h3 class="ui dividing header">
      <a id="filters" href="#filters">#</a> Filtering output
      <div class="sub header">What if I only want a specific field to reduce query time?</div>
    </h3>
    <p>
      The query API supports filtering and obtaining only specific fields. The
      supported fields are what you currently see in the JSON output above.
      The full lookup endpoint does things like reverse DNS lookups, so using
      this, will allow you to speed up potential queries if not all the
      information is being used.
      <br><br>

      The syntax is as follows:
      <pre class="block"><code>http[s]://geoip.pw/api/{query}/{comma-list-of-fields}</code></pre>
</p><div class="ui warning message">
        <strong>Warning:</strong> the <code class="inline">{comma-list-of-fields}</code>
        must be a set of valid field parameters, or the API will return an error.
        Errors returned by this endpoint will be plaintext (e.g.
        <code class="inline">error: too many filters supplied</code>), which
        should hopefully cater to those using this endpoint for shell scripts.
      </div>
      <br>

      The response to this query will be a pipe (<code class="inline">|</code>) separated
      array of values. Pipe is used primarily as it's very unlikely returned
      fields themselves will contain pipe symbols (compared to commas, for
      example), which allows things like shell scripts to easily parse the
      resulting output.
      <br><br>

      For example:
      <pre class="block"><code>$ curl https://geoip.pw/api/8.8.8.8/summary,country,host
United States, NA|United States|google-public-dns-a.google.com</code></pre>
    </p>

    <h3 class="ui dividing header">
      <a id="headers" href="#headers">#</a> Additional headers
      <div class="sub header">How much of a rate limit do I have? How old is the GeoIP data?</div>
    </h3>
    <p>
      With each API response, the API will tag a few headers that will help
      you better determine what your current rate limit is. If you want to
      obtain these headers without submitting a query for a given IP/hostname,
      simply use the below endpoint, which will not count towards your API
      limits:
      <pre class="block"><code>http[s]://geoip.pw/api/ping</code></pre>
</p><h4>Access-Control headers:</h4>
      GeoIP also has support to restrict connections based on the
      <code class="inline">Access-Control-Allow-Origin</code> header, as long
      as end browsers/clients support it. By default, this isn't enabled,
      however self-hosted versions may utilize this.
      <br><br>

      The following are the currently supported headers:
      <table class="ui small definition table">
        <thead>
          <tr>
            <th class="collapsing" />
            <th class="collapsing" data-tooltip="This indicates that the header will always be present in the given response">Always Present</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td class="collapsing"><code class="inline">X-Maxmind-Build</code></td>
            <td class="center aligned collapsing"><i class="large green checkmark icon" /></td>
            <td>The <a href="http://www.maxmind.com/" target="_blank">Maxmind</a> database version</td>
          </tr>
          <tr>
            <td class="collapsing"><code class="inline">X-Maxmind-Type</code></td>
            <td class="center aligned collapsing"><i class="large green checkmark icon" /></td>
            <td>The <a href="http://www.maxmind.com/" target="_blank">Maxmind</a> db type -- see <a href="https://dev.maxmind.com/geoip/geoip2/geolite2/" target="_blank">this page</a> for more info</td>
          </tr>
          <tr>
            <td class="collapsing"><code class="inline">X-Ratelimit-Limit</code></td>
            <td class="collapsing" />
            <td>Number of queries allowed within <code class="inline">reset</code> period</td>
          </tr>
          <tr>
            <td class="collapsing"><code class="inline">X-Ratelimit-Remaining</code></td>
            <td class="collapsing" />
            <td>Number of queries remaining within <code class="inline">limit</code></td>
          </tr>
          <tr>
            <td class="collapsing"><code class="inline">X-Ratelimit-Reset</code></td>
            <td class="collapsing" />
            <td>Seconds until <code class="inline">remaining</code> is reset</td>
          </tr>
        </tbody>
      </table>

      <span v-if="headers.length > 0">
        <br>
        The current values of these headers:
        <br>
        <table class="ui small striped celled table">

          <tbody>
            <tr v-for="header in headers" :key="header.name">
              <td class="collapsing"><code class="inline">{{ header.name }}</code></td>
              <td>{{ header.value[0] }}</td>
            </tr>
          </tbody>
        </table>
      </span>

      <div class="ui warning message">
        <strong>Warning:</strong> not all self-hosted versions of GeoIP may have
        rate-limiting enabled, and as such, the header may not always be present.
        (see the "supported headers" table, showing which headers always will be
        present.)
      </div>
      <br>

      An example response may look like this:
      <pre class="block"><code>$ curl -i https://geoip.pw/api/ping
HTTP/1.1 200 OK
Content-Type: application/json
X-Maxmind-Build: 6-1501721259
X-Maxmind-Type: GeoLite2-City
X-Ratelimit-Limit: 2000
X-Ratelimit-Remaining: 1993
X-Ratelimit-Reset: 3593

{"pong":true}</code></pre>
      <div class="ui info message">
        <strong>Note:</strong> the <code class="inline">/api/ping</code>
        endpoint is the only one which currently supports HTTP <code class="inline">HEAD</code>
        requests (which still contain the necessary headers.)
      </div>
    </p>
  </div>
</template>

<script>
export default {
  name: "apidocs",
  data: function () {
    return { headers: [] }
  },
  mounted: function() {
    this.$http.head("/api/ping").then((response) => {
      let out = [];

      for (name in response.headers.map) {
        let nice = "";
        for (let i = 0; i < name.length; i++) {
          if (i == 0 || name[i - 1] == "-") {
            nice += name[i].toUpperCase();
            continue;
          }

          nice += name[i];
        }

        if (name.includes("x-maxmind") || name.includes("x-rate")) {
          out.push({ name: nice, value: response.headers.map[name] })
        }
      }

      out.sort(function(a, b) {
        return a.name == b.name ? 0 : +(a.name > b.name) || -1;
      });

      this.headers = out;
    }, function(error) {});
  }
}
</script>

<style scoped>
code.inline {
  font-family: "Courier New", Courier, monospace;
  background-color: #e8e8e8;
  color: #383838;
  padding: 3px 4px;
  border-radius: 4px;
}
pre.block {
  background-color: #383838;
  color: white;
  padding: 10px;
  border-radius: 6px;
}
pre.block > code {
  font-family: "Courier New", Courier, monospace;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.ui.dividing.header { margin-top: 40px; }
</style>
