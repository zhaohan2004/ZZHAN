import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import './styles/main.css'
import App from './App.vue'
import { useThemeStore } from './stores/theme'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// 应用持久化主题（默认深色）
useThemeStore().init()

app.mount('#app')
