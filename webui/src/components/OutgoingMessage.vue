<template>
  <div class="outgoing-message">
    <div v-if="isForwarded" class="forwarded"> Forwarded Message</div>

    <div v-if="isPhoto">
      <img :src="content" alt="Sent Photo" class="message-photo" />
    </div>
    <span v-else class="content">{{ content }}</span>

    <!-- Message Status and Timestamp -->
    <div class="message-status">
      <span class="timestamp">{{ timestamp }}</span>
      <span v-if="fullyRead" class="status">✔✔</span>
      <span v-else-if="fullyReceived" class="status">✔</span>
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
  props: ["content", "timestamp", "isPhoto", "isForwarded", "reactions", "fullyReceived", "fullyRead", "username"],
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
        reactor: this.username,  // Use the username passed via props
        content: emoji
      };
      this.reactions.push(newReaction);
      this.isReacting = false; // Close the popup after selection

      // Optionally, send this reaction to the server to save it in the database
      // Example:
      // this.$emit('send-reaction', newReaction);
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
  right: 0;
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
