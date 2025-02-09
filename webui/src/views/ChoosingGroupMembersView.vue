<template>
  <div class="group-container">
    <h1 class="title">Members</h1>
    <div class="members-list">
      <UserCard v-for="user in selectedUsers" :key="user.id" :user="user" />
    </div>
    <div class="button-group">
      <button class="btn add" @click="navigateToSearch">Add</button>
      <button class="btn create" @click="createGroup">Create</button>
    </div>
  </div>
</template>

<script>
import UserCard from "@/components/UserCard.vue";
import { useRouter } from "vue-router";

export default {
  components: { UserCard },
  data() {
    return {
      selectedUsers: [], // Liste der hinzugefügten Benutzer
    };
  },
  methods: {
    navigateToSearch() {
      this.$router.push({ name: "search", query: { from: "group" } });
    },
    createGroup() {
      console.log("Creating group with members:", this.selectedUsers);
      // Hier würdest du eine API-Anfrage zum Erstellen der Gruppe senden
    },
  },
  mounted() {
    // Prüft, ob Benutzer in localStorage gespeichert wurden (z.B. nach dem Hinzufügen)
    const storedUsers = localStorage.getItem("selectedUsers");
    if (storedUsers) {
      this.selectedUsers = JSON.parse(storedUsers);
    }
  },
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