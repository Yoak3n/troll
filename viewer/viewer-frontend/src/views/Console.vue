<template>
    <div class="console-wrapper">
        <div class="console-header">
            状态 <n-badge :type="websocketStatus ? 'success' : 'error'" dot />
        </div>
        <div class="console-body">
            <div class="tasks-woker">
                <n-card title="任务池">
                    <template #header-extra>
                        <n-button @click="showTaskModal = true">添加任务</n-button>
                    </template>

                    <n-list v-if="taskPool.size != 0">
                        <n-list-item v-for="task in taskPool.values()" :key="task.id">
                            <TaskProcessBar :task="task" />
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


    </div>
    <n-modal v-model:show="showTaskModal">
        <TaskModal :submit-task-form="uploadTaskData" />
    </n-modal>
</template>


<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { NCollapse, NCollapseItem, NBadge, NCard, NList, NListItem, NButton, NModal } from 'naive-ui';
import type { LogItem, TaskForm, TaskProcess, WsMessage } from '../types';
import LogList from '../components/LogList/index.vue'
import TaskProcessBar from '../components/Task/index.vue'
import TaskModal from '../components/Task/TaskModal.vue';
import { random } from 'lodash-es';

const websocketStatus = ref(false)
const $websocket = ref<WebSocket | null>(null)
const showTaskModal = ref(false)
let clientId = ""
let heartBeatTimer: number
let reconnectTimer: number | undefined
let toClose = false
const initWebSocket = () => {
    //初始化weosocket
    clientId = random(1, 1000).toString()
    const wsuri = import.meta.env.VITE_APP_BASE_API + "/ws/" + clientId; //ws地址
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
        console.log(object);
        switch (object.action) {
            case "Log":
                const item: LogItem = object.data
                if (logsList.value.length > 500) {
                    logsList.value.shift()
                }
                logsList.value.push(item)
                break
            case "Ping":
                pongHandle()
                break
            case "Processes":
                const process = object.data
                Object.keys(process).forEach((key) => {
                    taskPool.value.set(key, process[key])
                })
                break
            case "Close":
                closeConnection()
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
const taskPool = ref<Map<string, TaskProcess>>(new Map())

const reconnectHandle = () => {
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
    }, 5000)
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

const pongHandle = () => {
    if (websocketStatus) {
        const pong = JSON.stringify({ action: "Pong", data: "pong" })
        $websocket.value?.send(pong)
    } else {
        clearInterval(heartBeatTimer)
    }
}
const closeConnection = () => {
    const closeMsg = JSON.stringify({ action: "Close", data: clientId })
    $websocket.value?.send(closeMsg)
    toClose = true
    $websocket.value?.close(1000)
    $websocket.value = null
}
onMounted(() => {
    initWebSocket()
})
onUnmounted(() => {
    closeConnection()
    clearInterval(heartBeatTimer)
    clearInterval(reconnectTimer)
})

const uploadTaskData = (task: TaskForm) => {
    const taskMessage: WsMessage = {
        action: 'Task',
        data: task
    }
    $websocket.value?.send(JSON.stringify(taskMessage))
    showTaskModal.value = false
}

</script>

<style scoped lang="less">
.console-wrapper {
    display: flex;
    flex-direction: column;
    margin: 1rem;

    .console-header {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-end;
        padding: .5rem 2rem;
    }

    .console-body {
        display: flex;
        flex-direction: column;
    }
}
</style>
