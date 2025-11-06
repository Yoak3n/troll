<template>
    <div class="log-container">
          <n-virtual-list :item-size="42" :items="logsList" item-resizable ref="virtualListInst">
            <template #default="{ item }">
                <div class="log-item">
                    [{{ item.time }}] {{ item.content }}
                </div>
            </template>
        </n-virtual-list>
        </div>
</template>

<script setup lang="ts">
import { NVirtualList } from 'naive-ui';
import type { VirtualListInst } from 'naive-ui'
import type { LogItem } from '../../types';
import { nextTick, ref, watch } from 'vue';

const {logsList} = defineProps<{
    logsList: LogItem[]
}>()
const virtualListInst = ref<VirtualListInst>()
const handleScrollToPosition = ()=> {
  virtualListInst.value?.scrollTo({ position: 'bottom' })
}
watch(()=>logsList,()=> nextTick(()=>handleScrollToPosition()))
</script>

<style scoped lang="less">
.log-container{
    background-color: beige;
    height: 100%;
}
</style>