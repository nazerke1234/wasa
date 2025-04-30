<template>
  <div class="search-container">
    <div class="search-box">
      <h2>User Search</h2>
      <input v-model="searchQuery" @input="debouncedSearch" placeholder="Enter name..." class="search-input" />

      <ul v-if="users.length" class="user-list">
        <li v-for="user in users" :key="user.ID" class="user-item">
          <div class="user-info" @click="selectUser(user)">
            <img v-if="user.photo" :src="user.photo" alt="Profile Photo" class="user-photo" />
            <span>{{ user.username }}</span>
          </div>
        </li>
      </ul>

      <p v-else-if="searchQuery.length >= 2" class="no-users">❌ Users not found</p>

      <div v-if="selectedUsers.length > 0" class="selected-users">
      <h3>Selected Users:</h3>
      <div class="user-cards">
        <div v-for="user in selectedUsers" :key="user.ID" class="user-card">
          <span>{{ user.username }}</span>
          <button @click="removeUserFromSelection(user)" class="delete-button">❌
            <i class="fas fa-trash"></i>
          </button>
        </div>
      </div>

      <button v-if="selectedUsers.length > 1" @click="createGroupChat" class="action-button">Create Group Chat</button>
      <button v-if="selectedUsers.length === 1" @click="createChatWithUser(selectedUsers[0])" class="action-button">Create Personal Chat</button>
    </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import _ from 'lodash';

export default {
  data() {
    return {
      searchQuery: '',
      users: [],
      selectedUsers: []
    };
  },
  methods: {
    async searchUsers() {
      if (this.searchQuery.length >= 2) {
        const token = localStorage.getItem('access_token')?.trim();
        if (!token) {
          console.error('Token missing! User is not authorized.');
          return;
        }
        axios.get(`http://localhost:8080/api/users?name=${this.searchQuery}`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        .then(response => {
          this.users = response.data;
        })
        .catch(error => {
          console.error('Error searching users:', error);
        });
      }
    },
    selectUser(user) {
      if (!this.selectedUsers.some(u => u.ID === user.ID)) {
        this.selectedUsers.push(user);
      }
    },
    removeUserFromSelection(user) {
      this.selectedUsers = this.selectedUsers.filter(u => u.ID !== user.ID);
    },
    createGroupChat() {
      const token = localStorage.getItem('access_token')?.trim();
      if (!token) return;
      if (this.selectedUsers.length < 2) {
        alert('At least two users are required for a group chat!');
        return;
      }
      const userIDs = this.selectedUsers.map(user => user.ID);
      const currentUserId = localStorage.getItem('user_id');
      if (currentUserId && !userIDs.includes(Number(currentUserId))) {
        userIDs.push(Number(currentUserId));
      }
      axios.post('http://localhost:8080/api/groupChats', { user_ids: userIDs, is_group: true }, { headers: { Authorization: `Bearer ${token}` } })
      .then(response => {
        alert('Group chat created!');
        this.$router.push(`/groupChats/${response.data.chat_id}`);
      })
      .catch(error => console.error('Error creating group chat:', error));
    },
    createChatWithUser(user) {
      const token = localStorage.getItem('access_token')?.trim();
      if (!token) return;
      localStorage.setItem('recipient_id', JSON.stringify(user.ID));
      axios.post('http://localhost:8080/api/personalChats', { user_id: user.ID, is_group: false }, { headers: { Authorization: `Bearer ${token}` } })
      .then(response => {
        if (response.data.message === "Chat already exists") {
          alert(`A chat with user ${user.username} already exists!`);
        } else {
          alert(`Chat with user ${user.username} created!`);
          this.$router.push(`/chats/${response.data.chat_id}`);
        }
      })
      .catch(error => console.error('Error creating chat:', error));
    }
  },
  created() {
    this.debouncedSearch = _.debounce(this.searchUsers, 500);
  }
};
</script>

<style scoped>
.search-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #1e3c3c, #2a5d5d);
  font-family: 'Poppins', sans-serif;
}
.search-box {
  background: white;
  padding: 30px;
  border-radius: 15px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  text-align: center;
  width: 400px;
}
.search-input {
  width: 90%;
  padding: 10px;
  border: 2px solid #2a5d5d;
  border-radius: 5px;
  margin-bottom: 15px;
}
.user-list {
  list-style: none;
  padding: 0;
}
.user-item {
  display: flex;
  justify-content: space-between;
  padding: 5px;
  background: white;
  border-radius: 5px;
  margin: 5px 0;
  cursor: pointer;
}
.selected-users {
  margin-top: 20px;
}
.user-cards {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.user-card span {
  font-size: 16px;
  font-weight: 500;
  margin-right: 10px;
}
.user-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
}

.delete-button {
  background: none;
  border: none;
  cursor: pointer;
}

.action-button {
  background: #2a5d5d;
  color: white;
  border: none;
  padding: 10px;
  border-radius: 5px;
  cursor: pointer;
}
.action-button:hover {
  background: #1e3c3c;
}
.delete-button i {
  color: white; /* White color for the cross icon */
  font-size: 18px; /* Adjust icon size for clarity */
}

.delete-button:focus {
  outline: none; /* Remove the focus outline */
}

.user-card {
  display: flex; /* Align text and button in a row */
  justify-content: space-between; /* Space between username and delete button */
  align-items: center; /* Center vertically */
  padding: 10px;
  margin-bottom: 8px;
  background-color: #f9f9f9; /* Light background */
  border-radius: 10px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1); /* Subtle shadow for depth */
}

</style>
