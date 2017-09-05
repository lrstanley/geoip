<template>
  <div>
    <div class="ui massive fixed inverted borderless menu">
      <div class="ui container">
        <a class="header item" href="#">GeoIP</a>
        <router-link exact-active-class="active" class="item" :to="{ name: 'hello'}">Hello</router-link>
        <router-link exact-active-class="active" class="item" :to="{ name: 'counter'}">Counter</router-link>
      </div>
    </div>
    <div class="ui main container">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
export default {
  name: 'app',
}
</script>

<style>
.main {
  margin-top: 50px;
}
</style>
