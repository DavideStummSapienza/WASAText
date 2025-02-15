<template>
  <div class="group-container">
    <h1 class="title">Members</h1>
    <div class="members-list">
      <UserCard 
        v-for="user in selectedUsers" 
        :key="user.id" 
        :user="user" 
        :forGroup="true"
        :disableClick="true"
      />
    </div>
    <div class="button-group">
      <button class="btn add" @click="navigateToSearch">Add</button>
      <button class="btn cancel" @click="cancelGroupCreation">Cancel</button>
      <button class="btn create" @click="createGroup">Create</button>
    </div>
  </div>
</template>

<script>
import UserCard from "@/components/UserCard.vue";

export default {
  components: { UserCard },
  data() {
    return {
      selectedUsers: [],
    };
  },
  methods: {
    navigateToSearch() {
      this.$router.push("/search");
    },
    createGroup() {
      console.log("Creating group with members:", this.selectedUsers);
      localStorage.removeItem("selectedUsers");
      this.$router.push("/chats");
    },
    cancelGroupCreation() {
      localStorage.removeItem("selectedUsers");
      this.$router.push("/");
    },
    loadSelectedUsers() {
      const storedUsers = localStorage.getItem("selectedUsers");
      this.selectedUsers = storedUsers ? JSON.parse(storedUsers) : [];
    }
  },
  mounted() {
    this.loadSelectedUsers();
  }
};
</script>


<style scoped>
.group-container {
  background: white;
  padding: 20px;
  border-radius: 10px;
  text-align: center;
}

.title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 20px;
}

.members-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.button-group {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}

.btn {
  padding: 10px 20px;
  border-radius: 5px;
  font-weight: bold;
  cursor: pointer;
}

.add {
  background: purple;
  color: white;
}

.create {
  background: darkblue;
  color: white;
}
</style>