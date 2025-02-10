<template>
  <div class="chat-view">
    <h1>Chat</h1>
    <div v-for="msg in messages" :key="msg.message_id">
      <IncomingMessage 
        v-if="msg.sender === this.$route.query.username" 
        :username="msg.sender"
        :content="msg.content"
        :timestamp="msg.timestamp"
        :is-photo="msg.is_photo"
        :photo-url="msg.photo_url"
        :is-forwarded="msg.is_forwarded"
        :reactions="msg.reactions"
      />
      <OutgoingMessage 
        v-else 
        :content="msg.content" 
        :timestamp="msg.timestamp"
        :is-photo="msg.is_photo"
        :photo-url="msg.photo_url"
        :is-forwarded="msg.is_forwarded"
        :reactions="msg.reactions"
        :fully-received="msg.fully_received"
        :fully-read="msg.fully_read"
      />
    </div>
    <MessageInput @send="handleSend" />
  </div>
</template>

<script>
import IncomingMessage from '@/components/IncomingMessage.vue';
import OutgoingMessage from '@/components/OutgoingMessage.vue';
import MessageInput from '@/components/MessageInput.vue';
import axios from "@/services/axios";

export default {
  components: { IncomingMessage, OutgoingMessage, MessageInput },
  data() {
    return {
      messages: [],
      currentUser: "currentUser" // TODO: Hier den aktuellen Nutzer setzen
    };
  },
  methods: {
    async fetchMessages() {
      try {
        const response = await axios.get("/chat", {
          params: { username: this.$route.query.username },
        });
        this.messages = response.data;
      } catch (error) {
        console.error("Error fetching messages:", error);
      }
    },
    async handleSend(content) {
      try {
        const newMessage = {
          fromUser: this.currentUser,
          toUser: this.$route.query.username,
          content,
          is_photo: false,
          photo_url: "",
          is_forwarded: false,
        };
        const response = await axios.post("/send-message", newMessage);
        this.messages.push(response.data); // Direkt hinzufügen
      } catch (error) {
        console.error("Error sending message:", error);
      }
    }
  },
  mounted() {
    this.fetchMessages();
  }
};
</script>

<style scoped>
.chat-view { background: #65558f; padding: 20px; }
</style>
