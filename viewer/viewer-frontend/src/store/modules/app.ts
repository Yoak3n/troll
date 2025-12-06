import { defineStore } from "pinia";
import { steps } from "../../components/Guidance/step";


export const useAppStore = defineStore('app', {
    state: ()=>({
        guidanceFinished: false,
        guidanceIndex: 0
    }),
    actions: {
        initGuidanceFinished() {
            this.guidanceFinished = localStorage.getItem('guidanceFinished') == 'true'
            this.guidanceIndex = parseInt(localStorage.getItem('guidanceIndex') || '1')
        },
        setGuidanceFinished(finished: boolean) {
            this.guidanceFinished = finished
            localStorage.setItem('guidanceFinished', finished.toString())
        },
        getGuidanceFinished() {
            return this.guidanceFinished
        },
        setGuidanceIndex(index: number) {
            this.guidanceIndex = index
            localStorage.setItem('guidanceIndex', index.toString())
        },
        incrementGuidanceIndex() {
            this.guidanceIndex++
            localStorage.setItem('guidanceIndex', this.guidanceIndex.toString())
            if (this.guidanceIndex == steps.length) {
                this.setGuidanceFinished(true)
            }
        },
        getGuidanceIndex() {
            return this.guidanceIndex
        }
    }
})
