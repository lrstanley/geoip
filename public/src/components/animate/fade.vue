<template>
  <transition appear mode="out-in">
    <slot></slot>
  </transition>
</template>

<style>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.1s ease-in-out;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
