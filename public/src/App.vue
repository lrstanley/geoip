<template>
  <div>
    <div class="ui page container">
      <div class="ui stackable grid">
        <div class="row">
          <div class="three wide column navigation">
            <div class="ui secondary stackable vertical pointing menu">
              <a class="header item brand" href="#">GeoIP <i class="blue world icon"></i></a>
              <router-link exact-active-class="active" class="item" :to="{ name: 'lookup' }">Lookup Address <i class="olive search icon"></i></router-link>
              <router-link exact-active-class="active" class="item" :to="{ name: 'bulkLookup' }">Bulk Lookup <i class="yellow database icon"></i></router-link>
              <router-link exact-active-class="active" class="item" :to="{ name: 'apidocs' }">API Documentation <i class="teal book icon"></i></router-link>
              <a href="https://github.com/lrstanley/geoip" class="item">Source on GitHub <i class="purple github icon"></i></a>
            </div>
          </div>
          <div class="ui ten wide column">
            <div class="ui segment main">
              <transition mode="out-in" name="fade" appear>
                <router-view></router-view>
              </transition>
            </div>

            <div class="footer">
              Location data from <a target="_blank" href="http://www.maxmind.com">Maxmind</a> &middot; GeoIP: <a target="_blank" href="https://github.com/lrstanley/geoip">FOSS</a> lookup service, made with <i class="red heart icon"></i>
            </div>
          </div>
        </div>
      </div>
    </div>
    <vue-progress-bar></vue-progress-bar>
  </div>
</template>

<script>
export default {
  name: 'app',
  mounted: function() {
    this.$Progress.finish()
  },
  created: function() {
    this.$Progress.start()

    this.$router.beforeEach((to, from, next) => {
      if (to.meta.progress !== undefined) {
        let meta = to.meta.progress
        this.$Progress.parseMeta(meta)
      }

      this.$Progress.start()

      next()
    })
    this.$router.afterEach((to, from) => {
      this.$Progress.finish()
    })
  }
}
</script>

<style>
.page {
  padding-top: 145px;
}

.page .navigation {
  margin-top: 75px !important;
}

@media screen and (max-width: 1199px) {
  .page {
    padding-top: 15px;
  }
  .page .navigation {
    margin-top: 0 !important;
  }
  .page .navigation .menu {
    width: inherit;
  }
}

.page .main > div {
  padding: 10px;
}

.brand {
  border-bottom: 1px solid #DEDEDE !important;
}

.footer {
  font-size: 13px;
  float: right;
}

.fade-enter-active { transition: opacity .2s; }
.fade-leave-active { transition: opacity .1s; }
.fade-enter, .fade-leave-to { opacity: 0; }

::-webkit-scrollbar {
  width: 10px;
  height: 6px;
}
::-webkit-scrollbar-track-piece {
  background-color: #F5F5F5;
  background-clip: padding-box;
}
::-webkit-scrollbar-thumb {
  background-color: #1678c2;
  background-clip: padding-box;
  border: 2px solid #FFFFFF;
  border-radius: 6px;
}
::-webkit-scrollbar-thumb:window-inactive {
  background-color: #1678c2;
}
</style>
