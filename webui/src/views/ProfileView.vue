<template>
  <div class="profile-settings">
    <h1 class="profile-title">Profile</h1>
    <div class="profile-card">
      <!-- Profilbild -->
      <img class="profile-avatar" :src="profilePhotoURL" alt="User Avatar" />

      <!-- Datei-Upload für das Profilbild -->
      <input type="file" @change="handleFileUpload" accept="image/*" />
      <button class="upload-button" @click="uploadProfilePicture">Upload Picture</button>

      <input type="text" class="profile-name" v-model="username" placeholder="Enter username" />
      <ErrorMsg v-if="errorMsg" :message="errorMsg" />
      <div class="button-group">
        <button class="save-button" @click="saveProfile">Save</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "@/services/axios";
import ErrorMsg from '@/components/ErrorMsg.vue';

export default {
  components: { ErrorMsg },
  data() {
    return {
      username: "",
      profilePhotoURL: 'https://ui-avatars.com/api/?name=empty&size=100',
      errorMsg: "",
      imageFile: null
    };
  },
  methods: {
    async fetchUser() {
      try {
        const storedUsername = sessionStorage.getItem("currentUser");
        
        if (!storedUsername) {
          console.error("No user found in sessionStorage.");
          return;
        }

        this.username = storedUsername;

        const response = await axios.get(`/users`, {
          params: { username: this.username },
        });


        if (response.data && response.data[0].username) {
          this.username = response.data[0].username;
          sessionStorage.setItem("currentUser", this.username);
          this.profilePhotoURL = response.data[0].profile_photo_url || this.profilePhotoURL;
        }
      } catch (error) {
        console.error("Error fetching user:", error);
      }
    },

    async saveProfile() {
      try {
        this.errorMsg = "";
        const response = await axios.put("/user-profile", { newusername: this.username });
        
        if (response.data.message === 'username successfully changed') {
          alert('Username updated successfully!');
        }
      } catch (error) {
        console.error('Error saving profile:', error);
        this.errorMsg = error.response?.data?.error || "An unexpected error occurred.";
      }
    },

    handleFileUpload(event) {
      this.imageFile = event.target.files[0];
    },

    async uploadProfilePicture() {
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

        // API-Call zum Aktualisieren des Profilbildes
        await this.changeProfilePicture(imageUrl);

      } catch (error) {
        console.error("Error uploading profile picture:", error);
        alert("Failed to upload profile picture.");
      }
    },

    async changeProfilePicture(newPhotoURL) {
      try {
        const response = await axios.put('/profile-picture', { photo_url: newPhotoURL });
        if (response.data.message === 'profile picture successfully updated') {
          this.profilePhotoURL = newPhotoURL;
          alert('Profile picture updated!');
        }
      } catch (error) {
        console.error('Error changing profile picture:', error);
        alert('Failed to update profile picture.');
      }
    },

  },

  mounted() {

    this.fetchUser();

  },
};
</script>

<<style scoped>
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
  width: 120px;
  height: 120px;
  border-radius: 50%;
  margin-bottom: 20px;
}

.profile-name {
  font-size: 24px;
  padding: 12px;
  border: 1px solid black;
  border-radius: 10px;
  text-align: center;
  width: 250px;
  margin-bottom: 20px;
}

.upload-button, .save-button {
  margin-top: 10px;
  padding: 12px 18px;
  border: none;
  border-radius: 5px;
  background-color: #007bff;
  color: white;
  font-size: 16px;
  cursor: pointer;
}

.upload-button:hover, .save-button:hover {
  background-color: #0056b3;
}

.button-group {
  display: flex;
  gap: 10px;
}
</style>