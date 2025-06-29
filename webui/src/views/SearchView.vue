<template>
  <div class="search-wrapper">
    <div class="toolbar">
      <h2>{{ displayName }}, search and chat</h2>
      <div class="toolbar-actions">
        <button @click="reload">Refresh</button>
        <button @click="signOut" class="logout-btn">Log Out</button>
        <button @click="goCreateGroup" class="create-btn">New Group</button>
      </div>
    </div>

    <div class="search-section">
      <form @submit.prevent="findUsers" class="search-form">
        <input
          v-model="query"
          type="text"
          class="search-box"
          placeholder="Search by username"
        />
        <button type="submit" class="search-btn">Search</button>
      </form>

      <div v-if="error" class="error-box">{{ error }}</div>
      <LoadingSpinner v-if="loading" />

      <div v-if="!loading && showResults" class="results">
        <h3>Results:</h3>
        <template v-if="users.length">
          <div
            v-for="user in users"
            :key="user.id"
            class="user-card"
          >
            <span class="username">@{{ user.name }}</span>
            <button @click="startConversation(user.id, user.name)" class="chat-btn">
              Chat
            </button>
          </div>
        </template>
        <template v-else>
          <p class="no-results">No users found for "{{ lastQuery }}"</p>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import axios from '../services/axios';
import LoadingSpinner from '../components/LoadingSpinner.vue';

export default {
  name: 'UserSearchView',
  components: { LoadingSpinner },
  data() {
    return {
      displayName: localStorage.getItem('name') || '',
      query: '',
      lastQuery: '',
      users: [],
      loading: false,
      showResults: false,
      error: ''
    };
  },
  mounted() {
    const token = localStorage.getItem('token');
    if (!token) this.$router.push('/');
  },
  methods: {
    async findUsers() {
      if (!this.query.trim()) {
        this.error = 'Enter a name to search.';
        this.showResults = false;
        return;
      }
      this.error = '';
      this.loading = true;
      try {
        const res = await axios.get('/search', {
          params: { username: this.query }
        });
        this.users = res.data;
        this.lastQuery = this.query;
        this.showResults = true;
      } catch (e) {
        console.error(e);
        this.error = 'Unable to search. Try again.';
      } finally {
        this.loading = false;
      }
    },
    async startConversation(userId, userName) {
      const senderId = localStorage.getItem('token');
      localStorage.setItem('conversationName', userName);
      try {
        const res = await axios.post('/conversations', { senderId, recipientId: userId });
        const convoId = res.data.conversationId;
        this.$router.push(`/conversations/${convoId}`);
      } catch (e) {
        console.error('Conversation error:', e);
      }
    },
    reload() {
      this.$router.push('/search');
    },
    signOut() {
      localStorage.clear();
      this.$router.push('/');
    },
    goCreateGroup() {
      this.$router.push('/new-group');
    }
  }
};
</script>

<style scoped>
.search-wrapper {
  max-width: 800px;
  margin: auto;
  padding: 2rem;
}
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}
.toolbar-actions {
  display: flex;
  gap: 0.5rem;
}
button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  background-color: #007bff;
  color: white;
  cursor: pointer;
}
button:hover {
  background-color: #0056b3;
}
.logout-btn {
  background-color: #6c757d;
}
.logout-btn:hover {
  background-color: #5a6268;
}
.create-btn {
  background-color: #28a745;
}
.create-btn:hover {
  background-color: #218838;
}
.search-section {
  text-align: center;
}
.search-form {
  display: flex;
  justify-content: center;
  margin-bottom: 1rem;
}
.search-box {
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 6px;
  font-size: 16px;
  width: 300px;
}
.search-btn {
  margin-left: 0.5rem;
  padding: 0.5rem 1rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}
.error-box {
  color: red;
  margin-top: 1rem;
}
.results {
  margin-top: 1.5rem;
}
.user-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 0.75rem;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}
.username {
  font-weight: bold;
  color: #007bff;
}
.chat-btn {
  background-color: #17a2b8;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
}
.chat-btn:hover {
  background-color: #138496;
}
.no-results {
  color: #666;
  font-size: 14px;
  margin-top: 1rem;
}
</style>
