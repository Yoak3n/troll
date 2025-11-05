<script lang="ts" setup>
import { h, ref,onMounted, type Component } from 'vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';
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
            { label: '用户查询', key: 'users', icon: renderIcon(PersonOutline) }
        ]
    },
    { label: '控制台', key: 'console', icon: renderIcon(BowlingBallOutline) },
    { label: '设置', key: 'setting', icon: renderIcon(CogOutline) },
    { label: '关于', key: "about", icon: renderIcon(HeartOutline) }
];
const $route = useRoute()
const $router = useRouter()
const collapsed = ref(true)
const activatedKey = ref('home')

const routerNameToMenuOptionMap: Map<string,string> = new Map<string,string>([
    ['user-search','users'],
    ['user','users']
])

onMounted(() => {
    const routeName = $route.name as string
    UpdateActivatedKey(routeName)
    $router.afterEach((to)=>{
        const name = to.name as string;    
        UpdateActivatedKey(name)
    })
})

const UpdateActivatedKey = (name:string)=>{
    activatedKey.value = routerNameToMenuOptionMap.has(name) ? routerNameToMenuOptionMap.get(name)!: name
}



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
        <n-layout-sider 
            bordered 
            collapse-mode="width" 
            :collapsed-width="52" 
            :width="240" 
            show-trigger="bar"
            :collapsed="collapsed"
            @collapse="collapsed = true" @expand="collapsed = false">
            <n-menu :options="MenuOptions" :render-label="renderMenuLabel" :expand-icon="expandIcon" v-model:value="activatedKey" />
        </n-layout-sider>
        <n-layout-content>
            <router-view />
        </n-layout-content>
    </n-layout>


</template>