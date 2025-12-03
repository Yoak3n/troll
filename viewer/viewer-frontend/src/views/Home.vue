<script lang="ts" setup>
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { NBreadcrumb, NBreadcrumbItem, NCard, NStatistic, NSpace,NDivider } from "naive-ui"
import { fetchDashboardStats, type DashboardStats } from "../api"
const $route = useRoute();
const pathStr = computed(() => {
    const hasQuery = $route.fullPath.includes("?");
    if (hasQuery){
        return $route.fullPath.split("?")[0]!.split("/").filter(i => decodeURI(i)).slice(1)
    }else{
        return $route.fullPath.split("/").filter(i => decodeURI(i)).slice(1);
    }
});
let stats = ref<DashboardStats>({
    topics: 0,
    videos: 0,
    users: 0,
    comments: 0
})


const RouterNameMap = new Map<string,string>([
    ['topics', '话题列表'],
    ['topic', '话题'],
    ['videos', '视频列表'],
    ['video', '视频'],
    ['users', '用户列表'],
    ['user', '用户'],
    ['user-search','用户搜索']
]);
onMounted(async () => {
    console.log('Current path segments:', pathStr.value);
    stats.value = await fetchDashboardStats()
});

</script>


<template>
    <div class="home-page">
        <n-breadcrumb separator=">" v-if="pathStr">
            <n-breadcrumb-item >
                <RouterLink :to="{name: 'home' }">首页</RouterLink>
            </n-breadcrumb-item>
            <n-breadcrumb-item v-for="n in pathStr">
                <RouterLink :to="{name: n }">
                {{RouterNameMap.has(n)?RouterNameMap.get(n):n  }}
                </RouterLink>
            </n-breadcrumb-item>
        </n-breadcrumb>
        <n-divider />
        <div class="home-dashboard" v-if="$route.name == 'home'">
            <n-card title="统计信息">
                <n-space>
                    <n-statistic label="话题总数" :value="stats.topics" />
                    <n-statistic label="视频总数" :value="stats.videos" />
                    <n-statistic label="用户总数" :value="stats.users" />
                    <n-statistic label="评论总数" :value="stats.comments" />
                </n-space>
            </n-card>
        </div>
        <router-view />
    </div>
</template>

<style scoped lang="less">
.home-page {
    padding: 16px;
    height: 100%;
}

</style>