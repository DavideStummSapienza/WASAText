<template>
  <div class="chat-list-container">
    <h1 class="chat-title">Chats</h1>
    <button class="profile-button">Profile Settings</button>
    <div class="chat-list">
      <ChatCard v-for="chat in conversations" :key="chat.id" :chat="chat" />
    </div>
    <button class="add-chat-button">+</button>
  </div>
</template>


<script>
import ChatCard from "@/components/ChatCard.vue";
import axios from "@/services/axios";

export default {
  components: {
    ChatCard, // Register the component
  },
  data() {
    return {
      searchQuery: '', // Search query input
      conversations: [], // List of user conversations
      searchResults: [], // Search results for users
      showCreateDialog: false, // Flag for showing the creation dialog
    };
  },
  methods: {
    // Fetch user's conversations from the API
    async fetchConversations() {
      try {
        const response = await axios.get("/user-profile");
        this.conversations = response.data || [];
      } catch (error) {
        console.error('Error fetching conversations', error);
        this.conversations = [];
      }
    },

    // Search for users (using the backend API)
    async searchUsers() {
      if (!this.searchQuery) {
        this.searchResults = [];
        return;
      }

      try {
        const response = await axios.get("/users", {
          params: { username: this.searchQuery }, // Pass the query parameter
        });
        this.searchResults = response.data;
      } catch (error) {
        console.error('Error searching users', error);
      }
    },

    // Start a new conversation with a user
    async startNewConversation(username) {
      try {
        const response = await axios.post(`/conversations/${username}`, {
          message: "Hi!", // Default message
        });
        this.$router.push({ name: 'conversation', params: { username } });
      } catch (error) {
        console.error('Error starting new conversation', error);
      }
    },

    // View an existing conversation
    viewConversation(conversationName) {
      this.$router.push({ name: 'conversation', params: { username: conversationName } });
    },

    // Show the dialog to create a new conversation or group
    openCreateDialog() {
      this.showCreateDialog = true;
    },

    // Create a new conversation
    async createConversation() {
      try {
        console.log('Creating a new conversation');
        // Logic to create a new conversation (expand as needed)
        this.closeCreateDialog();
      } catch (error) {
        console.error('Error creating conversation', error);
      }
    },

    // Create a new group
    async createGroup() {
      try {
        const groupName = prompt("Enter group name:");
        const selectedUsers = prompt("Enter usernames of users to add (comma separated):");

        // Prepare the request data
        const requestData = {
          groupName: groupName,
          names: selectedUsers.split(",").map(name => name.trim())
        };

        // Send a request to add users to the group
        const response = await axios.post('/groups', requestData);
        console.log(response.data.message); // Log the success message
        this.closeCreateDialog();
      } catch (error) {
        console.error('Error creating group', error);
      }
    },

    // Close the dialog
    closeCreateDialog() {
      this.showCreateDialog = false;
    },

    // Format timestamp to a readable string
    formatTimestamp(timestamp) {
      const date = new Date(timestamp);
      return date.toLocaleString(); // Formats the timestamp as a local string
    },
  },
  mounted() {
    this.fetchConversations();
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
</style>
