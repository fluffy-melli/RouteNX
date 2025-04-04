import { createApp } from 'vue'
import App from './App.vue'
import i18n from './language.js'
import './assets/style.scss'

const app = createApp(App)
app.use(i18n)
app.mount('#app')