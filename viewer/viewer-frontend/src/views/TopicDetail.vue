<template>
    <div class="topic-wrapper" v-if="$route.name == 'topic'">
        <n-flex justify="space-between" align="center">
            <h1>Topic: {{ topicName }}</h1>
            <n-flex class="action" align="center">
                <div class="action-button" v-if="Object.keys(checkedVideos).length > 0">
                    <n-button type="error" @click="deleteSelectedVideos" >删除选中视频</n-button>
                    <n-button type="info" @click="moveToNewTopic">修改话题分类</n-button>
                </div>
                <n-switch v-model:value="modifyMode" size="large">
                    <template #checked>
                        取消
                    </template>
                    <template #unchecked>
                        <p style="color: #444;">
                            编辑
                        </p>
                    </template>
                </n-switch>
            </n-flex>

        </n-flex>

        <n-grid x-gap="12" :cols="3">
            <n-gi v-for="video in videos" :key="video.avid" v-if="videos.length > 0">
                <n-card :title="EllipsisText(video.title, 15)" tag="div" hoverable content-class="video-card-content"
                    size="medium"
                    :content-style="{ 'display': 'flex', 'justify-content': 'end', 'max-height': '64rem' }">
                    <template #header-extra v-if="modifyMode">
                        <n-space>
                            <n-checkbox size="large" v-model:checked="checkedVideos[video.avid]" />
                        </n-space>
                    </template>
                    <router-link :to="{ name: 'video', query: { 'avid': video.avid, 'topicName': video.topic } }"
                        style="color: #1797FF;">包含{{
                            video.count }}条评论
                    </router-link>
                    <template #footer>
                        <NFlex justify="space-between">
                            <NFlex align="center">
                                <NAvatar :src="video.author.avatar" alt="author avatar" round />
                                {{ video.author.name }}
                            </NFlex>
                            <NFlex align="center">
                                更新时间：{{ video.update_at }}
                            </NFlex>
                        </NFlex>
                    </template>
                </n-card>

            </n-gi>
        </n-grid>
    </div>
    <RouterView />
</template>

<script setup lang="ts">
import { h, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { NGrid, NGi, NCard, NAvatar, NFlex, NSwitch, NCheckbox, NButton } from 'naive-ui';

import { deleteVideos, fetchVideosByTopic,updateTopicOfVideos } from '../api';
import type { VideoDataWithCommentsCount, TopicUpdateRequest } from '../types';
import { EllipsisText } from '../utils/name/show';

const $route = useRoute();
const topicName = ref<string>('');
const videos = ref<Array<VideoDataWithCommentsCount>>([]);
const isLoading = ref<boolean>(false);
const modifyMode = ref<boolean>(false);
const checkedVideos = ref<Record<number, boolean>>({});

onMounted(async () => {
    const topic = $route.query.topicName as string;
    topicName.value = topic
    await loadVideos();
});

const loadVideos = async () => {
    isLoading.value = true
    window.$loadingBar?.start()
    const ret = await fetchVideosByTopic(topicName.value);
    isLoading.value = false
    videos.value = ret;
    window.$loadingBar?.finish()
}

const deleteSelectedVideos = async () => {
    if (!modifyMode.value) {
        return;
    }
    const selectedVideos: number[] = Object.keys(checkedVideos.value).filter(avid => checkedVideos.value[Number(avid)]).map(avid => Number(avid));
    if (selectedVideos.length === 0) {
        window.$message?.warning('请选择要删除的视频');
        return;
    }
    window.$dialog?.warning({
        title: '确认删除选中视频吗？',
        content: `确定删除选中的 ${selectedVideos.length} 个视频吗？`,
        positiveText: '确认',
        negativeText: '取消',
        onPositiveClick: async () => {
            try {
                await deleteVideos(selectedVideos);
                window.$message?.success('删除成功');
                // 刷新视频列表
                loadVideos();
                // 清空选中状态
                checkedVideos.value = {};
            } catch (error) {
                window.$message?.error('删除失败');
            }
        }
    });
};

const newTopicName = ref<string>('');
import TopicDialogue from '../components/Topic/TopicDialogue.vue'
const moveToNewTopic = async () => {
    if (!modifyMode.value) {
        return;
    }
    const selectedVideos: number[] = Object.keys(checkedVideos.value).filter(avid => checkedVideos.value[Number(avid)]).map(avid => Number(avid));
    if (selectedVideos.length === 0) {
        window.$message?.warning('请选择要修改的视频');
        return;
    }
    window.$dialog?.create({
        title: '请输入新的话题名称',
        content: ()=>h(TopicDialogue, {
            modelValue: newTopicName.value,
            'onUpdate:modelValue': (val: string | undefined) => newTopicName.value = val || '',
            placeholder: '输入新的话题名称'
        }),
        positiveText: '确认',
        negativeText: '取消',
        onPositiveClick: async () => {
            if (newTopicName.value !== '') {
                const req:TopicUpdateRequest = {
                    avid: selectedVideos,
                    topic: newTopicName.value
                }
                await updateTopicOfVideos(req);
                window.$message?.success('修改成功');
                // 刷新视频列表
                loadVideos();
                // 清空选中状态
                checkedVideos.value = {};
            }else{
                window.$message?.error('请输入话题名称')
            }
        }
    });
}


</script>