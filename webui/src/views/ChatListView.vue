<template>
  <div>
    <h1 class="header">Chat List</h1>
    

    <!-- Search bar -->
    <div class="search-bar">
      <input 
        type="text" 
        v-model="searchQuery" 
        @input="searchUsers" 
        placeholder="Search for users or groups" 
      />
    </div>

    <!-- Chat List -->
    <div v-if="conversations.length" class="chat-list">
      <h2>Your Conversations</h2>
      <ul>
        <li v-for="conversation in conversations" :key="conversation.name" class="conversation-item">
          <div @click="viewConversation(conversation.name)" class="conversation">
            <img :src="conversation.photo_url" alt="Profile" class="avatar" />
            <div class="conversation-details">
              <p class="username">{{ conversation.username }}</p>
              <p class="last-message">{{ conversation.last_message }}</p>
              <p class="timestamp">{{ formatTimestamp(conversation.timestamp) }}</p>
            </div>
          </div>
        </li>
      </ul>
    </div>

    <!-- No conversations message -->
    <div v-else class="no-conversations">
      <p>No conversations yet</p>
    </div>

    <!-- Search Results for Users -->
    <div v-if="searchResults.length" class="search-results">
      <h2>Search Results</h2>
      <ul>
        <li v-for="user in searchResults" :key="user.username">
          <div @click="startNewConversation(user.username)" class="search-item">
            <p class="search-username">{{ user.username }}</p>
          </div>
        </li>
      </ul>
    </div>

    <!-- Button to add new conversation or group -->
    <button @click="openCreateDialog" class="create-button">+</button>

    <!-- Dialog for creating conversation/group -->
    <div v-if="showCreateDialog" class="create-dialog">
      <div class="dialog-actions">
        <button @click="createConversation">Create Conversation</button>
        <button @click="createGroup">Create Group</button>
        <button @click="closeCreateDialog">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script>

import axios from "@/services/axios";

export default {
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
/* Basic Styles */
body {
  font-family: Arial, sans-serif;
  background-color: #f4f4f9;
  color: #333;
  margin: 0;
  padding: 0;
}

.header {
  text-align: center;
  margin-top: 20px;
  font-size: 2em;
}

.description, .no-conversations {
  text-align: center;
  margin-bottom: 20px;
  font-size: 1.2em;
}

.search-bar {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.search-bar input {
  width: 80%;
  padding: 10px;
  font-size: 1em;
  border-radius: 5px;
  border: 1px solid #ddd;
}

.chat-list, .search-results {
  margin-top: 20px;
  margin-left: 20px;
  margin-right: 20px;
}

.chat-list ul, .search-results ul {
  list-style-type: none;
  padding: 0;
}

.conversation-item {
  display: flex;
  align-items: center;
  border-bottom: 1px solid #ccc;
  padding: 10px 0;
}

.conversation {
  display: flex;
  align-items: center;
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 15px;
}

.conversation-details {
  display: flex;
  flex-direction: column;
}

.username {
  font-weight: bold;
}

.last-message {
  color: #777;
}

.timestamp {
  font-size: 0.8em;
  color: #aaa;
}

.search-item {
  padding: 10px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
}

.create-button {
  position: fixed;
  bottom: 30px;
  right: 30px;
  padding: 15px;
  font-size: 1.5em;
  background-color: #007BFF;
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
}

.create-button:hover {
  background-color: #0056b3;
}

.create-dialog {
  position: fixed;
  bottom: 50px;
  left: 50%;
  transform: translateX(-50%);
  background-color: white;
  border: 1px solid #ddd;
  padding: 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.dialog-actions button {
  padding: 10px 20px;
  margin: 5px;
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.dialog-actions button:hover {
  background-color: #218838;
}
</style>
