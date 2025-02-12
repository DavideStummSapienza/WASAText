<template>
  <div class="group-settings">
    <h1 class="group-settings-title">Group Settings</h1>
    <div class="group-settings-card">
      <img class="group-photo" alt="User Avatar" src="https://ui-avatars.com/api/?name=empty&size=40" />
      <input type="text" class="group-name" v-model="groupname" />
      <div class="button-group">
        <button class="save-button">Save</button>
        <button class="leave-button">Leave Group</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      groupname: 'groupname'
    };
  }
};
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;500&display=swap');

.group-settings {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: #65558f;
  font-family: 'Roboto', sans-serif;
  color: white;
}

.group-settings-title {
  font-size: 40px;
  margin-bottom: 20px;
}

.group-settings-card {
  background: white;
  border-radius: 10px;
  padding: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.2);
}

.group-photo {
  width: 100px;
  height: 100px;
  margin-bottom: 20px;
}

.group-name {
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

.leave-button {
  background-color: red;
  color: white;
}
</style>
