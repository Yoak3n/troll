<template>
    <div class="topics-wrapper" v-if="$route.name=='topics'">
        <h1>topics</h1>
        <n-grid x-gap="12" :cols="3">
            <n-gi v-for="topic in topics" :key="topic.name" v-if="topics.length > 0">
                <n-card :title="topic.name" tag="button" hoverable :style="{ 'cursor': 'pointer' }" @click="() => {
                    $router.push({ name: 'topic', query: { topicName: topic.name } });
                }">
                    包含{{ topic.count }}个视频
                </n-card>
            </n-gi>>
        </n-grid>
    </div>
    <router-view></router-view>


</template>


<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { fetchTopics } from '../api';

import { NGrid, NGi, NCard } from 'naive-ui';
import type { TopicsList } from '../types';

const $router = useRouter();
const topics = ref<TopicsList[]>([]);

const fetchTopicsList = async () => {
    try {
        const response = await fetchTopics();
        topics.value = response;
    } catch (error) {
        console.error('Error fetching topics:', error);
    }
};

onMounted(() => {
    fetchTopicsList();
});
</script>

<style scoped lang="less"></style>