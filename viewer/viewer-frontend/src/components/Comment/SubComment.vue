<template>
    <div class="sub-comment-wrapper">
        <div class="sub-comment-item">
            <div class="user-info" @click="handleClickUsername" style="cursor: pointer;">
                <NAvatar lazy round class="avatar" :src="comment.owner.avatar" alt="user avatar" />
                {{ comment.owner.name }}
            </div>
            <div class="comment-content" @click="handleClickComment" style="cursor: pointer;">
                {{ comment.content }}
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { NAvatar } from 'naive-ui';
import type { CommentView } from '../../types';
import { inject } from 'vue';
const { comment } = defineProps<{
    comment: CommentView,
}>();
const popCommentContext: ({ x, y, commentId }: { x: number, y: number, commentId: number }) => void = inject('popCommentContext', () => { })
const popUserContext: ({ x, y, uid }: { x: number, y: number, uid: number }) => void = inject('popUserContext', () => {})

const handleClickUsername = (e: MouseEvent) =>popUserContext({ x: e.clientX, y: e.clientY, uid: comment.owner.uid })
const handleClickComment = (e:MouseEvent) =>  popCommentContext({ x: e.clientX, y: e.clientY, commentId: comment.id })
</script>

<style scoped lang="less">
.sub-comment-item {
    margin-top: .5rem;
    display: flex;
    .user-info {
        display: flex;
        margin-right: 1rem;
        align-items: start;
        color: #666;
        white-space: nowrap;
        font-weight: 400;
        .avatar {
            width: 24px;
            height: 24px;
            margin-right: .5rem;
        }
    }
    .comment-content{
        font-weight: 500;
        font-size: 15px;
    }
}
</style>