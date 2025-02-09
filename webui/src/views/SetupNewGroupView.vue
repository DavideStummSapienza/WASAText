<template>
  <div class="container">
    <div class="input-group">
      <label for="groupname">Groupname:</label>
      <input type="text" id="groupname" v-model="groupname" placeholder="Enter group name">
    </div>
    <div class="input-group">
      <label for="photoURL">Photo URL:</label>
      <input type="text" id="photoURL" v-model="photoURL" placeholder="Enter photo URL">
    </div>
    <button @click="createGroup">Next</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      groupname: '',
      photoURL: ''
    };
  },
  methods: {
    createGroup() {
      // Hier kannst du die Logik für das Erstellen der Gruppe hinzufügen
      console.log('Groupname:', this.groupname);
      console.log('Photo URL:', this.photoURL);
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
}

label {
  margin-right: 10px;
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
</style>