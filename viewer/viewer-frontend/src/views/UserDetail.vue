<template>
    <n-card title="comment">
        <n-collapse arrow-placement="right">
            <n-collapse-item :name="group.avid" v-for="group in comments" >
                <template #header>
                    <h3>{{ group.title }}</h3>
                </template>
                <template #header-extra>
                    共{{ group.comments.length }}条评论
                </template>
                <n-card v-for="comment in group.comments" :key="comment.id" embedded>
                    <div @click="jumpToReply(group.bvid,comment.id)" style="cursor: pointer;">{{ comment.content }}</div>
                </n-card>
            </n-collapse-item>
        </n-collapse>
    </n-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { NCollapse, NCollapseItem, NCard } from 'naive-ui';
import { fetchCommentsBySearchFilter } from '../api/modules/users';
import type { CommentViewWithVideo, SearchFilterRequest } from '../types';
import { jumpToReply} from '../utils/window/reply'
const $route = useRoute()
const comments = ref<CommentViewWithVideo[]>([])
onMounted(async () => {
    const uid = $route.query.uid ? Number($route.query.uid) : undefined
    const name = $route.query.name as string
    const rangeType = $route.query.rangeType as string
    const rangeData = $route.query.rangeData as string
    const data: SearchFilterRequest = {
        name,
        uid,
        rangeType,
        rangeData
    }
    console.log(data);

    const res = await fetchCommentsBySearchFilter(data)
    console.log(res);
    comments.value = res
})


</script>

<style scoped></style>