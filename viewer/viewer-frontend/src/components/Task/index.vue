<template>
    <div class="task-process-wrapper">
        <n-progress 
        :percentage="task.completed ? 100 : (task.current / task.total) * 100"  
        :processing="!task.completed" 
        :status="task.completed ? 'success' : 'default'"
        type="line">
            <n-ellipsis style="max-width: 240px">{{ task.label }}</n-ellipsis>
        </n-progress>
    </div>
</template>

<script setup lang="ts">
import { NProgress,NEllipsis } from 'naive-ui';
import type { TaskProcess } from '../../types'
const { task } = defineProps<{
    task: TaskProcess
}>()
</script>

<style scoped lang="less">
.task-process-wrapper {
    display: flex;
}
</style>