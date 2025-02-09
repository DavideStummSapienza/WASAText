<template>
  <div class="profile-settings">
    <h1 class="profile-title">Profile</h1>
    <div class="profile-card">
      <img class="profile-avatar" alt="User Avatar" src="/profile-avatar.png" />
      <input type="text" class="profile-name" v-model="username" />
      <div class="button-group">
        <button class="save-button">Save</button>
        <button class="delete-button">Delete Account</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: 'Tom'
    };
  }
};
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;500&display=swap');

.profile-settings {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: #65558f;
  font-family: 'Roboto', sans-serif;
  color: white;
}

.profile-title {
  font-size: 40px;
  margin-bottom: 20px;
}

.profile-card {
  background: white;
  border-radius: 10px;
  padding: 80px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.2);
}

.profile-avatar {
  width: 100px;
  height: 100px;
  margin-bottom: 20px;
}

.profile-name {
  font-size: 20px;
  padding: 10px;
  border: 1px solid black;
  border-radius: 10px;
  text-align: center;
  width: 200px;
  margin-bottom: 20px;
}

.button-group {
  display: flex;
  gap: 20px;
}

button {
  padding: 10px 20px;
  border: none;
  border-radius: 26px;
  font-size: 16px;
  cursor: pointer;
  font-weight: 500;
}

.save-button {
  background-color: #21005d;
  color: white;
}

.delete-button {
  background-color: red;
  color: white;
}
</style>
