import { createApp } from 'vue'
import App from './App.vue'

import store from './store'
import router from './router'

import axios from 'axios'
import './style.css'

axios.defaults.baseURL = 'http://' + import.meta.env.VITE_MICROSERVICES_URL

createApp(App).use(router).use(store).mount('#app')
