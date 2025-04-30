import { createRouter, createWebHistory } from 'vue-router'; 
import Login from '@/pages/LoginPage.vue';
import Chats from '@/components/shared/ChatList.vue';
import SendMessage from '@/components/shared/SendMessage.vue';
import Profile from '@/components/shared/ViewProfile.vue';
import SearchUsers from '@/components/shared/SearchUsers.vue';
import GroupView from '@/components/shared/GroupView.vue';

const routes = [
  { path: '/', component: Login },
  { path: '/chats', component: Chats },
  { path: '/chats/:chatId', component: SendMessage },
  { path: '/groupChats/:chatId', component: GroupView },
  { path: '/profile', component: Profile },  
  { path: '/search-users', component: SearchUsers },  
  { path: '/add-user', component: GroupView }
];

// Создаем новый роутер с использованием createWebHistory
const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
