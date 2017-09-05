<template>
  <div class="hello">
    <img class="logo" src="../assets/logo.png" />
    <h1>{{ msg }}</h1>
  </div>
</template>

<script>
export default {
  data () {
    return {
      msg: 'Hello World!'
    }
  }
}
</script>

<style scoped lang="less">
h1 {
  color: red;
}
</style>
