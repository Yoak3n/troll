<template>
    <div class="topic-wrapper">
        <h1>Topic: {{ topicName }}</h1>
        <n-grid x-gap="12" :cols="2">
            <n-gi v-for="video in videos" :key="video.avid" v-if="videos.length > 0">
                <n-card 
                :title="EllipsisText(video.title, 15)" 
                tag="div" hoverable 
                content-class="video-card-content"
                :content-style="{ 'display': 'flex', 'justify-content': 'end' }"
                :style="{ 'cursor': 'pointer' }">
                    {{ video.description }}
                    包含{{ video.count }}条评论
                    <template #footer>
                        <NFlex align="center">
                                <NAvatar :src="video.author.avatar" alt="author avatar" round />
                        {{ video.author.name }}
                        </NFlex>

                    </template>
                </n-card>

            </n-gi>
        </n-grid>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { NGrid, NGi, NCard,NAvatar,NFlex } from 'naive-ui';
import { fetchVideosByTopic } from '../api';
import type { VideoDataWithCommentsCount } from '../types';
import { EllipsisText } from '../utils/name/show';
const $route = useRoute();
const topicName = ref<string>('');
const videos = ref<Array<VideoDataWithCommentsCount>>([]);
onMounted(async () => {
    topicName.value = $route.query.topicName as string;
    console.log(topicName.value);

    const ret = await fetchVideosByTopic(topicName.value);
    console.log(ret);
    
    videos.value = ret;
});


</script>