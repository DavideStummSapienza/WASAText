<template>
  <div class="incoming-message">

    <!-- Sender Username -->
    <span class="username">{{ username }}</span>

    <!-- Forwarded -->
    <div v-if="isForwarded" class="forwarded">Forwarded Message</div>

    <!-- Photo display -->
    <div v-if="isPhoto">
      <img :src="content" alt="Received Photo" class="message-photo" />
    </div>
    <span v-else class="content">{{ content }}</span>

    <!-- Timestamp -->
    <div class="message-info">
      <span class="timestamp">{{ timestamp }}</span>
    </div>

    <!-- Reaction Button -->
    <button @click="toggleReactionPopup" class="reaction-button">+</button>

    <!-- Reactions Display -->
    <div class="reactions">
      <div v-for="(reaction, index) in reactions" :key="index">
        <span>{{ reaction.reactor }}: </span><span>{{ reaction.content }}</span>
      </div>
    </div>

    <!-- Emoji Selection Popup -->
    <div v-if="isReacting" class="reaction-popup">
      <button @click="addReaction(':D')">:D</button>
      <button @click="addReaction('D:')">D:</button>
      <button @click="addReaction(':|')">:|</button>
    </div>


  </div>
</template>

<script>
export default {
  props: ["username", "content", "timestamp", "isPhoto", "isForwarded", "reactions"],
  data() {
    return {
      isReacting: false,  // Flag to toggle the emoji popup visibility
    };
  },
  methods: {
    toggleReactionPopup() {
      // Toggle the visibility of the emoji popup
      this.isReacting = !this.isReacting;
    },
    addReaction(emoji) {

      // Add the selected emoji to the reactions list along with the current username
      const newReaction = {
        content: emoji
      };

      // emit newReaction
      this.$emit("reaction-added", newReaction);

      this.isReacting = false; // close Popup

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

.reaction-button {
  background-color: #005047;
  color: white;
  padding: 5px 10px;
  border-radius: 5px;
  cursor: pointer;
  border: none;
}

.reaction-popup {
  position: absolute;
  bottom: 50px;
  left: 0;
  background: white;
  border: 1px solid #ccc;
  border-radius: 5px;
  padding: 10px;
  display: flex;
  gap: 10px;
}

.reaction-popup button {
  background-color: transparent;
  border: none;
  font-size: 20px;
  cursor: pointer;
}

.reaction-popup button:hover {
  color: #005047;
}
</style>
