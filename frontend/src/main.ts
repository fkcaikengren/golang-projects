import { createApp } from 'vue'

import App from './App.vue'
import { installAppProviders } from './app/providers'
import './app/styles.css'

const app = createApp(App)

installAppProviders(app)

app.mount('#app')
