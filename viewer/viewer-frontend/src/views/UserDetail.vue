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
                    <div @click="(e)=>onClickedComment({x:e.clientX,y:e.clientY,commentId:comment.id,bvid:group.bvid})" style="cursor: pointer;">{{ comment.content }}</div>
                </n-card>
            </n-collapse-item>
        </n-collapse>
    </n-card>
        <comment-context
        :context="{ x: xRef, y: yRef, showDropdown: showCommentContext, OnClickoutside:()=> showCommentContext=false, HandleSelect: () => { } }"
        :comment-id="clickedCommentId" :bv-id="currentVideoBvid" />
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute} from 'vue-router';
import { NCollapse, NCollapseItem, NCard } from 'naive-ui';
import CommentContext from '../components/Comment/CommentContext.vue'
import { fetchCommentsBySearchFilter } from '../api';
import type { CommentViewWithVideo, SearchFilterRequest } from '../types';

const $route = useRoute()
const comments = ref<CommentViewWithVideo[]>([])

const xRef = ref(-1)
const yRef = ref(-1)
const clickedCommentId = ref(-1)
const showCommentContext = ref(false)
const currentVideoBvid = ref('')
const onClickedComment = ({x,y,commentId,bvid}:{x:number,y:number,commentId: number, bvid: string})=>{
    xRef.value = x
    yRef.value = y
    clickedCommentId.value = commentId
    currentVideoBvid.value = bvid
    showCommentContext.value = true

}
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
    const res = await fetchCommentsBySearchFilter(data)
    comments.value = res
})


</script>

<style scoped></style>