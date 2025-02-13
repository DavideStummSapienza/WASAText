<template>
  <div class="chat-list-container">
    <h1 class="chat-title">Chats</h1>
    <button class="profile-button" @click="goToProfile">Profile Settings</button>
    <div class="chat-list">
      <ChatCard v-for="chat in conversations" :key="chat.name" :chat="chat" />
    </div>
    <button class="add-chat-button" @click="openCreatedDialog">+</button>

    <!-- Dialog-Box -->
    <div v-if="showCreatedDialog" class="dialog-overlay">
      <div class="dialog-box">
        <h2>Create New</h2>
        <button class="dialog-button" @click="createConversation">Start New Conversation</button>
        <button class="dialog-button" @click="createGroup">Create New Group</button>
        <button class="dialog-button cancel" @click="closeCreatedDialog">Cancel</button>
      </div>
    </div>
  </div>
</template>


<script>
import ChatCard from "@/components/ChatCard.vue";
import axios from "@/services/axios";

export default {
  components: {
    ChatCard,
  },
  data() {
    return {
      conversations: [],
      showCreatedDialog: false,
      updateInterval: null,
    };
  },
  methods: {
    async fetchConversations() {
      try {
        const response = await axios.get("/user-profile");
        const newConversations = response.data || [];

        // Only reloads if there is new data
        if (JSON.stringify(this.conversations) !== JSON.stringify(newConversations)) {
          this.conversations = newConversations;
        }
      } catch (error) {
        console.error("Error fetching conversations", error);
        this.conversations = [];
      }
    },
    async goToProfile() {
      this.$router.push("/profile")
      //this.closeCreatedDialog();
    },
    openCreatedDialog() {
      this.showCreatedDialog = true;
    },
    async createConversation() {
      this.$router.push("/search")
      //this.closeCreatedDialog();
    },
    async createGroup() {
      this.$router.push("/create-group")
      //this.closeCreatedDialog();
    },
    closeCreatedDialog() {
      this.showCreatedDialog = false;
    },
  },
  mounted() {
    this.fetchConversations();

    this.updateInterval = setInterval(() => {
      this.fetchConversations();
    }, 5000);
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

.chat-list-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: #65558f;
  font-family: 'Roboto', sans-serif;
  color: white;
  position: relative;
}

.chat-title {
  font-size: 40px;
  margin-bottom: 20px;
}

.profile-button {
  position: absolute;
  top: 20px;
  right: 20px;
  background: #21005d;
  color: white;
  padding: 10px 20px;
  border-radius: 15px;
  cursor: pointer;
  border: none;
}

.chat-list {
  background: white;
  border-radius: 10px;
  padding: 10px;
  width: 80%;
  height: 70vh;
  overflow-y: auto;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.2);
}

.add-chat-button {
  position: absolute;
  bottom: 40px;
  right: 40px;
  background: #005047;
  color: white;
  font-size: 32px;
  border-radius: 50%;
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border: none;
}

/* Dialog Box */
.dialog-box {
  background: #21005d;
  padding: 20px;
  border-radius: 10px;
  text-align: center;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.3);
}

/* Dialog Buttons */
.dialog-button {
  background: #005047;
  color: white;
  border: none;
  padding: 10px 20px;
  margin: 10px;
  border-radius: 5px;
  cursor: pointer;
  font-size: 24px;
}

.dialog-button.cancel {
  background: #b22222;
}
</style>
