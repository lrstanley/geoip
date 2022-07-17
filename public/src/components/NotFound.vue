<template>
    <div class="not-found">
        <h2 class="ui center aligned dividing icon header">
            <i class="circular search icon" />
            Uhoh! Page not found!
        </h2>
        <!-- <h2 class="ui dividing header">Uhoh! Page not found!</h2> -->

        <p>
            The page you were looking for doesn't exist. Maybe try heading back to the <router-link :to="{ name: 'lookup' }">home page</router-link>, or the
            <router-link :to="{ name: 'apidocs' }">API Documentation</router-link>?
        </p>
    </div>
</template>

<script></script>

<style scoped>
.not-found {
    padding-top: 70px !important;
    padding-bottom: 70px !important;
}
</style>
