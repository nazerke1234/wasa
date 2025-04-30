import { createApp } from 'vue'; 
import App from './App.vue';
import router from './router';  // Import router

// Creating app
createApp(App)
  .use(router)
  .mount('#app'); // Mounting the app to an element with the id="app"
