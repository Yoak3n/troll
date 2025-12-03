<template>
    <div class="topic-wrapper" v-if="$route.name == 'topic'">
        <h1>Topic: {{ topicName }}</h1>
        <n-grid x-gap="12" :cols="3">
            <n-gi v-for="video in videos" :key="video.avid" v-if="videos.length > 0">
                <n-card :title="EllipsisText(video.title, 15)" tag="div" hoverable content-class="video-card-content"
                    size="medium"
                    :content-style="{ 'display': 'flex', 'justify-content': 'end', 'max-height': '64rem' }">
                    <router-link :to="{ name: 'video', query: { 'avid': video.avid, 'topicName': video.topic } }" style="color: #1797FF;">包含{{
                        video.count}}条评论
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
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { NGrid, NGi, NCard, NAvatar, NFlex } from 'naive-ui';
import { fetchVideosByTopic } from '../api';
import type { VideoDataWithCommentsCount } from '../types';
import { EllipsisText } from '../utils/name/show';
const $route = useRoute();
const topicName = ref<string>('');
const videos = ref<Array<VideoDataWithCommentsCount>>([]);
const isLoading = ref<boolean>(false);
onMounted(async () => {
    const topic = $route.query.topicName as string;
    topicName.value = topic
    isLoading.value = true
    window.$loadingBar?.start()
    const ret = await fetchVideosByTopic(topicName.value);
    window.$loadingBar?.finish()
    isLoading.value = false
    videos.value = ret;
});


</script>