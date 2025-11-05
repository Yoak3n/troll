<template>
    <n-dropdown placement="bottom-start" trigger="manual" :x="context.x" :y="context.y" :options="context.options"
        :show="context.showDropdown" :render-label="renderLabel" :on-clickoutside="context.OnClickoutside"
        @select="context.HandleSelect" />
</template>

<script setup lang="ts">
import { h, onMounted, type VNodeChild } from 'vue';
import { NDropdown } from 'naive-ui';
import type { DropdownOption } from 'naive-ui';
import type { DropdownOptionProps } from '../../types/option'

interface BasicDropdownOptionProps {
    context: DropdownOptionProps
}

const { context } = defineProps<BasicDropdownOptionProps>();
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
onMounted(()=>{
    console.log("popa",context.showDropdown);
})
</script>
