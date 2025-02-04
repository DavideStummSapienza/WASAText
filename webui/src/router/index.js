import { createRouter, createWebHashHistory } from "vue-router";
import HomeView from "@/views/HomeView.vue";
import LoginView from "@/views/LoginView.vue"; // Login-Seite
import ChatListView from "@/views/ChatListView.vue"; // Chatliste (später erstellen)

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: "/", component: HomeView },
    { path: "/login", component: LoginView }, // Login-Seite
    { path: "/chats", component: ChatListView }, // Chatliste nach Login
  ],
});

export default router;