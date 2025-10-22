<template>
  <div>
    <div class="page-header">
      <h1 class="page-heading">{{ userName }}, find and message users!</h1>
      <div class="action-controls">
        <div class="control-group">
          <button type="button" class="control-btn secondary" @click="reloadPage">Reload</button>
          <button type="button" class="control-btn secondary" @click="signOut">Sign Out</button>
        </div>
        <div class="control-group">
          <button type="button" class="control-btn primary" @click="createGroup">Create Group</button>
        </div>
      </div>
    </div>

    <div class="search-interface">
      <form @submit.prevent="findUsers" class="search-form">
        <input
          id="user-search"
          v-model="searchTerm"
          class="search-input"
          type="text"
          placeholder="Enter username to search"
        />
        <button class="search-submit-btn" type="submit">Find Users</button>
      </form>
      <div v-if="errorMessage" class="error-alert">
        {{ errorMessage }}
      </div>
      <div v-if="isLoading">
        <LoadingSpinner />
      </div>
      <div v-if="!isLoading && displayResults" class="results-container">
        <h2 class="results-heading">Search Results:</h2>
        <template v-if="foundUsers.length > 0">
          <div v-for="user in foundUsers" :key="user.id" class="user-result-card">
            <h5 class="user-username">
              @{{ user.name }}
            </h5>
            <button
              class="message-btn"
              @click="startChat(user.id, user.name)"
            >
              Message
            </button>
          </div>
        </template>
        <template v-else>
          <p class="no-results-message">No users found for "{{ previousSearchTerm }}"</p>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import LoadingSpinner from "../components/LoadingSpinner.vue";

export default {
  name: "SearchView",
  components: {
    LoadingSpinner,
  },
  data() {
    return {
      userName: localStorage.getItem("name"),
      searchTerm: "",
      previousSearchTerm: "",
      foundUsers: [],
      isLoading: false,
      displayResults: false,
      errorMessage: "",
    };
  },
  methods: {
    async findUsers() {
      if (!this.searchTerm.trim()) {
        this.errorMessage = "Please enter a valid search term.";
        this.displayResults = false;
        return;
      }
      this.isLoading = true;
      this.errorMessage = "";
      this.foundUsers = [];
      this.displayResults = false;
      try {
        const apiResponse = await axios.get(`/search`, {
          params: { username: this.searchTerm },
        });
        this.foundUsers = apiResponse.data;
        this.previousSearchTerm = this.searchTerm;
        this.displayResults = true;
      } catch (apiError) {
        const errorCode = apiError.response?.status;
        const errorDetails = apiError.response?.data?.message || "Unable to retrieve users.";
        this.errorMessage = `Error ${errorCode}: ${errorDetails}`;
      } finally {
        this.isLoading = false;
      }
    },
    async startChat(recipientIdentifier, recipientDisplayName) {
      localStorage.setItem("conversationName", recipientDisplayName);
      const senderIdentifier = localStorage.getItem("token");
      try {
        const apiResponse = await axios.post(`/chats`, { 
          senderId: senderIdentifier, 
          recipientId: recipientIdentifier 
        });
        const chatId = apiResponse.data.conversationId;
        this.$router.push({
          path: `/chats/${chatId}`
        });
      } catch (apiError) {
        console.error("Error initiating conversation:", apiError);
        this.errorMessage = "Unable to start conversation. Please try again.";
      }
    },
    reloadPage() {
      this.$router.push({ path: "/search" });
    },
    signOut() {
      localStorage.clear();
      this.$router.push({ path: "/" });
    },
    createGroup() {
      this.$router.push({ path: "/new-group" });
    }
  },
  mounted() {
    const authToken = localStorage.getItem("token");
    if (!authToken) {
      this.$router.push({ path: "/" });
    }
  }
};
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  align-items: center;
  padding-top: 1rem;
  padding-bottom: 0.5rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid #e0e0e0;
}

.page-heading {
  font-size: 1.75rem;
  font-weight: 600;
  color: #333;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.action-controls {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.control-group {
  display: flex;
  gap: 0.5rem;
}

.control-btn {
  padding: 0.375rem 0.75rem;
  border: 1px solid;
  border-radius: 0.375rem;
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.2s ease;
}

.control-btn.secondary {
  background-color: transparent;
  color: #6b7280;
  border-color: #6b7280;
}

.control-btn.secondary:hover {
  background-color: #6b7280;
  color: white;
}

.control-btn.primary {
  background-color: transparent;
  color: #3b82f6;
  border-color: #3b82f6;
}

.control-btn.primary:hover {
  background-color: #3b82f6;
  color: white;
}

.search-interface {
  text-align: center;
  padding: 2rem;
  max-width: 600px;
  margin: 0 auto;
}

.search-form {
  display: flex;
  justify-content: center;
  margin-bottom: 1.5rem;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.search-input {
  padding: 0.75rem 1rem;
  border: 2px solid #e5e7eb;
  border-radius: 0.5rem;
  font-size: 1rem;
  width: 300px;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.search-submit-btn {
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #3b82f6 0%, #1e40af 100%);
  color: #fff;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.search-submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
}

.error-alert {
  background-color: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
  border-radius: 0.5rem;
  padding: 1rem;
  margin: 1.5rem 0;
  text-align: center;
}

.results-container {
  margin-top: 2rem;
}

.results-heading {
  font-size: 1.5rem;
  font-weight: 600;
  color: #374151;
  margin-bottom: 1rem;
}

.user-result-card {
  padding: 1rem;
  margin: 0.75rem 0;
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 0.75rem;
  font-size: 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.user-result-card:hover {
  background-color: #f9fafb;
  transform: translateY(-1px);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.user-username {
  margin: 0;
  font-size: 1.125rem;
  color: #3b82f6;
  font-weight: 500;
}

.user-username:hover {
  text-decoration: underline;
  cursor: pointer;
}

.message-btn {
  padding: 0.5rem 1rem;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: #fff;
  border: none;
  border-radius: 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.message-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(16, 185, 129, 0.4);
}

.no-results-message {
  font-size: 1rem;
  color: #6b7280;
  margin-top: 1.5rem;
  font-style: italic;
}

@media (max-width: 640px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .action-controls {
    width: 100%;
    justify-content: space-between;
  }
  
  .search-form {
    flex-direction: column;
    align-items: center;
  }
  
  .search-input {
    width: 100%;
    max-width: 300px;
  }
  
  .user-result-card {
    flex-direction: column;
    gap: 0.75rem;
    text-align: center;
  }
}
</style>


