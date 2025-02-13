<template>
  <div class="profile-settings">
    <h1 class="profile-title">Profile</h1>
    <div class="profile-card">
      <!-- Profilbild -->
      <img class="profile-avatar" :src="profilePhotoURL" alt="User Avatar" />
      <input type="text" class="profile-name" v-model="username" placeholder="Enter username" />
      <div class="button-group">
        <!-- Speichern Button -->
        <button class="save-button" @click="saveProfile">Save</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'; // Wir verwenden axios für API-Anfragen

export default {
  data() {
    return {
      username: 'Tom', // Der aktuelle Benutzername
      profilePhotoURL: 'https://ui-avatars.com/api/?name=Tom&size=100', // Das Profilbild-URL (Initialwert)
    };
  },
  methods: {
    // Methode, um den Benutzernamen zu speichern
    async saveProfile() {
      try {
        // API-Aufruf zum Ändern des Benutzernamens
        const response = await axios.put('/user-profile', {
          newusername: this.username,
        });

        if (response.data.message === 'username successfully changed') {
          alert('Username updated successfully!');
        }
      } catch (error) {
        console.error('Error saving profile:', error);
        alert('Failed to save profile.');
      }
    },

    // Methode, um das Profilbild zu ändern
    async changeProfilePicture(newPhotoURL) {
      try {
        // API-Aufruf zum Ändern des Profilbildes
        const response = await axios.put('/profile-picture', {
          photo_url: newPhotoURL,
        }, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('authToken')}` // Authentifizierungs-Token
          }
        });

        if (response.data.message === 'profile picture successfully updated') {
          this.profilePhotoURL = newPhotoURL; // Profilbild im Frontend aktualisieren
          alert('Profile picture updated!');
        }
      } catch (error) {
        console.error('Error changing profile picture:', error);
        alert('Failed to update profile picture.');
      }
    },
  },
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
</style>
