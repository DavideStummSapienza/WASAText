<template>
  <div class="home">
    <h1>Welcome to the WASAText</h1>

    <!-- Register and Login Buttons -->
    <div class="auth-buttons">
      <button @click="goToLogin" class="auth-button">Login</button>
      <button @click="goToRegister" class="auth-button">Register</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HomeView',
  methods: {
    // Navigate to the Login view
    goToLogin() {
      this.$router.push({ name: 'login' });
    },
    // Navigate to the Register view
    goToRegister() {
      this.$router.push({ name: 'register' });
    }
  }
};
</script>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

.auth-buttons {
  margin-top: 20px;
}

.auth-button {
  padding: 10px 20px;
  font-size: 16px;
  margin: 10px;
  cursor: pointer;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
}

.auth-button:hover {
  background-color: #0056b3;
}
</style>
