<template>
  <div class="incoming-message">
    <span class="username">{{ username }}</span>

    <div v-if="isForwarded" class="forwarded"> Forwarded Message</div>

    <div v-if="isPhoto">
      <img :src="photoUrl" alt="Received Photo" class="message-photo" />
    </div>
    <span v-else class="content">{{ content }}</span>

    <div class="message-info">
      <span class="timestamp">{{ timestamp }}</span>
    </div>

    <button @click="react" class="reaction-button">+</button>

    <div class="reactions">
      <span v-for="(reaction, index) in reactions" :key="index">{{ reaction }}</span>
    </div>
  </div>
</template>

<script>
export default {
  props: ["username", "content", "timestamp", "isPhoto", "photoUrl", "isForwarded", "reactions"],
  methods: {
    react() {
      alert("Reacting to message!");
    }
  }
};
</script>

<style scoped>
.incoming-message {
  background: #E6DFFF;
  padding: 10px;
  border-radius: 10px;
  text-align: left;
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

.message-info {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
  color: #555;
}

.username {
  font-weight: bold;
  color: #4B0082;
  display: block;
}

.reactions {
  margin-top: 5px;
  font-size: 14px;
}
</style>
