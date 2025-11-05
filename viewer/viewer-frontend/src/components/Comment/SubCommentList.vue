<template>
    <div class="sub-comment-list-wrapper" v-if="subComments && subComments.length > 0">
        <n-button text-color="#888" text @click="()=>isExpanded=!isExpanded">{{isExpanded? "收起":`共有${ subComments.length}条子评论`}}</n-button>
        <SubComment v-show="isExpanded" v-for="subComment in subComments" :key="subComment.id" :comment="subComment"  />
    </div>
</template>

<script setup lang="ts">
import { ref} from "vue"
import { NButton } from 'naive-ui';
import SubComment from './SubComment.vue';
import type { CommentView } from '../../types';

const { subComments } = defineProps<{
    subComments?: CommentView[],
}>();
const isExpanded = ref(false);
</script>

<style scoped>
</style>