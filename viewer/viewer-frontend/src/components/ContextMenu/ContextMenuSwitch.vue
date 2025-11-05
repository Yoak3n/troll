<template>
  <user-context :context="context" v-if="type=='user'" :uid="uid!"/>
  <comment-context :context="context" v-else  :comment-id="commentId!" :bv-id="bvId!"/>
</template>

<script setup lang="ts">
import UserContext from '../User/UserContext.vue'
import CommentContext from '../Comment/CommentContext.vue'
import type {DropdownOptionProps} from '../../types/option'

interface ContextMenuSwitchProps {
    x: number;
    y: number;
    type: string
    showDropdown: boolean;
    OnClickoutside: () => void;
    uid?: number
    bvId?: string
    commentId?: number
}
const {type, x, y, showDropdown, OnClickoutside,uid,bvId,commentId} = defineProps<ContextMenuSwitchProps>()
const context: DropdownOptionProps = {
    x,
    y,
    showDropdown,
    OnClickoutside,
    HandleSelect: ()=>{}
}
</script>
