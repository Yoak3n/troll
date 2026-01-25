<template>
    <div class="user-search-wrapper">
        <n-tooltip trigger="hover" placement="left">
            <template #trigger>
                <n-float-button position="absolute" shape="circle" :width="60" :height="60" @click="handleOpenSignedUsers" right="30" bottom="50">
                    <n-icon :size="30" color="#4098FC">
                        <Bookmarks />
                    </n-icon>
                </n-float-button>
            </template>
            查看标记用户列表
        </n-tooltip>
        <div class="search-container">
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
                                @update-value="() => searchFormData.searchRangeData = []"
                                v-model:value="searchFormData.searchRange"></n-select>
                        </n-gi>
                        <n-gi :span="2" v-if="searchFormData.searchRange != 'all'">
                            <n-select v-model:value="searchFormData.searchRangeData" :options="searchRangeDataOption"
                                multiple remote clearable filterable :loading="loadingRef"
                                @search="handleQueryClosestRangeData" @focus="handleRangeDataFocused"
                                :clear-filter-after-select="false"
                                :placeholder="searchFormData.searchRange == 'video' ? '请选择视频' : '请选择话题'" />
                        </n-gi>
                    </n-grid>

                </div>

            </div>
        </div>


        <n-drawer v-model:show="active" :width="502">
            <n-drawer-content title="标记用户列表">
                <n-list hoverable clickable>
                    <n-list-item v-for="user in signedUsers" :key="user.uid" @click="handleUserClick(user)">
                        <template #suffix>
                            <n-button type="error" size="small" @click.stop="handleDeleteUser(user.uid)">删除</n-button>
                        </template>
                        <n-thing :title="user.name" content-style="margin-top: 10px;">
                            <template #avatar>
                                <n-avatar :src="user.avatar" />
                            </template>
                            <template #description>
                                <n-space size="small" style="margin-top: 4px">
                                    <n-tag :bordered="false" type="info" size="small">
                                        UID: {{ user.uid }}
                                    </n-tag>
                                    <n-tag :bordered="false" type="success" size="small" v-if="user.location">
                                        {{ user.location }}
                                    </n-tag>
                                </n-space>
                            </template>
                        </n-thing>
                    </n-list-item>
                </n-list>
            </n-drawer-content>
        </n-drawer>
    </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { throttle } from 'lodash-es'
import { NInput, NButton, NSelect, NGrid, NGi, NDrawer, NDrawerContent, NList, NListItem, NThing, NAvatar, NTag, NSpace, NIcon, NFloatButton, NTooltip } from 'naive-ui';
import { Bookmarks } from '@vicons/ionicons5'
import type { SelectOption } from 'naive-ui'
import { ref } from 'vue';
import { fetchDataRange, fetchSignedUsers, unsignUser } from '../api'
import type { SearchOptionRequest, RangeDataType, User } from '../types';

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
const active = ref(false)
const signedUsers = ref<User[]>([])

const handleOpenSignedUsers = async () => {
    try {
        const users = await fetchSignedUsers()
        signedUsers.value = users
        active.value = true
    } catch (error) {
        console.error('Failed to fetch signed users:', error)
        window.$message?.error('获取标记用户列表失败')
    }
}

const handleDeleteUser = async (uid: number) => {
    try {
        await unsignUser([uid])
        signedUsers.value = signedUsers.value.filter(u => u.uid !== uid)
        window.$message?.success('删除成功')
    } catch (error) {
        console.error('Failed to delete signed user:', error)
        window.$message?.error('删除失败')
    }
}

const handleUserClick = (user: User) => {
    $router.push({
        name: 'user',
        query: {
            uid: user.uid
        }
    })
    active.value = false
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
let searchRangeData: RangeDataType[] = []
const searchRangeDataOption = ref<SelectOption[]>([])
const loadingRef = ref(false)
const fetchClosestRangeData = throttle(async (): Promise<RangeDataType[]> => fetchSpecityRangeData(searchFormData.value.searchType, searchFormData.value.search, searchFormData.value.searchRange), 500)
const fetchSpecityRangeData = async (typ: string, data: string, rangeType: string): Promise<RangeDataType[]> => {
    let searchOption: SearchOptionRequest = {
        rangeType
    }
    if (data != "") {
        if (typ != 'uid') {
            searchOption.name = data
        } else {
            const uid = parseInt(data)
            searchOption.uid = uid
        }
    }
    const ret = await fetchDataRange(searchOption)
    if (ret.type == rangeType) {
        return ret.options
    }
    return []
}
const handleRangeDataFocused = async () => {
    searchRangeData = await fetchClosestRangeData()
    searchRangeDataOption.value = searchRangeData as { label: string, value: string }[]
}
const handleQueryClosestRangeData = throttle(async (query: string) => {
    console.log("search.....");

    loadingRef.value = true
    try {
        const options = searchRangeData as { label: string, value: string }[]
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
    $router.push({
        name: 'user', query: {
            uid: searchFormData.value.searchType == 'uid' ? parseInt(searchFormData.value.search) : undefined,
            name: searchFormData.value.searchType == 'username' ? searchFormData.value.search : undefined,
            rangeType: searchFormData.value.searchRange,
            rangeData: searchFormData.value.searchRangeData.join(",")
        }
    })
}


</script>

<style scoped lang="less">
.user-search-wrapper {
    height: 100%;
    width: 100%;

}

.search-container {
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