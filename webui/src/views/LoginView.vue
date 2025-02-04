<template>
  <div class="login-container">
    <h1>Login</h1>
    <input v-model="username" placeholder="Enter your username" />
    <button @click="login">Login</button>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<script>
import axios from "@/services/axios"; // Importiere die Axios-Instanz

export default {
  data() {
    return {
      username: "",
      error: null,
    };
  },
  methods: {
    async login() {
      if (!this.username) {
        this.error = "Username is required";
        return;
      }
      try {
        await axios.post("/session", { username: this.username });
        this.$router.push("/chats"); // Nach erfolgreichem Login zur Chatliste weiterleiten
      } catch (err) {
        this.error = "Login failed. Please try again.";
      }
    },
  },
};
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 100px auto;
  text-align: center;
}

input {
  width: 100%;
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 5px;
}

button {
  width: 100%;
  padding: 10px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}

.error {
  color: red;
  margin-top: 10px;
}
</style>
