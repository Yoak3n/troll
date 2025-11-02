<template>
    <h1>topics</h1>
    <n-grid x-gap="12" :cols="3">
        <n-gi v-for="topic in topics" :key="topic.name">
            <n-card 
            :title="topic.name"
            tag="button" 
            hoverable 
            @click="()=>{

            }">
                {{ topic.count }}
            </n-card>

        </n-gi>>
    </n-grid>


</template>


<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { fetchTopics } from '../api';

import { NGrid, NGi, NCard } from 'naive-ui';
import type { TopicsList } from '../types/topic';

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