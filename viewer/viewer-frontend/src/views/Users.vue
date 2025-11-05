<template>
    <div class="users-wrapper">
            <div v-if="$route.name=='users'">
        <h1>users</h1>
    </div>
    <router-view/>
    </div>

</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';

const $route = useRoute()

</script>

<style scoped>
.users-wrapper{
    height: 100%;
}
</style>