import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import './styles/global.scss'

const app = createApp(App)

// main.ts 是前端应用启动入口。
// 这里统一挂载 Pinia、Router、Element Plus 和全局样式，所有页面都会共享这些能力。

// 注册 Element Plus 图标，业务页面可以直接使用 <el-icon> 中的图标组件。
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.mount('#app')
