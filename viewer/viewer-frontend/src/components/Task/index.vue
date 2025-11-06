<template>
    <div class="task-process-wrapper">
        <n-progress :percentage="task.completed ? 100 : (task.current / task.total) * 100">
            {{ task.label }}
        </n-progress>
    </div>
</template>

<script setup lang="ts">
import { NProgress } from 'naive-ui';
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