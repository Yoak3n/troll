<template>
    <n-steps :current="activeStep" :status="currentStatus">
        <n-step v-for="step in steps" :key="step.title" :title="step.title" :description="step.content"/>
    </n-steps>
    <n-space justify="center">
        <n-button @click="()=>nextStep()" :type="currentStatus == 'finish' ? 'primary' : 'default'">
            {{ currentStatus == 'finish' ? '完成' : '下一步' }}
        </n-button>
    </n-space>

</template>


<script setup lang="ts">
import { computed } from 'vue'
import { NSteps, NStep, NButton,NSpace } from 'naive-ui';
import { steps } from './step'
const {activeStep} = defineProps({
    activeStep: {
        type: Number,
        default: 1
    },
    nextStep: {
        type: Function,
        default: () => {}
    }
})


// const activeStep = ref(1)
const currentStatus = computed(() => {
    if (activeStep == steps.length) {
        return 'finish'
    }
    return 'process'
})



</script>
