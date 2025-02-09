<template>
  <div class="incoming-message">
    <span class="timestamp">{{ timestamp }}</span>
    <span class="username">{{ username }}</span>
    <span class="content">{{ content }}</span>
    <button @click="react" class="reaction-button">+</button>
  </div>
</template>

<script>
export default {
  props: ['username', 'content', 'timestamp'],
  methods: {
    react() {
      alert('Reacting to message!');
    }
  }
};
</script>

<style scoped>
.incoming-message { background: #E6DFFF; padding: 10px; border-radius: 10px; }
</style>