<template>
  <div 
    class="user-card" 
    :class="{ 'disabled': forGroup && disableClick }"
    @click="passChoosenUser"
  >
    <img 
      :src="user.profile_photo_url || `https://ui-avatars.com/api/?name=${encodeURIComponent(user.username)}&size=40`" 
      alt="Profile Picture" 
      class="profile-picture" 
    />
    <span class="username">{{ user.username }}</span>
  </div>
</template>

<script>
export default {
  props: {
    user: Object,
    forGroup: Boolean,
    disableClick: Boolean, // ⬅️ Neue Prop für Deaktivierung
  },
  methods: {
    passChoosenUser() {
      if (this.disableClick) return; // ⬅️ Blockiert Klick in ChooseMemberView

      if (this.forGroup) {
        // Bestehende Benutzer laden
        let selectedUsers = JSON.parse(localStorage.getItem("selectedUsers")) || [];

        // Prüfen, ob der Benutzer schon in der Liste ist
        const exists = selectedUsers.some(u => u.id === this.user.id);
        if (!exists) {
          selectedUsers.push(this.user);
          localStorage.setItem("selectedUsers", JSON.stringify(selectedUsers));
        }

        // Zurück zu `ChooseMemberView`
        this.$router.push("/choose-members");
      } else {
        this.$router.push({ path: '/chat', query: { username: this.user.username } });
      }
    },
  },
};
</script>

<style scoped>
.user-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #d6c3ff;
  padding: 15px;
  border-radius: 10px;
  font-size: 18px;
  font-weight: bold;
  cursor: pointer;
}

.user-card.disabled {
  pointer-events: none;
  opacity: 0.5;
}
</style>
