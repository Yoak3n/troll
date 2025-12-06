import { defineStore } from "pinia";

export const useRouterStore = defineStore('router', {
    state: ()=>({
        currentPath: ''
    }),
    actions: {
        setCurrentPath(path: string) {
            this.currentPath = path
        },
        getCurrentPath() {
            return this.currentPath
        }
    }
})
