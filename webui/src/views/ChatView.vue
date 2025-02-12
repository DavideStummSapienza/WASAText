<template>
  <div class="chat-view">
    <h1 class="chat-title">{{this.$route.query.username}}</h1>
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
      currentUser: "currentUser" // TODO: set current user maybe not needed tho
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
        let photoUrl = "";

        // Überprüfen, ob der Inhalt eine URL zu einem Bild ist
        const urlMatch = content.match(/https?:\/\/.*\.(jpg|jpeg|png|gif|bmp|svg)/i);
        if (urlMatch) {
          isPhoto = true;
          photoUrl = urlMatch[0]; // Die URL extrahieren
          content = "[Bild]"; // Platzhaltertext, wenn es ein Bild ist
        }

        const newMessage = {
          message: content,
          is_photo: isPhoto,
          photo_url: photoUrl, // save Foto-URL 
        };

        const response = await axios.post(`/conversations/${partnerUsername}`, newMessage);
        this.messages.push(response.data); // Direkt zur Chat-Liste hinzufügen
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
