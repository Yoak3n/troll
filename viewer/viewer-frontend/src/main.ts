import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'

const app = createApp(App)
app.use(router)
import pinia from './store'
app.use(pinia)

app.mount('#app')


