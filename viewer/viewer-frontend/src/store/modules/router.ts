import { defineStore } from "pinia";

export const useRouterStore = defineStore('router', {
    state: ()=>({
        currentPath: '',
        currentTopicName: window.localStorage.getItem('currentTopicName') ||'',
    }),
    actions: {
        setCurrentPath(path: string) {
            this.currentPath = path
        },
        getCurrentPath() {
            return this.currentPath
        },

        setCurrentTopicName(topic: string) {
            this.currentTopicName = topic
            window.localStorage.setItem('currentTopicName',topic)
        },
        getCurrentTopicName():string{
            return this.currentTopicName
        }
    }
})
