<template>
  <div class="outgoing-message">
    <div v-if="isForwarded" class="forwarded">🔄 Forwarded Message</div>

    <div v-if="isPhoto">
      <img :src="photoUrl" alt="Sent Photo" class="message-photo" />
    </div>
    <span v-else class="content">{{ content }}</span>

    <div class="message-status">
      <span class="timestamp">{{ timestamp }}</span>
      <span v-if="fullyRead" class="status">✔✔</span>
      <span v-else-if="fullyReceived" class="status">✔</span>
    </div>

    <button @click="react" class="reaction-button">+</button>

    <div class="reactions">
      <span v-for="(reaction, index) in reactions" :key="index">{{ reaction }}</span>
    </div>
  </div>
</template>

<script>
export default {
  props: ["content", "timestamp", "isPhoto", "photoUrl", "isForwarded", "reactions", "fullyReceived", "fullyRead"],
  methods: {
    react() {
      alert("Reacting to message!");
    }
  }
};
</script>

<style scoped>
.outgoing-message {
  background: #D6C5F0;
  padding: 10px;
  border-radius: 10px;
  text-align: right;
  position: relative;
}

.message-photo {
  max-width: 100%;
  border-radius: 10px;
}

.forwarded {
  font-size: 12px;
  color: gray;
}

.message-status {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
  color: #555;
}

.status {
  font-weight: bold;
  color: green;
}

.reactions {
  margin-top: 5px;
  font-size: 14px;
}
</style>
