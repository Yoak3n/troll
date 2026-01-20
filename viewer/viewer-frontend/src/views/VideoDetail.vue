<template>
    <div class="video-wrapper">
        <div>
            {{ }}
        </div>
        <n-spin :show="isLoading" description="" size="large" stroke="#333">
            <template #description>
                <span style="color: #444;">
                    正在获取评论数据...
                </span>
            </template>
            <div class="comments">
                <n-list class="comment-list">
                    <Comment v-for="c in commentsList" :comment="c" />
                    <n-back-top :right="40" :bottom="160" />
                </n-list>
                <n-pagination v-if="length > 100" v-model:page="page" :page-size="100" :page-count="pageCout"
                    show-size-picker @update:page="(p) => page = p" />
            </div>
        </n-spin>
        <user-context :uid="clickedUid"
            :context="{ x: xRef, y: yRef, showDropdown: showUserContext, OnClickoutside: () => showUserContext = false, HandleSelect: () => { } }" />
        <comment-context
            :context="{ x: xRef, y: yRef, showDropdown: showCommentContext, OnClickoutside: () => showCommentContext = false, HandleSelect: () => { } }"
            :comment-id="clickedCommentId" :bv-id="currentVideoInfo.bvid" />
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, provide, watch, nextTick, reactive } from 'vue';
import { useRoute } from 'vue-router';
import { NList, NPagination, NBackTop,NSpin } from 'naive-ui';
import { fetchCommentsByVideo } from '../api';
import type { CommentView, CommentViewWithVideo, VideoData } from '../types';
import UserContext from '../components/User/UserContext.vue';
import CommentContext from '../components/Comment/CommentContext.vue'
import Comment from '../components/Comment/index.vue';
import { jumpToReply } from '../utils/window/reply';


const isLoading = ref(false)
// 需要分页，不然一次性渲染上万条评论
const commentsList = ref<CommentView[]>([]);
const xRef = ref<number>(-1)
const yRef = ref<number>(-1)
const showUserContext = ref(false)
const showCommentContext = ref(false)
const contextType = ref('')

let currentVideoInfo = reactive<VideoData>({avid: 0, bvid: '', title: '', topic: '',description: '',author : {uid:0 ,name: '',avatar: '', location: ''}})
const clickedUid = ref<number>(-1)
const clickedCommentId = ref<number>(-1)
let commentsData: CommentViewWithVideo;
const length = ref(0)
const page = ref(0)
const pageCout = ref(0)
const popUserContext = ({ x, y, uid }: { x: number, y: number, uid: number }) => {
    contextType.value = 'user'
    xRef.value = x
    yRef.value = y
    clickedUid.value = uid
    showUserContext.value = true
}
const popCommentContext = ({ x, y, commentId }: { x: number, y: number, commentId: number }) => {
    contextType.value = 'comment'
    xRef.value = x
    yRef.value = y
    clickedCommentId.value = commentId
    showCommentContext.value = true
}

const jumpToCommentLocation = (commentId: number) => jumpToReply(commentsData.bvid, commentId)
provide('popUserContext', popUserContext)
provide('popCommentContext', popCommentContext)
provide('jumpToCommentLocation', jumpToCommentLocation)
const renderCommentsList = async (data: CommentView[]) => {
    if (length.value > 100) {
        const start = (page.value - 1) * 100
        let end = start + 100
        if (length.value < end) {
            end = length.value
        }
        commentsList.value = data.slice(start, end);
    } else {
        commentsList.value = data;
    }
    await nextTick(() => window.scrollTo({
        top: 0,
        behavior: 'smooth'
    })
    )
}
watch((page), (n, o) => {
    if (n == o) return
    renderCommentsList(commentsData.comments)
})
onMounted(async () => {
    const $route = useRoute();
    const avid = Number($route.query.avid);
    isLoading.value = true
    commentsData = await fetchCommentsByVideo(avid);
    isLoading.value = false
    currentVideoInfo.bvid = commentsData.bvid
    Object.assign(currentVideoInfo,commentsData)
    const pageQuery = Number($route.query.page)
    pageQuery ? page.value = pageQuery : page.value = 1
    length.value = commentsData.comments.length
    pageCout.value = Math.ceil(length.value / 100)
    renderCommentsList(commentsData.comments)
});

</script>

<style scoped></style>
