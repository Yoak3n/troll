<script lang="ts" setup>
import { h, ref,onMounted } from 'vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';
import { NLayout, NLayoutSider, NLayoutContent, NMenu, type MenuOption, NIcon } from 'naive-ui';
import {
    CaretDownOutline,
} from '@vicons/ionicons5'
import { routerNameToMenuOptionMap } from '../assets/map'
import { MenuOptions } from '../assets/option/MenuOption'

const $route = useRoute()
const $router = useRouter()
const collapsed = ref(true)
const activatedKey = ref('home')


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