<template>
    <n-dropdown placement="bottom-start" trigger="manual" :x="context.x" :y="context.y" :options="options"
        :show="context.showDropdown" :on-clickoutside="onClickoutside" :renderLabel="renderLabel"
        @select="handleSelect" />
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { NDropdown } from 'naive-ui';
import type { DropdownOption } from 'naive-ui';
import type { DropdownOptionProps } from '../../types/option';
import { jumpToReply } from '../../utils/window/reply';
import { h, type VNodeChild } from 'vue';

interface Props {
    context: DropdownOptionProps
    commentId: number
    bvId: string
}

const $router = useRouter();
const { context, commentId, bvId } = defineProps<Props>();

const options: Array<DropdownOption> = [
    { label: '跳转该评论', key: 'jump', type: 'normal' },
    { label: '查找类似评论', key: 'searchSimilar', type: 'normal' },
    { label: '取消', key: 'cancel', type: 'warning' },
];
const renderLabel = (option: DropdownOption) => {
    if (option.type === 'warning') {
        return h(
            'a', { style: { color: 'red' } },
            {
                default: () => option.label as VNodeChild
            }
        )
    }
    return option.label as VNodeChild
}

const onClickoutside = () => context.OnClickoutside!();
const handleSelect = async (key: string) => {
    switch (key) {
        case 'jump':
            jumpToReply(bvId, commentId)
            break
        case 'searchSimilar':
            $router.back()
            break
        case 'cancel':
            break
    }
    context.OnClickoutside!()
};

</script>