<template>
    <n-card title="comment">
        <n-collapse arrow-placement="right">
            <n-collapse-item :name="group.topic" v-for="group in topics">
                <template #header>
                    <h3>{{ group.topic }}</h3>
                </template>
                <template #header-extra>
                    共{{ group.videos.length }}个评论
                </template>
                <n-card v-for="video in group.videos" :key="video.avid" embedded>
                    <!-- <div @click="jumpToReply(group.bvid,comment.id)" style="cursor: pointer;">{{ comment.content }}</div> -->
                     ""
                </n-card>
            </n-collapse-item>
        </n-collapse>
        <!-- 这不和topic差不多嘛 -->
    </n-card>
</template>

<script setup lang="ts">
import { NCollapse,NCollapseItem,NCard } from 'naive-ui';
import type { VideoDataGroupyByTopic } from '../types';
import { ref } from 'vue';
const topics = ref<VideoDataGroupyByTopic[]>([])
</script>