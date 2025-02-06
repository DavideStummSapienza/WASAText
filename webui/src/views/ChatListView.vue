<template>
  <div>
    <h1>Chat List</h1>
    <p>Here will be the list of your chats.</p>

    <!-- Search bar -->
    <div>
      <input 
        type="text" 
        v-model="searchQuery" 
        @input="searchUsers" 
        placeholder="Search for users or groups" 
      />
    </div>

    <!-- Chat List -->
    <div v-if="conversations.length">
      <h2>Your Conversations</h2>
      <ul>
        <li v-for="conversation in conversations" :key="conversation.name">
          <div @click="viewConversation(conversation.name)">
            <img :src="conversation.photo_url" alt="Profile" />
            <p>{{ conversation.username }}</p>
            <p>{{ conversation.last_message }}</p>
            <p>{{ formatTimestamp(conversation.timestamp) }}</p> <!-- Display timestamp -->
          </div>
        </li>
      </ul>
    </div>

    <!-- Search Results for Users -->
    <div v-if="searchResults.length">
      <h2>Search Results</h2>
      <ul>
        <li v-for="user in searchResults" :key="user.username">
          <div @click="startNewConversation(user.username)">
            <p>{{ user.username }}</p>
          </div>
        </li>
      </ul>
    </div>

    <!-- Button to add new conversation or group -->
    <button @click="openCreateDialog">+</button>

    <!-- Dialog for creating conversation/group -->
    <div v-if="showCreateDialog">
      <div>
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
        const response = await axios.get('/api/conversations');
        this.conversations = response.data;
      } catch (error) {
        console.error('Error fetching conversations', error);
      }
    },

    // Search for users (using the backend API)
    async searchUsers() {
      if (!this.searchQuery) {
        this.searchResults = [];
        return;
      }

      try {
        const response = await axios.get(`/api/users`, {
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
        const response = await axios.post(`/api/conversations/${username}`, {
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
        const response = await axios.post('/api/groups', requestData);
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
/* Optional: Add styles for the chat list view */
button {
  margin-top: 20px;
}

div {
  margin-bottom: 10px;
}
</style>
