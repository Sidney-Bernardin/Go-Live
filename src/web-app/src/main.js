import { createApp } from 'vue'
import store from './store'
import router from './router'
import './style.css'
import App from './App.vue'
import axios from 'axios'

axios.defaults.baseURL = import.meta.env.VITE_MICROSERVICES_URL

createApp(App).use(router).use(store).mount('#app')
