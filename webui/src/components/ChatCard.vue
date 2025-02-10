<template>
  <div class="chat-card">
    <img 
      :src="chat.photo_url || `https://ui-avatars.com/api/?name=${encodeURIComponent(chat.name)}&size=40`" 
      alt="Profile Picture" 
      class="profile-picture" 
    />
    <div class="chat-info">
      <strong>{{ chat.name }}</strong>
      <p>{{ chat.last_message }}</p>
    </div>
    <span class="chat-time">{{ formatTime(chat.last_message_time) }}</span>
  </div>
</template>

<script>
export default {
  props: {
    chat: Object
  },
  methods: {
    formatTime(timestamp) {
      if (!timestamp) return "";
      return new Date(timestamp).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    }
  }
};
</script>

<style scoped>
.chat-card {
  display: flex;
  align-items: center;
  background: white;
  border-radius: 10px;
  padding: 10px;
  margin-bottom: 10px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
}

.profile-picture {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover; /* Stellt sicher, dass das Bild gut aussieht */
  margin-right: 10px;
}

.chat-info {
  flex-grow: 1;
}

.chat-time {
  font-size: 12px;
  color: gray;
}
</style>
