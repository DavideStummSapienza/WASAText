import { createRouter, createWebHashHistory } from "vue-router";
import HomeView from "@/views/HomeView.vue";
import ProfileView from "@/views/ProfileView.vue";
import GroupSettingsView from "@/views/GroupSettingsView.vue";
import LoginView from "@/views/LoginView.vue"; 
import ChatListView from "@/views/ChatListView.vue"; 
import ChatView from "@/views/ChatView.vue"; 
import SearchView from "@/views/SearchView.vue";
import ChoosingGroupMembersView from "@/views/ChoosingGroupMembersView.vue";
import SetupNewGroupView from "@/views/SetupNewGroupView.vue";
import SignUpView from "@/views/SignUpView.vue";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: "/", component: HomeView }, // Start-Site, decide Login or Sign Up
    { path: "/login", component: LoginView }, // Login-Site
    { path: "/chats", component: ChatListView }, // Chatlist after Login
    { path: "/profile", component: ProfileView }, // Profile Settings-Site
    { path: "/group-settings", component: GroupSettingsView }, // Group Settings-Site
    { path: "/chat", component: ChatView }, // specific Chat
    { path: "/search", component: SearchView }, // Search Users-Site
    { path: "/choose-members", component: ChoosingGroupMembersView }, // Choosing Group Members-Site
    { path: "/create-group", component: SetupNewGroupView }, // Group Setup-Site
    { path: "/sign-up", component: SignUpView }, // Sign Up-Site
  ],
});

export default router;