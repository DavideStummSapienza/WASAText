<template>
  <div class="search-container">
    <h1 class="title">Choose User</h1>
    <input v-model="searchQuery" type="text" placeholder="Search..." class="search-bar" @input="fetchUsers" />
    <div class="user-list">
      <UserCard v-for="user in users" :key="user.username" :user="user" />
    </div>
  </div>
</template>

<script>
import axios from "@/services/axios";
import UserCard from "@/components/UserCard.vue";

export default {
  components: {
    UserCard,
  },
  data() {
    return {
      searchQuery: "",
      users: [],
    };
  },
  methods: {
    async fetchUsers() {
      try {
        const response = await axios.get("/users", {
          params: { username: this.searchQuery },
        });
        this.users = response.data;
      } catch (error) {
        console.error("Error fetching users", error);
      }
    },
  },
  mounted() {
    this.fetchUsers();
  },
};
</script>

<style scoped>
.search-container {
  text-align: center;
  background-color: #65558f;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 20px;
}

.title {
  font-size: 32px;
  color: white;
  font-style: italic;
}

.search-bar {
  width: 80%;
  padding: 10px;
  font-size: 18px;
  border-radius: 10px;
  border: none;
  margin: 10px 0;
}

.user-list {
  width: 80%;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>