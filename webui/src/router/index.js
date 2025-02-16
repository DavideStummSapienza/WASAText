import { createRouter, createWebHashHistory } from "vue-router";
import ProfileView from "@/views/ProfileView.vue";
import GroupSettingsView from "@/views/GroupSettingsView.vue";
import LoginView from "@/views/LoginView.vue"; 
import ChatListView from "@/views/ChatListView.vue"; 
import ChatView from "@/views/ChatView.vue"; 
import SearchView from "@/views/SearchView.vue";
import CreateGroupView from "@/views/CreateGroupView.vue";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: "/", component: LoginView }, // Login-Site
    { path: "/chats", component: ChatListView }, // Chatlist after Login
    { path: "/profile", component: ProfileView }, // Profile Settings-Site
    { path: "/group-settings", component: GroupSettingsView }, // Group Settings-Site
    { path: "/chat", component: ChatView }, // specific Chat
    { path: "/search", component: SearchView }, // Search Users-Site
    { path: "/create-group", component: CreateGroupView }, // Create New Group-Site
  ],
});

export default router;