<template>
  <div>
    <div class="ui page container">
      <div class="ui stackable grid">
        <div class="row">
          <div class="three wide column navigation">
            <div class="ui secondary stackable vertical pointing menu">
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

.fade-enter-active { transition: opacity .3s; }
.fade-leave-active { transition: opacity .2s; }
.fade-enter, .fade-leave-to { opacity: 0; }
</style>
