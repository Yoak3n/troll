<template>
    <div class="console-wrapper">
        状态 <n-badge :type="websocketStatus ? 'success' : 'error'" dot />
        <div class="tasks-woker">
            <n-card title="任务池" >
                <n-list v-if="taskPool.length != 0">
                    <n-list-item v-for="task in taskPool" :key="task.id">
                        <TaskProcessBar :task="task"/>
                    </n-list-item>
                </n-list>
            </n-card>
        </div>
        <n-collapse arrow-placement="right" default-expanded-names="log" class="">
            <n-collapse-item name="log">
                <template #header>
                    <h4>运行日志</h4>
                </template>
                <div class="log-box">
                    <LogList :logs-list="logsList" />
                </div>
            </n-collapse-item>
        </n-collapse>
    </div>
</template>


<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { NCollapse, NCollapseItem, NBadge, NCard,NList,NListItem } from 'naive-ui';
import type { LogItem, TaskProcess, WsMessage } from '../types';
import LogList from '../components/LogList/index.vue'
import TaskProcessBar from '../components/Task/index.vue'
import { random } from 'lodash-es';

const websocketStatus = ref(false)
const $websocket = ref<WebSocket | null>(null)
let clientId = ""
let heartBeatTimer: number
let reconnectTimer: number | undefined
let toClose = false
const initWebSocket = () => {
    //初始化weosocket
    const wsuri = import.meta.env.VITE_APP_BASE_API + "/ws/"+ clientId; //ws地址
    $websocket.value = new WebSocket(wsuri);
    $websocket.value.onopen = () => {
        console.log('连接成功')
        websocketStatus.value = true
        heartBeatTimer = setInterval(pingHandle, 20000)
    };
    $websocket.value.onerror = (err) => {
        console.log('连接出错', err)
        websocketStatus.value = false
        clearInterval(heartBeatTimer)
    };
    $websocket.value.onmessage = (event) => {
        if (event.data == "success") { return }
        const object: WsMessage = JSON.parse(event.data)
        switch (object.action) {
            case "Log":
                const item: LogItem = object.data
                if (logsList.value.length > 500) {
                    logsList.value.shift()
                }
                logsList.value.push(item)
                break
            default:
                break
        }
    };
    $websocket.value.onclose = () => {
        console.log('连接关闭')
        websocketStatus.value = false
        clearInterval(heartBeatTimer)
        if (!reconnectTimer && !toClose) {
            reconnectTimer = reconnectHandle()
        }
    };
}
const logsList = ref<LogItem[]>([])
const taskPool = ref<TaskProcess[]>([])
const reconnectHandle = () => {
    console.log();
    let index = 0
    const timer = setInterval(() => {
        initWebSocket()
        if (websocketStatus.value || index >= 30) {
            clearInterval(timer)
            reconnectTimer = undefined
        } else {
            index++
            console.log(`重新连接，第${index}次重试...`, timer);
        }
    }, 1000)
    return timer
}

const pingHandle = () => {
    if (websocketStatus) {
        const ping = JSON.stringify({ action: "Ping", data: "ping" })
        $websocket.value?.send(ping)
    } else {
        clearInterval(heartBeatTimer)
    }
}

onMounted(() => {
    clientId = random(1,1000).toString()
    initWebSocket()
})
onUnmounted(() => {
    closeConnection()
    clearInterval(heartBeatTimer)
    clearInterval(reconnectTimer)
})
const closeConnection = ()=>{
    const closeMsg = JSON.stringify({action: "Close",data: clientId})
    $websocket.value?.send(closeMsg)
    toClose = true
    $websocket.value?.close(1000)
    $websocket.value = null
}

</script>

<style scoped lang="less">
.console-wrapper {
    margin: 1rem;
}
</style>