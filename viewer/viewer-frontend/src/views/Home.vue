<script lang="ts" setup>
import { computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { NBreadcrumb, NBreadcrumbItem } from "naive-ui"

const $route = useRoute();
const pathStr = computed(() => {
    return $route.fullPath.split("/").filter(i => i);
});
onMounted(() => {
    console.log('Current path segments:', pathStr.value);
});

</script>


<template>
    <div class="home-page">
        <n-breadcrumb separator=">" v-if="pathStr.length == 1 && pathStr[0] !== 'home'">
            <n-breadcrumb-item >
                <RouterLink :to="{name: 'home' }">Home</RouterLink>
            </n-breadcrumb-item>
            <n-breadcrumb-item v-for="n in pathStr">
                {{ n }}
            </n-breadcrumb-item>
        </n-breadcrumb>
        <router-view />
    </div>

</template>

<style scoped lang="less">
.home-page {
    padding: 16px;
}

</style>