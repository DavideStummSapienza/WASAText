<template>
  <div class="message-input">
    <button @click="addImage">📎</button>
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
        this.$emit('send', `<img src='${url}' />`);
      }
    }
  }
};
</script>

<style scoped>
.message-input { display: flex; gap: 10px; }
</style>