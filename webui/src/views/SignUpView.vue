<template>
  <div class="signUp-container">
    <h1>Sign Up</h1>
    <!-- Input field for the username -->
    <input v-model="username" placeholder="Enter your username" />
    <!-- Button to trigger the login function -->
    <button @click="signUp">Sign Up</button>
    <!-- Display an error message if there is an error -->
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<script>
import axios from "@/services/axios"; // Import the Axios instance for making HTTP requests

export default {
  // Data properties for the component
  data() {
    return {
      username: "", // Stores the username entered by the user
      error: null, // Stores any error message to be displayed to the user
    };
  },
  methods: {
    // Login function to handle user login
    async signUp() {
      // Check if username is provided
      if (!this.username) {
        this.error = "Username is required"; // Set error message if username is missing
        return;
      }
      try {
        
        // Make a POST request to the /session endpoint with the username
        const response = await axios.post("/session", { username: this.username });

        // safe Identifier in sessionStorage
        sessionStorage.setItem("identifier", response.data.identifier);

        // On success, redirect to the chat list page
        this.$router.push("/chats");
      } catch (err) {
        // Error handling for different error types
        if (err.response) {
          // Error coming from the backend (API response)
          if (err.response.data && err.response.data.error) {
            this.error = err.response.data.error; // Custom error message from the API
          } else {
            this.error = "An unexpected error occurred"; // General error message for unknown backend errors
          }
        } else if (err.request) {
          // Error due to network issues (e.g., no response from the server)
          this.error = "Network error, please try again later.";
        } else {
          // Other errors (e.g., misconfigured request or unexpected client-side issues)
          this.error = "An unexpected error occurred";
        }
      }
    },
  },
};
</script>

<style scoped>
/* Styling for the login container */
.signUp-container {
  max-width: 400px; /* Set max width for the container */
  margin: 100px auto; /* Center the container vertically and horizontally */
  text-align: center; /* Center text inside the container */
}

/* Styling for the input field */
input {
  width: 100%; /* Full width input */
  padding: 10px; /* Padding inside the input */
  margin: 10px 0; /* Margin above and below the input */
  border: 1px solid #ccc; /* Light gray border */
  border-radius: 5px; /* Rounded corners */
}

/* Styling for the login button */
button {
  width: 100%; /* Full width button */
  padding: 10px; /* Padding inside the button */
  background-color: #007bff; /* Blue background */
  color: white; /* White text */
  border: none; /* Remove border */
  border-radius: 5px; /* Rounded corners */
  cursor: pointer; /* Pointer cursor on hover */
}

/* Hover effect for the button */
button:hover {
  background-color: #0056b3; /* Darker blue on hover */
}

/* Styling for error message */
.error {
  color: red; /* Red color for error messages */
  margin-top: 10px; /* Margin above the error message */
}
</style>
