<template>
    <div class="user-search-wrapper">
        <div class="search-bar">
            <h1 style="text-align: center;">搜索指定用户</h1>
            <div class="search-actions">
                <n-input v-model:value="searchFormData.search" size="large" placeholder="请输入用户的对应信息"></n-input>
                <n-button type="info" size="large" @click="handleSearchButtonClicked">搜索</n-button>
            </div>
            <div class="search-options">
                <n-grid :cols="2" :y-gap="2">
                    <n-gi>
                        <n-select :options="searchTypeOptions" :default-value="initSearchFormData.searchType"
                            v-model:value="searchFormData.searchType"></n-select>
                    </n-gi>
                    <n-gi>
                        <n-select :options="searchRangeOption" :default-value="initSearchFormData.searchRange"
                            @update-value="()=>searchFormData.searchRangeData=[]"
                            v-model:value="searchFormData.searchRange"></n-select>
                    </n-gi>
                    <n-gi :span="2" v-if="searchFormData.searchRange != 'all'">
                        <n-select v-model:value="searchFormData.searchRangeData" :options="searchRangeDataOption"
                            multiple remote clearable filterable
                            :loading="loadingRef" 
                            @search="handleQueryClosestRangeData" 
                            @focus="handleRangeDataFocused"
                            :clear-filter-after-select="false"
                            :placeholder="searchFormData.searchRange == 'video' ? '请选择视频' : '请选择话题'" />
                    </n-gi>
                </n-grid>

            </div>

        </div>
    </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { throttle } from 'lodash-es'
import { NInput, NButton, NSelect, NGrid, NGi } from 'naive-ui';
import type { SelectOption } from 'naive-ui'
import { ref } from 'vue';
import {fetchDataRange} from '../api'
import type { SearchOptionRequest,RangeDataType } from '../types';

interface SearchForm {
    search: string,
    searchType: string,
    searchRange: string,
    searchRangeData: string[]
}

const initSearchFormData = {
    search: '',
    searchType: 'uid',
    searchRange: 'all',
    searchRangeData: []
}
const searchFormData = ref<SearchForm>(initSearchFormData)
const searchTypeOptions: SelectOption[] = [
    {
        label: '按UID搜索',
        value: 'uid'
    }, {
        label: '按用户名搜索',
        value: 'username'
    }
]
const searchRangeOption: SelectOption[] = [
    {
        label: '不限范围',
        value: 'all'
    }, {
        label: '指定话题',
        value: 'topic'
    }, {
        label: '指定视频',
        value: 'video'
    }
]
let searchRangeData:RangeDataType[] = []
const searchRangeDataOption = ref<SelectOption[]>([])
const loadingRef = ref(false)
const fetchClosestRangeData = throttle(async(): Promise<RangeDataType[]> => fetchSpecityRangeData(searchFormData.value.searchType,searchFormData.value.search,searchFormData.value.searchRange),500)
const fetchSpecityRangeData = async (typ: string, data: string, rangeType: string): Promise<RangeDataType[]>  => {
    let searchOption:SearchOptionRequest = {
        rangeType
    }
    if (data != "") {
        if (typ != 'uid') {
            searchOption.name = data
        }else{
            const uid = parseInt(data)
            searchOption.uid = uid
        }
    }
    const ret = await fetchDataRange(searchOption)
    if (ret.type == rangeType){
        return ret.options
    }
    return []
}
const handleRangeDataFocused = async() => {
    searchRangeData = await fetchClosestRangeData()
    searchRangeDataOption.value = searchRangeData as {label: string, value: string}[]
}
const handleQueryClosestRangeData =  throttle(async (query: string) => {
    console.log("search.....");
    
    loadingRef.value = true
    try {
        const options = searchRangeData as {label: string, value: string}[]
        searchRangeDataOption.value = options.filter(item => item.label.toLowerCase().includes(query.toLowerCase()))
    } catch (error) {
        console.error('查询失败:', error)
        searchRangeDataOption.value = []
    } finally {
        loadingRef.value = false
    }
}, 500)
const $router = useRouter()
const handleSearchButtonClicked = () => {
    $router.push({name: 'user', query: {
        uid: searchFormData.value.searchType == 'uid' ? parseInt(searchFormData.value.search):undefined, 
        name: searchFormData.value.searchType == 'username' ? searchFormData.value.search: undefined,
        rangeType: searchFormData.value.searchRange,
        rangeData: searchFormData.value.searchRangeData.join(",")
    }})
}


</script>

<style scoped lang="less">
.user-search-wrapper {
    height: 100%;
    width: 100%;
    justify-content: center;
    align-items: center;
    display: flex;

    .search-bar {
        width: 50%;
        height: 50%;
        .search-actions {
            display: flex;
        }
        .search-options {
            margin-top: .5rem;
            padding: 0 .5rem;
        }
    }
}
</style>