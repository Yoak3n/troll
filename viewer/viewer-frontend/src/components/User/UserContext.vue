<template>
    <n-dropdown placement="bottom-start" trigger="manual" :x="context.x" :y="context.y" :options="options"
        :show="context.showDropdown" :render-label="renderLabel" :on-clickoutside="onClickoutside"
        @select="handleSelect" />
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { NDropdown } from 'naive-ui';
import type { DropdownOption } from 'naive-ui';
import type { DropdownOptionProps } from '../../types/option'
import { h, type VNodeChild } from 'vue';

interface Props {
    context: DropdownOptionProps,
    uid: number
}

const $router = useRouter();
const $route = useRoute()
const { context, uid } = defineProps<Props>();

const options: Array<DropdownOption> = [
    { label: '所有评论', key: 'viewProfile', type: 'normal' },
    { label: '去TA的个人空间', key: 'userSpace', type: 'normal' },
    { label: '取消', key: 'cancel', type: 'warning' },
];

const onClickoutside = () => context?.OnClickoutside();
const handleSelect = async (key: string) => {
    switch (key) {
        case 'viewProfile':
            $router.push({ name: 'user', query: { uid: uid.toString(), topicName: $route.query.topicName as string } });
            break
        case 'userSpace':
            window.open(`https://space.bilibili.com/${uid}`, '_blank');
            break
        case 'cancel':
            break
    }
    context.OnClickoutside()
};
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


</script>