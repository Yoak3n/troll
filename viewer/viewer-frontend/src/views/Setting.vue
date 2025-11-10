<template>
    <div class="setting-wrapper">
        <n-card title="账号设置">
            <div class="cookie-setting" v-for="(item, index) in SettingData.cookies">
                <div class="setting-item">
                    <n-input placeholder="请输入账号cookie" v-model:value="item.data" @blur="UpdateSetting"></n-input>
                    <n-button-group>
                        <n-button @click="() => SettingData.cookies.splice(index, 1)">
                            <n-icon>
                                <Remove />
                            </n-icon>
                        </n-button>
                        <n-button v-if="index == (SettingData.cookies.length - 1)"
                            @click="() => {SettingData.cookies.push({type: 'cookie',data:''});UpdateSetting()}" round>
                            <n-icon>
                                <Add />
                            </n-icon>
                        </n-button>
                    </n-button-group>
                </div>
            </div>
        </n-card>
        <n-card title="代理设置">
            <div class="proxy-setting" v-for="(item, index) in SettingData.proxies">
                <div class="setting-item">
                    <n-input placeholder="请输入代理地址" v-model:value="item.data" @blur="UpdateSetting"></n-input>
                    <n-button-group>
                        <n-button @click="() => {SettingData.proxies.splice(index, 1);UpdateSetting()}">
                            <n-icon>
                                <Remove />
                            </n-icon>
                        </n-button>
                        <n-button v-if="index == (SettingData.proxies.length - 1)"
                            @click="() => SettingData.proxies.push({type: 'proxy',data:''})" round>
                            <n-icon>
                                <Add />
                            </n-icon>
                        </n-button>
                    </n-button-group>
                </div>
            </div>
        </n-card>
        <n-card title="请求策略">
            <n-grid :cols="14" :x-gap="20">
                <n-gi span="1">
                    <div class="basic">
                        基础间隔
                    </div>
                </n-gi>
                <n-gi span="6">
                    <n-input-number  :default-value="1"
                        disabled
                        v-model:value="SettingData.interval.basic">
                        <template #suffix>
                            秒
                        </template>
                    </n-input-number>
                </n-gi>
                <n-gi span="1">
                    <div class="random">
                        随机范围
                    </div>
                </n-gi>
                <n-gi span="6">
                    <n-input-number :default-value="3"
                        disabled
                        v-model:value="SettingData.interval.random">
                        <template #suffix>
                            秒
                        </template>
                    </n-input-number>
                </n-gi>
            </n-grid>


        </n-card>
    </div>
</template>


<script setup lang="ts">
import { onMounted,reactive } from 'vue';
import { NCard, NGrid, NGi, NInputNumber, NInput, NButtonGroup, NButton, NIcon } from 'naive-ui';
import { Add, Remove } from '@vicons/ionicons5'
import type { Setting } from '../types';
import { fetchSetting,updateSetting } from '../api';

const SettingData = reactive<Setting>({
    cookies: [{type: 'cookie',data: ''}],
    proxies: [{type: 'proxy',data:''}],
    interval: {
        basic: 1,
        random: 3
    }
})

onMounted(()=>{
    syncSetting()
})

const syncSetting =  async() =>{
    const setting = await fetchSetting()
  
    SettingData.cookies = setting.cookies
    SettingData.proxies = setting.proxies
    if (SettingData.cookies.length == 0){
        SettingData.cookies.push({type: 'cookie',data: ''})
    }
    if (SettingData.proxies.length == 0){
        SettingData.proxies.push({type: 'proxy',data: ''})
    }
}

const UpdateSetting  = () =>{
    console.log(SettingData)
    const uploadData: Setting = {
        cookies: SettingData.cookies.filter((v)=>v.data != ''),
        proxies: SettingData.proxies.filter(v=>v.data != ''),
        interval: SettingData.interval
    } 
    updateSetting(uploadData)
}


</script>

<style scoped lang="less">
.setting-item {
    display: flex;
    width: 80%;
}
</style>