<template>
    <n-card style="width: 80%;" :bordered="false" size="huge" role="dialog" aria-modal="true" title="添加任务">
        <n-form :ref="formRef" :model="formData">
            <n-grid :cols="2" :x-gap="24">
                <n-form-item-gi path="type" label="任务类型">
                    <n-select v-model:value="formData.type" :options="taskTypeOption" default-value="topic" />
                </n-form-item-gi>
                <n-gi>
                    <n-grid :cols="1">
                        <n-form-item-gi v-for="(_, index) in formData.data" :path="`data[${index}]`"
                            :label="`${formData.type == 'topic' ? '关键词' : '视频号(BV)'}${formData.data.length > 1 ? index + 1 : ''}`"
                            :key="index">
                            <n-input-group>
                                <n-input v-model:value="formData.data[index]" clearable
                                    :placeholder="formData.type == 'topic' ? '请输入搜索视频的关键词' : '请输入视频的BV号'" />
                            </n-input-group>
                            <n-button-group>
                                <n-button @click="() => formData.data.splice(index, 1)">
                                    <n-icon>
                                        <Remove />
                                    </n-icon>
                                </n-button>
                                <n-button v-if="index == (formData.data.length - 1)"
                                    @click="() => formData.data.push('')" round>
                                    <n-icon>
                                        <Add />
                                    </n-icon>
                                </n-button>
                            </n-button-group>
                        </n-form-item-gi>
                    </n-grid>
                </n-gi>
                <n-form-item-gi label="话题分类" path="topic">
                    <n-input v-model:value="formData.topic" placeholder="指定任务" />
                </n-form-item-gi>
                <n-form-item-gi label="搜索结果页数" path="page" v-if="formData.type == 'topic'">
                    <n-input-number v-model:value="formData.page"></n-input-number>
                </n-form-item-gi>
            </n-grid>
        </n-form>
        <div style="display: flex;justify-content: center;"> 
            <n-button @click="submitTaskForm(formData)" style="width: 50%; margin: 0 auto" type="success">提交</n-button>
        </div>

    </n-card>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import type { VNodeRef } from 'vue';
import {
    NCard,
    NForm,
    NFormItemGi,
    NGrid, NGi,
    NSelect,
    NButtonGroup,
    NButton,
    NInputGroup, NInputGroupLabel, NInput, NInputNumber,
    NIcon
} from 'naive-ui';
import type { SelectOption } from 'naive-ui';
import { Add, Remove } from '@vicons/ionicons5'
import type { TaskForm } from '../../types';

const { submitTaskForm } = defineProps<{
    submitTaskForm: (task: TaskForm) => void
}>()
const formRef = ref<VNodeRef>()
const formData = reactive<TaskForm>({
    type: 'topic',
    data: [''],
    topic: '',
    page: 1
})

const taskTypeOption: SelectOption[] = [
    {
        "label": "搜索关键词",
        "value": "topic"
    }, {
        "label": "指定视频",
        "value": "video"
    }
]



</script>

<style scoped></style>