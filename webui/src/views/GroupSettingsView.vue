<template>
  <div class="group-settings">
    <h1 class="group-settings-title">Group Settings</h1>
    <div class="group-settings-card">
      <img class="group-photo" :src="groupPhotoURL" alt="Group Avatar" />
      <input type="file" @change="handleFileUpload" accept="image/*" />
      <button class="upload-button" @click="uploadGroupPhoto">Upload Group Photo</button>

      <input type="text" class="group-name" v-model="groupname" />
      <div class="button-group">
        <button class="save-button" @click="saveGroupName">Save</button>
        <button class="add-user-button" @click="addUserToGroup">Add User</button>
        <button class="leave-button" @click="leaveGroup">Leave Group</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "@/services/axios";

export default {
  data() {
    return {
      groupname: this.$route.query.groupname || "Unnamed Group",
      groupPhotoURL: "https://ui-avatars.com/api/?name=Group&size=100",
      imageFile: null,
    };
  },
  async mounted() {
    await this.fetchGroupData();
  },
  methods: {
    async fetchGroupData() {
      try {
        const response = await axios.get(`/groups/${this.groupname}`);
        if (response.data) {
          this.groupPhotoURL = response.data.group_photo_url || this.groupPhotoURL;
        }
      } catch (error) {
        console.error("Error fetching group data:", error);
      }
    },
    handleFileUpload(event) {
      this.imageFile = event.target.files[0];
    },
    async uploadGroupPhoto() {
      if (!this.imageFile) {
        alert("Please select an image first.");
        return;
      }
      try {
        const formData = new FormData();
        formData.append("image", this.imageFile);

        const uploadResponse = await axios.post("/upload", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        const imageUrl = uploadResponse.data.imageUrl;

        await axios.put(`/groups/${this.groupname}/photo`, { photo_url: imageUrl });

        this.groupPhotoURL = imageUrl;
        alert("Group photo updated!");
      } catch (error) {
        console.error("Error uploading group photo:", error);
        alert("Failed to upload group photo.");
      }
    },
    async saveGroupName() {
      try {
        await axios.put(`/groups/${this.groupname}/rename`, { newname: this.groupname });
        alert("Group name updated!");
      } catch (error) {
        console.error("Error updating group name:", error);
        alert("Failed to update group name.");
      }
    },
    async addUserToGroup() {
      const newUser = prompt("Enter username to add:");
      if (!newUser) return;
      try {
        await axios.post(`/groups/${this.groupname}/add-user`, { username: newUser });
        alert(`${newUser} added to the group!`);
      } catch (error) {
        console.error("Error adding user to group:", error);
        alert("Failed to add user.");
      }
    },
    async leaveGroup() {
      if (!confirm("Are you sure you want to leave this group?")) return;
      try {
        await axios.post(`/groups/${this.groupname}/leave`);
        alert("You have left the group.");
        this.$router.push("/home");
      } catch (error) {
        console.error("Error leaving group:", error);
        alert("Failed to leave group.");
      }
    },
  },
};
</script>

<style scoped>
.group-settings {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: #65558f;
  font-family: "Roboto", sans-serif;
  color: white;
}

.group-settings-title {
  font-size: 40px;
  margin-bottom: 20px;
}

.group-settings-card {
  background: white;
  border-radius: 10px;
  padding: 80px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.2);
}

.group-photo {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  margin-bottom: 20px;
}

.group-name {
  font-size: 20px;
  padding: 10px;
  border: 1px solid black;
  border-radius: 10px;
  text-align: center;
  width: 250px;
  margin-bottom: 20px;
}

.button-group {
  display: flex;
  gap: 10px;
}

button {
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  font-size: 16px;
  cursor: pointer;
}

.save-button {
  background-color: #21005d;
  color: white;
}

.add-user-button {
  background-color: #007bff;
  color: white;
}

.leave-button {
  background-color: red;
  color: white;
}
</style>
