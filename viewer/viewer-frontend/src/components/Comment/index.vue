<template>
    <n-list-item>
        <n-thing content-style="display:flex;width:100%">
            <div class="user-info">
                <NAvatar lazy class="avatar" :src="comment.owner.avatar" alt="user avatar" size="large" />
                <div class="username" style="font-weight: bold;cursor: pointer;word-wrap: break-word;margin-top: .2rem;"
                    @click="onClickUser">{{ comment.owner.name }}</div>
            </div>
            <div class="comment-content" style="padding: 0.2rem 1rem 0.2rem 1rem;font-weight:500;width: 95%;">
                <n-grid :cols="1" :y-gap="12">
                    <n-gi>
                        <div style="white-space: pre-wrap;cursor: pointer;" @click="onClickComment">
                            {{ comment.content}}
                        </div>
                    </n-gi>
                    <n-gi>
                        <div style="display: flex; justify-items: end;">
                            <SubCommentList :sub-comments="comment.children"  />
                        </div>
                    </n-gi>
                </n-grid>
            </div>
        </n-thing>

    </n-list-item>
</template>

<script setup lang="ts">
import { inject } from 'vue';
import SubCommentList from './SubCommentList.vue';
import { NListItem, NAvatar, NThing, NGrid, NGi } from 'naive-ui';
import type { CommentView } from '../../types';
const { comment } = defineProps<{
    comment: CommentView;
}>();

const popCommentContext: ({ x, y, commentId }: { x: number, y: number, commentId: number }) => void = inject('popCommentContext', () => { })
const popUserContext: ({ x, y, uid }: { x: number, y: number, uid: number }) => void = inject('popUserContext', () => { })

const onClickComment = (e: MouseEvent) => popCommentContext({ x: e.clientX, y: e.clientY, commentId: comment.id })
const onClickUser = (e: MouseEvent) => popUserContext({ x: e.clientX, y: e.clientY, uid: comment.owner.uid })
</script>

<style lang="less" scoped>
.user-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-bottom: 8px;
    width: 5%;
}
</style>