<template>
    <n-input v-model:value="inputModelValue" placeholder="请输入新话题名称"></n-input>
</template>

<script setup lang="ts">
import { NInput } from 'naive-ui';

const inputModelValue = defineModel<string>('modelValue')


</script>