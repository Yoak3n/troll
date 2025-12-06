<template>
    <n-dropdown placement="bottom-start" trigger="manual" :x="context.x" :y="context.y" :options="options"
        :show="context.showDropdown" :render-label="renderLabel" :on-clickoutside="onClickoutside"
        @select="handleSelect" />
</template>


<script setup lang="ts">
import { h, ref, type VNodeChild } from 'vue'
import { useRouter } from 'vue-router'
import { NDropdown } from 'naive-ui'
import type { DropdownOption } from 'naive-ui';
import type { DropdownOptionProps } from '../../types/option'
import { deleteTopic, updateTopic } from '../../api/modules/topics'
import mitt from 'mitt';
import TopicDialogue from './TopicDialogue.vue'

interface Props {
    context: DropdownOptionProps,
    topicName: string
}

const $mitt = mitt()

const $router = useRouter()
const options: Array<DropdownOption> = [
    { label: '查看话题', key: 'viewTopic', type: 'normal' },
    { label: '重命名话题', key: 'renameTopic', type: 'normal' },
    { label: '删除话题', key: 'deleteTopic', type: 'warning' },
    { label: '取消', key: 'cancel', type: 'normal' },
];
const renameTopicValue = ref('')
const updatedTopicName = (value: string| undefined) =>renameTopicValue.value = value || ''
const { context, topicName } = defineProps<Props>();
const onClickoutside = () => context?.OnClickoutside();
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
const handleSelect = async (key: string) => {
    switch (key) {
        case 'viewTopic':
            $router.push({ name: 'Topic', query: { topicName } })
            break
        case 'renameTopic':
            window.$dialog?.create({
                title: '重命名话题',
                content: ()=> h(TopicDialogue, { modelValue: renameTopicValue.value, 'onUpdate:modelValue': updatedTopicName }),
                positiveText: '确认',
                onPositiveClick: async () => {
                    if (renameTopicValue.value != '') {
                        await updateTopic(topicName, renameTopicValue.value)
                        $mitt.emit('topicUpdated')
                    }else{
                        window.$message?.error('请输入话题名称')
                    }
                }
                })
            break
        case 'deleteTopic':
            window.$dialog?.error({
                title: '删除话题',
                content: ()=> h('p', {}, `确定删除话题 ${topicName} 吗？`),
                positiveText: '确认',
                negativeText: '取消',
                onPositiveClick: async () => {
                    await deleteTopic(topicName)
                    $mitt.emit('topicDeleted')
                }
            })
            break
        case 'cancel':
            break
    }
    context.OnClickoutside()
};
</script>