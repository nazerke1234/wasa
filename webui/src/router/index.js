import { createRouter, createWebHashHistory } from "vue-router";
import ProfileView from "../views/ProfileView.vue";
import GroupView from "../views/GroupView.vue"
import CreateGroupView from "../views/CreateGroupView.vue"
import EditGroupView from "../views/EditGroupView.vue"
import LoginView from "../views/LoginView.vue";
import HomeView from "../views/HomeView.vue";
import SearchView from "../views/SearchView.vue";
import ChatView from "../views/ChatView.vue";


const routes = [
  { path: "/", component: LoginView },
  { path: "/home", component: HomeView },
  { path: "/search", component: SearchView },
  { path: "/chats/:uuid", name: "ChatView", component: ChatView, props: true },
  { path: "/me", component: ProfileView},
  { path: "/groups", component: GroupView},
  { path: "/new-group", component: CreateGroupView},
  { path: "/groups/:uuid", name: "GroupEditView", component: EditGroupView, props: true}
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});


export default router;
