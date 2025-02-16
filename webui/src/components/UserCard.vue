<template>
  <div 
    class="user-card" 
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
  },
  methods: {
    passChoosenUser() {
        this.$router.push({ path: '/chat', query: { username: this.user.username } });
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
