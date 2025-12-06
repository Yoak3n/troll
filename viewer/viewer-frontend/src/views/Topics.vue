<template>
    <div class="topics-wrapper" v-if="$route.name=='topics'">
        <h1>topics</h1>
        <n-grid x-gap="12" :cols="3">
            <n-gi v-for="topic in topics" :key="topic.name" v-if="topics.length > 0">
                <n-card 
                :title="topic.name" 
                tag="button" hoverable 
                :style="{ 'cursor': 'pointer' }" 
                @click="() => $router.push({ name: 'topic', query: { topicName: topic.name } })"
                @contextmenu.prevent="(e:MouseEvent) => popTopicContext({ x: e.clientX, y: e.clientY, topicName: topic.name })"
                >
                    包含{{ topic.count }}个视频
                </n-card>
            </n-gi>>
        </n-grid>
    </div>
    <router-view></router-view>
    <TopicContext  :context="{ x: xRef, y: yRef, showDropdown: showTopicContext, OnClickoutside: ()=>showTopicContext = false, HandleSelect: () => { }}" :topicName="topicNameRef"/>

</template>


<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router'
import { NGrid, NGi, NCard } from 'naive-ui';
import mitt from 'mitt';


import { fetchTopics } from '../api';
import type { TopicsList } from '../types';
import TopicContext from '../components/Topic/TopicContext.vue'


const $router = useRouter()
const topics = ref<TopicsList[]>([]);
const fetchTopicsList = async () => {
    try {
        const response = await fetchTopics();
        topics.value = response;
    } catch (error) {
        console.error('Error fetching topics:', error);
    }
};

// popup topic context menu
const xRef = ref(-1)
const yRef = ref(-1)
const topicNameRef = ref('')
const showTopicContext = ref(false)

const popTopicContext = ({ x, y, topicName }: { x: number, y: number, topicName: string }) => {
    xRef.value = x
    yRef.value = y
    topicNameRef.value = topicName
    showTopicContext.value = true
}
const $mitt = mitt()
onMounted(() => {
    $mitt.on('topicUpdated', fetchTopicsList)
    fetchTopicsList();
});
</script>

<style scoped lang="less"></style>