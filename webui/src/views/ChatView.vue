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
        @reaction-added="handleReaction(msg.message_id, $event)"
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
        @reaction-added="handleReaction(msg.message_id, $event)"
      />
    </div>

    <!-- New Message Input -->
    <MessageInput @send="handleSend" @send-image="handleImageUpload" />

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
      updateInterval: null,
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

    async handleImageUpload(imageFile) {
      try {
        const formData = new FormData();
        formData.append("image", imageFile);

        // Upload an das Backend
        const uploadResponse = await axios.post("/upload", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        const imageUrl = uploadResponse.data.imageUrl;

        // Bild als Nachricht senden
        this.handleSend(imageUrl);
      } catch (error) {
        console.error("Error uploading image:", error);
      }
    },

    async handleReaction(messageId, reaction) {

      // Find message and add Reaktion
      const message = this.messages.find(msg => msg.message_id === messageId);

      if (!message) {
        console.error(`Message with ID ${messageId} not found!`);
        return;
      }

      // Sicherstellen, dass `reactions` existiert
      if (!message.reactions) {
        message.reactions = [];
      }

      // Reaktion hinzufügen
      message.reactions.push(reaction);

      // Send reaction to server
      try {

        await axios.put(`/conversations/messages/${messageId}/comment`, reaction);

        await this.fetchMessages();

      } catch (error) {

        console.error("Error sending reaction:", error);

      }
    },

  },

  mounted() {
    this.fetchMessages();

    this.updateInterval = setInterval(() => {
      this.fetchMessages();
    }, 10000);
  },

  beforeUnmount() {
    // Delete the Interval
    if (this.updateInterval) {
      clearInterval(this.updateInterval);
    }
  },
};
</script>

<style scoped>
.chat-view { background: #65558f; padding: 20px; }
</style>
