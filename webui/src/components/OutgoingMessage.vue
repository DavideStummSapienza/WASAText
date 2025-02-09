<template>
  <div class="outgoing-message">
    <span class="content">{{ content }}</span>
    <span class="timestamp">{{ timestamp }}</span>
    <button @click="react" class="reaction-button">+</button>
  </div>
</template>

<script>
export default {
  props: ['content', 'timestamp'],
  methods: {
    react() {
      alert('Reacting to message!');
    }
  }
};
</script>

<style scoped>
.outgoing-message { background: #D6C5F0; padding: 10px; border-radius: 10px; text-align: right; }
</style>