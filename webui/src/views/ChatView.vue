<template>
  <div class="chat-view">
    <h1>Chat</h1>
    <div v-for="msg in messages" :key="msg.id">
      <IncomingMessage v-if="msg.type === 'incoming'" :username="msg.username" :content="msg.content" :timestamp="msg.timestamp" />
      <OutgoingMessage v-else :content="msg.content" :timestamp="msg.timestamp" />
    </div>
    <MessageInput @send="handleSend" />
  </div>
</template>

<script>
import IncomingMessage from '@/components/IncomingMessage.vue';
import OutgoingMessage from '@/components/OutgoingMessage.vue';
import MessageInput from '@/components/MessageInput.vue';

export default {
  components: { IncomingMessage, OutgoingMessage, MessageInput },
  data() {
    return {
      messages: []
    };
  },
  methods: {
    handleSend(content) {
      this.messages.push({ id: Date.now(), content, timestamp: new Date().toLocaleTimeString(), type: 'outgoing' });
    }
  }
};
</script>

<style scoped>
.chat-view { background: #65558f; padding: 20px; }
</style>