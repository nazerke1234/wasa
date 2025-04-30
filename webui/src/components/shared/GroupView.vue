<template>
  <div class="chat-container">
    <!-- Group header with name and photo -->
    <div class="chat-header">
      <div class="group-info">
        <img v-if="groupPhoto" :src="groupPhoto" alt="Group Photo" class="group-photo" />
        <h2>{{ groupName }} ({{ usersCount }} participants)</h2>
      </div>
      <div class="header-actions">
        <button @click="toggleEditName">Change Name</button>
        <button @click="toggleEditPhoto">Change Photo</button>
        <button v-if="showExitButton" class="exit-btn" @click="confirmExitGroup">Exit</button>
      </div>
    </div>

    <!-- Name edit form -->
    <div v-if="isEditingName" class="edit-box">
      <input v-model="newGroupName" type="text" placeholder="New group name" />
      <button @click="updateGroupName" class="btn-success">Save</button>
      <button @click="isEditingName = false" class="btn-danger">Cancel</button>
    </div>

    <!-- Photo edit form -->
    <div v-if="isEditingPhoto" class="edit-box">
      <input v-model="newPhotoUrl" type="text" placeholder="URL of new photo" />
      <button @click="saveGroupPhoto" class="btn-success">Save</button>
      <button @click="cancelPhotoEdit" class="btn-danger">Cancel</button>
    </div>

    <!-- User search -->
    <div class="search-section">
      <button @click="showSearch = !showSearch" v-if="!showSearch">➕ Add user</button>
      <div v-if="showSearch">
        <input v-model="searchQuery" @input="debouncedSearch" placeholder="Search users..." />
      </div>
    </div>

    <!-- Found users -->
    <ul v-if="users.length" class="user-list">
      <li v-for="user in users" :key="user.ID">
        <img v-if="user.photo" :src="user.photo" class="user-photo" />
        <span @click="selectUser(user)">{{ user.username }}</span>
      </li>
    </ul>

    <!-- Selected users -->
    <div v-if="selectedUsers.length" class="selected-users">
      <h3>Selected Users:</h3>
      <ul>
        <li v-for="user in selectedUsers" :key="user.ID">
          {{ user.username }} <button @click="removeUserFromSelection(user)">❌</button>
        </li>
      </ul>
      <button @click="addUsersToGroup" :disabled="selectedUsers.length === 0">✅ Add</button>
    </div>

    <!-- Chat messages -->
    <div class="messages">
      <div
        v-for="message in messages"
        :key="message.ID"
        :class="{ 'my-message': message.isMine, 'other-message': !message.isMine }"
      >
        <div class="message-content">
          <span>{{ message.Content }}</span>

          <button v-if="message.isMine" @click="deleteMessage(message.ID)" class="msg-action-btn">
            <font-awesome-icon icon="trash" />
          </button>

          <button @click="message.showForwardInput = !message.showForwardInput" class="msg-action-btn">
            <font-awesome-icon icon="share" />
          </button>

          <!-- Comment add/remove button -->
          <button v-if="!message.showCommentInput" @click="toggleCommentInput(message)" class="msg-action-btn">
            <font-awesome-icon :icon="message.Comment ? 'times' : 'comment'" />
          </button>
        </div>

        <!-- Forward message input field -->
        <div v-if="message.showForwardInput" class="forward-box">
          <input v-model="message.forwardChatName" type="text" placeholder="Chat name..." />
          <button @click="forwardMessage(message)">
            <font-awesome-icon icon="paper-plane" />
          </button>
        </div>

        <!-- Comment input field -->
        <div v-if="message.showCommentInput" class="comment-box">
          <input v-model="message.newComment" type="text" placeholder="Enter comment..." />
          <button @click="addOrRemoveComment(message)">
            <font-awesome-icon icon="check" />
          </button>
        </div>

        <!-- Display comment -->
        <div v-if="message.Comment" class="comment">
          <strong>Comment:</strong> {{ message.Comment }}
        </div>
      </div>
    </div>

    <!-- Message input form -->
    <div class="message-input">
      <input 
        v-model="newMessage"
        type="text"
        placeholder="Enter message..."
        @keyup.enter="sendMessage"
      />
      <button @click="sendMessage">
        <font-awesome-icon icon="paper-plane" />
      </button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import _ from 'lodash';
import { messageMixin } from '../../mixins/messageMixin';
import { library } from "@fortawesome/fontawesome-svg-core";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { faTrash, faShare, faComment, faTimes, faCheck, faPaperPlane } from "@fortawesome/free-solid-svg-icons";

library.add(faTrash, faShare, faComment, faTimes, faCheck, faPaperPlane);

export default {
  components: { FontAwesomeIcon },
  mixins: [messageMixin],
  data() {
    return {
      searchQuery: '',
      users: [],
      selectedUsers: [],
      chatId: this.$route.params.chatId, // Get chatId from URL
      showSearch: false,  // Flag for search box visibility
      isEditingName: false, // Flag for editing group name form
      newGroupName: '', // Group name
      usersCount: 0, // Number of participants
      isEditingPhoto: false, // Flag for photo editing
      newPhotoUrl: '',
      showExitButton: true, // Flag for showing "Exit group" button
      groupPhoto: '', // Group photo
      newMessage: '', // New message text
      messages: [] 
    };
  },
  methods: {
    // Toggle comment input visibility
    toggleCommentInput(message) {
      message.showCommentInput = !message.showCommentInput;
    },

    // Add or remove comment
    async addOrRemoveComment(message) {
      if (message.Comment) {
        await this.uncommentMessage(message.ID);
      } else {
        await this.commentMessage(message.ID, message.newComment);
      }
      message.showCommentInput = false;
    },
    toggleEditName() {
      this.isEditingName = !this.isEditingName;
    },
    toggleEditPhoto() {
      this.isEditingPhoto = !this.isEditingPhoto;
    },
    async updateGroupName() {
      const token = localStorage.getItem('access_token')?.trim();
      if (!token) {
        console.error('Token missing! User not authorized.');
        return;
      }

      if (!this.newGroupName) {
        alert('Group name cannot be empty!');
        return;
      }

      // Send request to update group name
      axios.put(`http://localhost:8080/api/group/${this.chatId}/updateName`, {
        name: this.newGroupName
      }, {
        headers: { Authorization: `Bearer ${token}` }
      })
      .then(() => {
        alert('Group name updated successfully!');
        this.chatName = this.newGroupName; // Update group name
        this.isEditingName = false; // Close edit form
      })
      .catch(error => {
        console.error('Error updating group name:', error);
      });
    },

    async getMessages() {
      const token = localStorage.getItem('access_token')?.trim();
      if (!token) {
        console.error('❌ Token missing! User not authorized.');
        return;
      }

      try {
        const response = await axios.get(`http://localhost:8080/api/chats/${this.chatId}`, {
          headers: { Authorization: `Bearer ${token}` }
        });

        // Update messages state and correctly set isMine flag
        this.messages = response.data.map(message => ({
          ...message,
          isMine: message.SenderID === JSON.parse(localStorage.getItem('user_id')), // Check if it's your message
          showForwardInput: false,  // Hide forward input by default
          forwardChatName: '',     // Clear field for new chat
          newComment: '',          // Clear field for new comment
        }));
      } catch (error) {
        console.error('❌ Error fetching messages:', error);
      }
    },

    enablePhotoEdit() {
      this.isEditingPhoto = true;
      this.newPhotoUrl = this.groupPhoto; // Pre-fill with current photo
    },

    saveGroupPhoto() {
      if (!this.newPhotoUrl) {
        alert('Photo URL cannot be empty');
        return;
      }

      const token = localStorage.getItem('access_token');
      axios.put(`http://localhost:8080/api/group/${this.chatId}/updatePhoto`, {
        photo_url: this.newPhotoUrl
      }, {
        headers: { Authorization: `Bearer ${token}` }
      })
      .then(() => {
        alert('Group photo updated successfully');
        this.groupPhoto = this.newPhotoUrl; // Update photo in UI
        this.isEditingPhoto = false; // Close edit form
      })
      .catch(error => {
        console.error('Error updating group photo:', error);
      });
    },

    cancelPhotoEdit() {
      this.isEditingPhoto = false;
      this.newPhotoUrl = ''; // Clear input field
    },

    async searchUsers() {
      if (this.searchQuery.length >= 2) {
        const token = localStorage.getItem('access_token')?.trim();
        if (!token) {
          console.error('Token missing! User not authorized.');
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
      // Add selected user to list
      if (!this.selectedUsers.some(u => u.ID === user.ID)) {
        this.selectedUsers.push(user);
      }

      // Clear search results after selection
      this.users = [];
      this.searchQuery = '';
      this.showSearch = false; // Hide search after adding user
    },

    removeUserFromSelection(user) {
      // Remove user from list
      this.selectedUsers = this.selectedUsers.filter(u => u.ID !== user.ID);
    },

    addUsersToGroup() {
      const token = localStorage.getItem('access_token')?.trim();
      if (!token) {
        console.error('Token missing! User not authorized.');
        return;
      }

      // Get selected user IDs
      const userIds = this.selectedUsers.map(user => user.ID);

      // Send request to add users to group
      axios.post(`http://localhost:8080/api/group/${this.chatId}/add`, {
        user_ids: userIds
      }, {
        headers: { Authorization: `Bearer ${token}` }
      })
      .then(() => {
        alert(`Users added to group!`);
        // Clear selected users list
        this.selectedUsers = [];
        this.$router.push({ path: '/chats' }).catch(() => {});
      })
      .catch(error => {
        console.error('Error adding users to group:', error);
        alert('User already added!')
      });
    },

    confirmExitGroup() {
      // Ask user for confirmation
      const confirmExit = window.confirm("Are you sure you want to exit the group?");
      if (confirmExit) {
        this.exitGroup(); // If confirmed, exit group
      }
    },

    exitGroup() {
      const token = localStorage.getItem('access_token')?.trim();
      if (!token) {
        console.error('Token missing! User not authorized.');
        return;
      }

      const chatId = this.$route.params.chatId;
      console.log(chatId)

      // Send request to exit group
      axios.post(`http://localhost:8080/api/group/${chatId}/leave`, {}, {
        headers: { Authorization: `Bearer ${token}` }
      })
      .then(() => {
        alert('You have exited the group');
        // Remove chat for current user and redirect to chat list
        this.$router.push('/chats');
      })
      .catch(error => {
        console.error('Error exiting group:', error);
      });
    },

    // Get group info (name and participant count)
    getGroupInfo() {
      const token = localStorage.getItem('access_token')?.trim();
      if (!token) {
        console.error('Token missing! User not authorized.');
        return;
      }

      axios.get(`http://localhost:8080/api/group/${this.chatId}/info`, {
        headers: { Authorization: `Bearer ${token}` }
      })
      .then(response => {
        this.groupName = response.data.name;
        this.groupPhoto = response.data.photo;
        this.usersCount = response.data.users_count;
      })
      .catch(error => {
        console.error('Error fetching group info:', error);
      });
    }
  },
  created() {
    this.debouncedSearch = _.debounce(this.searchUsers, 500); // Debounce searches
    this.getGroupInfo(); // Fetch group info on component creation
    this.getMessages();
  }
};
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #1e2d2b; 
  color: white;
}
.chat-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: #223432;
  padding: 15px;
  border-bottom: 2px solid #444;
}

.group-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
}

.group-photo {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  border: 2px solid #6a5acd;
}

h2 {
  font-size: 20px;
  margin: 5px 0;
  color: #ddd;
}

.header-actions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.header-actions button {
  background: #6a5acd;
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: 0.3s;
  font-size: 14px;
}

.header-actions button:hover {
  background: #5a4bb8;
}

.exit-btn {
  background: #ff5555 !important;
}

.exit-btn:hover {
  background: #cc4444 !important;
}

/* Общий стиль контейнера */
.form-group {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

/* Поля ввода */
input[type="text"] {
  width: 250px;
  padding: 8px;
  border: 2px solid #6a5acd;
  border-radius: 6px;
  background: #2a2a3a;
  color: #fff;
  font-size: 14px;
  outline: none;
  transition: 0.3s;
}

input[type="text"]:focus {
  border-color: #8c7ae6;
  box-shadow: 0 0 5px rgba(138, 43, 226, 0.8);
}

button {
  padding: 8px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: 0.3s;
}

.btn-success {
  background: #4caf50;
  color: white;
  box-shadow: 0px 4px 10px rgba(76, 175, 80, 0.3);
}

.btn-success:hover {
  background: #388e3c;
  box-shadow: 0px 4px 15px rgba(76, 175, 80, 0.5);
}

.btn-danger {
  background: #ff5555;
  color: white;
  box-shadow: 0px 4px 10px rgba(255, 85, 85, 0.3);
}

.btn-danger:hover {
  background: #d32f2f;
  box-shadow: 0px 4px 15px rgba(255, 85, 85, 0.5);
}

.add-user-btn {
  background: #6a5acd;
  color: white;
  padding: 10px 15px;
  font-weight: bold;
  display: inline-block;
}

.add-user-btn:hover {
  background: #5a4bb8;
}

.search-section input {
  width: 100%;
  max-width: 300px;
  padding: 10px;
  border: 2px solid #6a5acd;
  border-radius: 8px;
  background: #2a2a3a;
  color: white;
  font-size: 14px;
  transition: 0.3s;
}

.search-section input:focus {
  border-color: #8c7ae6;
  box-shadow: 0 0 10px rgba(138, 43, 226, 0.6);
}

.selected-users {
  background: #2a2a3a;
  padding: 15px;
  border-radius: 8px;
  margin-top: 10px;
}

.selected-users h3 {
  text-align: center;
  color: white;
}

.selected-users ul {
  padding: 0;
  list-style: none;
}

.selected-users li {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #3a3a4a;
  padding: 10px;
  border-radius: 6px;
  color: white;
  margin-bottom: 5px;
}

.selected-users button {
  background: transparent;
  border: none;
  color: #fafafa;
  font-size: 16px;
  cursor: pointer;
  transition: 0.3s;
}

.selected-users button:hover {
  color: #459435;
}
.user-list {
  list-style: none;
  padding: 10px;
  background: #2a2a3a;
  border-radius: 10px;
}

.user-list li {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-bottom: 1px solid #444;
  cursor: pointer;
  transition: 0.3s;
}

.user-list li:hover {
  background: #3a3a4a;
}

.user-photo {
  width: 30px;
  height: 30px;
  border-radius: 50%;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  display: flex;
  flex-direction: column;
}

.my-message, .other-message {
  padding: 8px 12px;
  border-radius: 12px;
  max-width: 60%;
  margin-bottom: 10px;
  word-wrap: break-word;
}

.my-message {
  background: #6a5acd; 
  align-self: flex-end;
  color: white;
}

.other-message {
  background: #444; 
  align-self: flex-start;
  color: white;
}

.comment {
  font-size: 14px;
  color: #bbb;
  margin-top: 5px;
}

.message-input {
  display: flex;
  padding: 10px;
  background: #223432;
  border-top: 1px solid #333;
}

.message-input input {
  flex: 1;
  padding: 10px;
  border: none;
  border-radius: 20px;
  background: #2e3f3c;
  color: white;
  outline: none;
}

.message-input button {
  padding: 10px 15px;
  margin-left: 10px;
  border: none;
  background: #6a5acd;
  color: white;
  border-radius: 10px;
  cursor: pointer;
}

.msg-action-btn {
  background: none;
  border: none;
  color: white;
  font-size: 12px;
  padding: 5px;
  cursor: pointer;
  transition: color 0.3s;
}

.forward-box {
  display: flex;
  gap: 5px;
  margin-top: 5px;
}

.forward-box input {
  flex: 1;
  padding: 6px;
  border-radius: 5px;
  background: #2e3f3c;
  color: white;
  border: none;
}

.forward-box button {
  background: #6a5acd;
  color: white;
  border: none;
  padding: 6px;
  border-radius: 5px;
  cursor: pointer;
}

</style>