<script lang="ts" setup>
import { h, ref,onMounted, type Component } from 'vue';
import { RouterLink, useRoute } from 'vue-router';
import { NLayout, NLayoutSider, NLayoutContent, NMenu, type MenuOption, NIcon } from 'naive-ui';
import {
    ChatbubbleEllipsesOutline,
    CaretDownOutline,
    CaretForwardCircleSharp,
    HomeOutline,
    PersonOutline,
    BowlingBallOutline,
    CogOutline,
    HeartOutline
} from '@vicons/ionicons5'



const MenuOptions: MenuOption[] = [
    {
        label: '首页', key: 'home', icon: renderIcon(HomeOutline), children: [
            { label: '话题列表', key: 'topics', icon: renderIcon(ChatbubbleEllipsesOutline) },
            { label: '视频列表', key: 'videos', icon: renderIcon(CaretForwardCircleSharp) },
            { label: '用户查询', key: 'user', icon: renderIcon(PersonOutline) }
        ]
    },
    { label: '控制台', key: 'console', icon: renderIcon(BowlingBallOutline) },
    { label: '设置', key: 'setting', icon: renderIcon(CogOutline) },
    { label: '关于', key: "about", icon: renderIcon(HeartOutline) }
];
const $route = useRoute();
const collapsed = ref(true)
const activatedKey = ref('home')

onMounted(() => {
    const routeName = $route.name as string;
    console.log($route.fullPath);
    activatedKey.value = routeName
})

function renderMenuLabel(option: MenuOption) {
    return h(
        RouterLink,
        {
          to: {
            name: option.key as string,
          }
        },
        {default: () => option.label }
    )
}

function renderIcon(icon: Component) {
    return () => h(NIcon, null, { default: () => h(icon) })
}

function expandIcon() {
    return h(NIcon, null, { default: () => h(CaretDownOutline) })
}


</script>

<template>
    <n-layout has-sider style="height: 100vh;">
        <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="240" show-trigger="bar"
            @collapse="collapsed = true" @expand="collapsed = false">
            <n-menu :options="MenuOptions" :render-label="renderMenuLabel" :expand-icon="expandIcon" v-model:value="activatedKey" />
        </n-layout-sider>
        <n-layout-content>
            <router-view />
        </n-layout-content>
    </n-layout>


</template>