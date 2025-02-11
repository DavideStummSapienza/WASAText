<template>
  <div class="message-input">
    <button @click="addImage">Send Image</button>
    <input v-model="message" placeholder="Type a message..." />
    <button @click="sendMessage">Send</button>
  </div>
</template>

<script>
export default {
  data() {
    return { message: '' };
  },
  methods: {
    sendMessage() {
      this.$emit('send', this.message);
      this.message = '';
    },
    addImage() {
      const url = prompt('Enter image URL:');
      if (url) {
        // Check if valid image URL
        const urlMatch = url.match(/https?:\/\/\S+\.(jpg|jpeg|png|gif|bmp|svg|webp|tiff|ico|apng|webm|mp4|jpeg)/i);
        if (urlMatch) {
          this.$emit('send', url); // Send URL
        } else {
          alert('Invalid image URL');
        }
      }
    }
  }
};
</script>

<style scoped>
.message-input { display: flex; gap: 10px; }
</style>