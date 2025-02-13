<template>
  <div class="chat-card" @click="goToConversation">
    <img 
      :src="chat.photo_url || `https://ui-avatars.com/api/?name=${encodeURIComponent(chat.name)}&size=40`" 
      alt="Profile Picture" 
      class="profile-picture" 
    />
    <div class="chat-info">
      <strong class="convo-name">{{ chat.name }}</strong>
      <p class="last-message">{{ chat.last_message }}</p>
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
    },
    goToConversation() {
      this.$router.push({ path: '/chat', query: { username: this.chat.name } });
    }
  }
};
</script>

<style scoped>
.chat-card {
  display: flex;
  align-items: center; /* Vertikal zentrieren */
  background: #E8DEF8;
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
  display: flex;
  flex-direction: row;
  justify-content: flex-start; /* Platzierung der Namen und Nachricht in einer Reihe */
  align-items: center; /* Vertikale Ausrichtung der Elemente */
  margin-right: 10px;
  flex-grow: 1; /* Stellt sicher, dass der Textbereich flexiblen Platz einnimmt */
}

.last-message {
  font-size: 18px;
  color: #1D192B;
  align-self: center;
  overflow: hidden; /* Verhindert das Überlaufen der Nachricht */
  text-overflow: ellipsis; /* Wenn die Nachricht zu lang ist, wird sie abgeschnitten */
  white-space: nowrap; /* Verhindert das Umbruch der Nachricht */
  max-width: calc(100% - 60px); /* Verhindert, dass die Nachricht zu groß wird */
}

.convo-name {
  font-size: 20px;
  font-weight: bold;
  color: #1D192B;
  margin-right: 400px; /* Fügt einen kleinen Abstand zwischen dem Namen und der Nachricht hinzu */
}

.chat-time {
  font-size: 12px;
  color: #1D192B;
  margin-left: auto; /* Positioniert den Zeitstempel rechts */
}
</style>
