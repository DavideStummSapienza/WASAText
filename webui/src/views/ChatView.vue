<template>
  <div class="chat-view">
    <!-- Partner Username as Title -->
    <h1 class="chat-title">{{this.$route.query.username}}</h1>

    <!-- Show the messages List -->
    <div v-for="msg in messages" :key="msg.message_id">
      <IncomingMessage 
        v-if="msg.sender === this.$route.query.username" 
        :username="msg.sender"
        :content="msg.content"
        :timestamp="msg.timestamp"
        :is-photo="msg.is_photo"
        :is-forwarded="msg.is_forwarded"
        :reactions="msg.reactions"
      />
      <OutgoingMessage 
        v-else 
        :content="msg.content" 
        :timestamp="msg.timestamp"
        :is-photo="msg.is_photo"
        :is-forwarded="msg.is_forwarded"
        :reactions="msg.reactions"
        :fully-received="msg.fully_received"
        :fully-read="msg.fully_read"
      />
    </div>

    <!-- New Message Input -->
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
    };
  },
  methods: {
    async fetchMessages() {
      try {
        const partnerUsername = this.$route.query.username;
        
        const response = await axios.get(`/conversations/${partnerUsername}`);
        this.messages = response.data;
      } catch (error) {
        console.error("Error fetching messages:", error);
      }
    },
    async handleSend(content) {
      try {
        const partnerUsername = this.$route.query.username;
        let isPhoto = false;

        // Überprüfen, ob der Inhalt eine URL zu einem Bild ist
        const urlMatch = content.match(/https?:\/\/.*\.(jpg|jpeg|png|gif|bmp|svg)/i);
        if (urlMatch) {
          isPhoto = true;
          content = urlMatch[0];
        }

        const newMessage = {
          message: content,
          is_photo: isPhoto, 
        };

        const response = await axios.post(`/conversations/${partnerUsername}`, newMessage);
        if (!this.messages) {
          this.messages = [];
        }
    
        this.messages.push(response.data);
        
      } catch (error) {
        console.error("Error sending message:", error);
      }
    },
  },
  mounted() {
    this.fetchMessages();
  }
};
</script>

<style scoped>
.chat-view { background: #65558f; padding: 20px; }
</style>
