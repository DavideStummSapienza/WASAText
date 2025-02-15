<template>
  <div class="container">
    <div class="input-group">
      <label for="groupname">Groupname:</label>
      <input type="text" id="groupname" v-model="groupname" placeholder="Enter group name">
      <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
    </div>
    <button @click="createGroup">Next</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      groupname: '',
      errorMessage: ''
    };
  },
  methods: {
    createGroup() {
      if (!this.groupname.trim()) {
        this.errorMessage = "Group name cannot be empty!";
        return;
      }
      this.errorMessage = "";
      console.log('Groupname:', this.groupname);
      this.$router.push({ path: "/choose-members", query: { groupname: this.groupname.trim() } });
    }
  }
};
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 20px;
}

.input-group {
  margin-bottom: 15px;
  display: flex;
  flex-direction: column;
}

label {
  margin-bottom: 5px;
}

input {
  padding: 5px;
  width: 200px;
}

button {
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  cursor: pointer;
}

button:hover {
  background-color: #45a049;
}

.error-message {
  color: red;
  font-size: 14px;
  margin-top: 5px;
}
</style>
