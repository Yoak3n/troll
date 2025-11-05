<script lang="ts" setup>
import { computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { NBreadcrumb, NBreadcrumbItem } from "naive-ui"

const $route = useRoute();
const pathStr = computed(() => {
    const hasQuery = $route.fullPath.includes("?");
    if (hasQuery){
        return $route.fullPath.split("?")[0]!.split("/").filter(i => decodeURI(i)).slice(1)
    }else{
        return $route.fullPath.split("/").filter(i => decodeURI(i)).slice(1);
    }
});
const RouterNameMap = new Map<string,string>([
    ['topics', '话题列表'],
    ['topic', '话题'],
    ['video', '视频'],
    ['users', '用户列表'],
    ['user', '用户'],
    ['user-search','用户搜索']
]);
onMounted(() => {
    console.log('Current path segments:', pathStr.value);
});

</script>


<template>
    <div class="home-page">
        <n-breadcrumb separator=">" v-if="pathStr && pathStr.length != 1">
            <n-breadcrumb-item >
                <RouterLink :to="{name: 'home' }">首页</RouterLink>
            </n-breadcrumb-item>
            <n-breadcrumb-item v-for="n in pathStr">
                <RouterLink :to="{name: n }">
                {{RouterNameMap.has(n)?RouterNameMap.get(n):n  }}
                </RouterLink>
            </n-breadcrumb-item>
        </n-breadcrumb>
        <router-view />
    </div>

</template>

<style scoped lang="less">
.home-page {
    padding: 16px;
    height: 100%;
}

</style>