<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-4 border-bottom">
      <h1 class="h2 conversation-header">{{ userName }}, search to start a chat</h1>
      
      <div class="btn-toolbar mb-2 mb-md-0 button-group-container">
        <button type="button" class="btn-chat-secondary" @click="newGroup">
          New group
        </button>

        <button type="button" class="btn-chat-tertiary" @click="logOut">
          Log Out
        </button>
      </div>
    </div>
    <div class="search-container">
  
      <form @submit.prevent="searchUsers" class="search-form">
        <input
          id="username"
          v-model="query"
          class="search-box"
          type="text"
          placeholder="Search by username"
        />
        <button class="search-button btn-chat-primary" type="submit">Search</button>
      </form>
      
      <div v-if="error" class="error-box">
        {{ error }}
      </div>
      
      <div v-if="loading">
        <LoadingSpinner />
      </div>
      
      <div v-if="!loading && showResults" class="results-section">
        <h2 class="contacts-title">Contacts</h2>
        
        <template v-if="users.length > 0">
          <div v-for="user in users" :key="user.id" class="user-card-new">
            <h5 class="user-name-new">
              @{{ user.name }} <span v-if="user.name === userName">(You)</span>
            </h5>
            <button
              class="text-button btn-chat-primary"
              @click="navigateToConversation(user.id, user.name)"
            >
              Start Chat
            </button>
          </div>
        </template>
        
        <template v-else>
          <p class="no-results">No users found matching "{{ lastQuery }}"</p>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import LoadingSpinner from "../components/LoadingSpinner.vue";

export default {
  name: "SearchPeopleView",
  components: {
    LoadingSpinner,
  },
  data() {
    return {
      userName: localStorage.getItem("name"),
      query: "",
      lastQuery: "",
      users: [],
      loading: false,
      showResults: false,
      error: "",
    };
  },
  methods: {
    async searchUsers() {
      if (!this.query.trim()) {
        this.error = "Please enter a valid search query.";
        this.showResults = false;
        return;
      }
      this.loading = true;
      this.error = "";
      this.users = [];
      this.showResults = false;
      try {
        const response = await axios.get(`/search`, {
          params: { username: this.query },
        });
        this.users = response.data;
        this.lastQuery = this.query;
        this.showResults = true;
      } catch (err) {
        const status = err.response?.status;
        const reason = err.response?.data?.message || "Failed to fetch users.";
        this.error = `Status ${status}: ${reason}`;
      } finally {
        this.loading = false;
      }
    },
    navigateToConversation(recipientId, recipientName) {
      this.error = "";
       
      localStorage.setItem("conversationName", recipientName);
      const senderId = localStorage.getItem("token");
      
      axios
        .post(`/conversations`, { senderId, recipientId })
        .then((response) => {
          const conversationId = response.data.conversationId;
          this.$router.push({
            path: `/conversations/${conversationId}`
          });
        })
        .catch((error) => {
          console.error("Error starting conversation:", error);
          if (error.response && error.response.data && error.response.data.error) {
            this.error = error.response.data.error; // "You cannot start a conversation with yourself"
          } else {
            this.error = "Failed to start conversation. Please try again.";
          }

        });
    },
    
    logOut() {
      localStorage.clear();
      this.$router.push({ path: "/" });
    },
    newGroup() {
      this.$router.push({ path: "/new-group" });
    }
  },
  mounted() {
    const token = localStorage.getItem("token");
    if (!token) {
      this.$router.push({ path: "/" });
    }
  }
};
</script>

<style scoped>
.conversation-header {
  font-size: 2rem;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}
.button-group-container {
    display: flex;
    gap: 10px;
}


.btn-chat-primary {
    background-color: #4a90e2;
    color: white;
    border: 1px solid #4a90e2;
    padding: 10px 18px;
    font-weight: 600;
    border-radius: 8px;
    transition: background-color 0.2s, box-shadow 0.2s;
    cursor: pointer;
}

.btn-chat-primary:hover {
    background-color: #3a74b6;
    box-shadow: 0 4px 10px rgba(74, 144, 226, 0.3);
}


.btn-chat-secondary {
    background-color: white;
    color: #4a90e2;
    border: 1px solid #4a90e2;
    padding: 10px 18px;
    font-weight: 600;
    border-radius: 8px;
    transition: background-color 0.2s, color 0.2s;
    cursor: pointer;
}

.btn-chat-secondary:hover {
    background-color: #e6f1fc;
    color: #3a74b6;
}


.btn-chat-tertiary {
    background-color: transparent;
    color: #7f8c8d;
    border: 1px solid #bdc3c7;
    padding: 10px 18px;
    font-weight: 500;
    border-radius: 8px;
    transition: background-color 0.2s, color 0.2s;
    cursor: pointer;
}

.btn-chat-tertiary:hover {
  background-color: #ecf0f1;
  color: #34495e;
  border-color: #34495e;
}
.search-container {
  text-align: center;
  padding: 40px 20px 20px;
  max-width: 700px;
  margin: 0 auto;
}

.search-form {
  display: flex;
  gap: 10px;
  justify-content: center;
  margin-bottom: 20px;
}

search-box {
  padding: 12px 15px;
  border: 1px solid #bdc3c7;
  border-radius: 8px;
  font-size: 16px;
  width: 100%; 
  max-width: 400px;
  transition: border-color 0.2s;
}


.search-box:focus {
  border-color: #4a90e2;
  outline: none;
}

.results-section {
  text-align: left;
  margin-top: 20px;
}

.contacts-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 15px;
  border-bottom: 2px solid #eee;
  padding-bottom: 5px;
}

.user-card-new {
  padding: 15px 20px;
  margin: 10px 0;
  background-color: #fff;
  border: 1px solid #e1e8ed;
  border-radius: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
  transition: all 0.2s;
}

.user-card-new:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}
.user-name-new {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 500;
  color: #2c3e50;
}

.text-button.btn-chat-primary {
  padding: 8px 15px;
  font-size: 14px;
}
.error-box {
  background-color: #f8d7da;
  color: #842029;
  border: 1px solid #f5c2c7;
  border-radius: 8px;
  padding: 10px;
  margin: 20px 0;
  text-align: center;
}



.no-results {
  font-size: 16px;
  color: #7f8c8d;
  margin-top: 20px;
  text-align: center;
}
</style>

