<template>
  <div>
    <div class="ui page container">
      <div class="ui grid">
        <div class="row">
          <div class="three wide column navigation">
            <div class="ui secondary vertical pointing menu">
              <a class="header item brand" href="#">GeoIP <i class="blue world icon"></i></a>
              <router-link exact-active-class="active" class="item" :to="{ name: 'lookup'}">Lookup Address <i class="olive search icon"></i></router-link>
              <router-link exact-active-class="active" class="item" :to="{ name: 'apidocs'}">API Documentation <i class="teal book icon"></i></router-link>
              <router-link exact-active-class="active" class="item" :to="{ name: 'about'}">About <i class="red help circle icon"></i></router-link>
            </div>
          </div>
          <div class="ui ten wide column">
            <div class="ui segment main">
              <transition mode="out-in" name="fade" appear>
                <router-view></router-view>
              </transition>
            </div>

            <!-- TODO: footer here -->
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
// TODO:
//   - https://github.com/hilongjw/vue-progressbar
//   - http://whois.arin.net/rest/ip/8.8.8.8.json
export default {
  name: 'app',
  // mounted: function () {
  //   console.log("mounted")
  //   this.$http.get('/api/8.8.8.8').then(response => {
  //     console.log(response)
  //     console.info(response.body.ip)
  //   }, response => {
  //     console.log(response)
  //   });
  // }
}
</script>

<style>
.page {
  padding-top: 145px;
}

.page .navigation {
  margin-top: 75px !important;
}

.page .main > div {
  padding: 10px;
}

.brand {
  border-bottom: 1px solid #DEDEDE !important;
}

.fade-enter-active { transition: opacity .3s; }
.fade-leave-active { transition: opacity .2s; }
.fade-enter, .fade-leave-to { opacity: 0; }
</style>
