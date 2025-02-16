<template>
  <div class="container">
    <div class="input-group">
      <label for="groupname">Groupname:</label>
      <input type="text" id="groupname" v-model="groupname" placeholder="Enter group name">
      <label for="members">Add Members:</label>

      <div class="add-members">
        <input type="text" id="members" v-model="newMember" placeholder="Enter an username">
        <button @click="addMember"> Add </button>
      </div>

      <label for="memberlist">Added Members:</label>
      <p id="memberlist" v-for="(member, index) in addedMembers" :key="index">
        {{ member }}
      </p>

      <ErrorMsg v-if="errorMessage" :message="errorMessage" />
    </div>
    <button @click="createGroup">Create</button>
  </div>
</template>

<script>
import axios from "@/services/axios";
import ErrorMsg from '../components/ErrorMsg.vue';

export default {
  components: {ErrorMsg},
  data() {
    return {
      addedMembers: [],
      groupname: '',
      newMember:"",
      errorMessage: ''
    };
  },
  methods: {
    async createGroup() {
      try {

        if (!this.groupname.trim()) {
          this.errorMessage = "Group name cannot be empty!";
          return;
        }

        this.errorMessage = "";
        console.log('Groupname:', this.groupname);

        const response = await axios.post("/groups", {
          groupname: this.groupname,
          names: this.addedMembers
        })


        this.$router.push({ path: "/chats", query: { groupname: this.groupname.trim() } });

      } catch(error) {
        
        if (error.response && error.response.data.error) {
          this.errorMessage = error.response.data.error;
        } else {
          this.errorMessage = "An unexpected error occurred.";
        }

      }
    },
    addMember() {
    if (!this.newMember.trim()) {
      this.errorMessage = "Member name cannot be empty!";
      return;
    }
    this.errorMessage = "";
    this.addedMembers.push(this.newMember.trim()); // add member
    this.newMember = ""; // clear input field
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

.add-members {
  display: flex;
  flex-direction: row;
  gap: 20px
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
