<template>
  <div>
    <div class="page-header">
      <h1 class="main-title">{{ userDisplayName }}, your conversations</h1>
      <div class="controls-container">
        <div class="control-group">
          <button type="button" class="control-btn secondary" @click="reloadData">Reload</button>
          <button type="button" class="control-btn secondary" @click="signOut">Sign Out</button>
        </div>
        <div class="control-group">
          <button type="button" class="control-btn primary" @click="initiateGroup">Create Group</button>
        </div>
      </div>
    </div>
    <ErrorMsg v-if="errorMessage" :msg="errorMessage" />
    <div>
      <div v-if="chatList.length === 0">
        <p>No conversations available.</p>
      </div>
      <div v-else class="chat-list">
        <div
          v-for="chat in chatList"
          :key="chat.id"
          class="chat-item"
          @click="openChat(chat.id, chat.name)"
        >
          <div class="chat-avatar">
            <img
              v-if="chat.conversationPhoto.String"
              :src="'data:image/png;base64,' + chat.conversationPhoto.String"
              alt="Chat Avatar"
              class="avatar-image"
            />
          </div>
          <div class="chat-info">
            <h4>{{ chat.name }}</h4>
            <p v-if="chat.lastMessage" class="recent-message">
              Last message from {{ chat.lastMessage.senderName }}:
              <img v-if="chat.lastMessage.attachment"
                   :src="'data:image/*;base64,' + chat.lastMessage.attachment"
                   class="message-thumbnail"
                   alt="Media attachment">
              <span v-if="checkIfForwarded(chat.lastMessage)" v-html="formatMessageContent(chat.lastMessage)"></span>
              <span v-else>{{ formatMessageContent(chat.lastMessage) }}</span>
              at {{ new Date(chat.lastMessage.timestamp).toLocaleString() }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "HomeView",
  components: {
    ErrorMsg,
  },
  data() {
    localStorage.removeItem("recipientId");
    return {
      userDisplayName: "",
      errorMessage: null,
      isLoading: false,
      chatList: [],
      autoRefreshTimer: null,
    };
  },
  methods: {
    async fetchChats() {
      this.errorMessage = null;
      this.isLoading = true;
      try {
        const authToken = localStorage.getItem("token");
        if (!authToken) {
          this.$router.push({ path: "/" });
          return;
        }
        const apiResponse = await this.$axios.get("/conversations", {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        });
        this.chatList = apiResponse.data || [];
      } catch (apiError) {
        console.error("Error fetching conversations:", apiError);
        this.errorMessage = "Unable to load conversations. Please try again.";
      } finally {
        this.isLoading = false;
      }
    },
    openChat(chatId, chatTitle) {
      localStorage.setItem("conversationName", chatTitle);
      this.$router.push({
        path: `/conversations/${chatId}`,
      });
    },
    shortenText(content, maxLength = 50, suffix = '...') {
      if (!content || content.length <= maxLength) {
        return content;
      }
      const lastSpacePosition = content.substring(0, maxLength).lastIndexOf(' ');
      if (lastSpacePosition === -1) {
        return content.substring(0, maxLength) + suffix;
      }
      return content.substring(0, lastSpacePosition) + suffix;
    },
    checkIfForwarded(message) {
      return message.content.includes("<strong>Forwarded from");
    },
    formatMessageContent(message) {
      if (this.checkIfForwarded(message)) {
        return message.content;
      }
      return this.shortenText(message.content);
    },
    reloadData() {
      this.fetchChats();
    },
    signOut() {
      this.$router.push({ path: "/" });
    },
    initiateGroup() {
      this.$router.push({ path: "/new-group" });
    },
  },
  mounted() {
    this.userDisplayName = localStorage.getItem("name") || "User";
    this.fetchChats();
    this.autoRefreshTimer = setInterval(() => {
      this.fetchChats();
    }, 1000);
  },
  unmounted() {
    clearInterval(this.autoRefreshTimer);
  },
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

.main-title {
  font-size: 1.75rem;
  font-weight: 600;
  color: #333;
}

.controls-container {
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

.chat-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.chat-item {
  background-color: #f8f9fa;
  padding: 1rem;
  cursor: pointer;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: background-color 0.2s ease;
  border: 1px solid #e9ecef;
}

.chat-item:hover {
  background-color: #e9ecef;
}

.chat-avatar {
  flex-shrink: 0;
  width: 75px;
  height: 75px;
}

.avatar-image {
  width: 75px;
  height: 75px;
  object-fit: cover;
  border-radius: 50%;
  border: 2px solid #dee2e6;
}

.chat-info {
  flex: 1;
  min-width: 0;
}

.chat-info h4 {
  margin: 0 0 0.5rem 0;
  color: #2d3748;
  font-size: 1.125rem;
}

.recent-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0;
  color: #6b7280;
  font-size: 0.875rem;
  line-height: 1.4;
  flex-wrap: wrap;
}

.message-thumbnail {
  width: 20px;
  height: 20px;
  object-fit: cover;
  border-radius: 0.25rem;
  flex-shrink: 0;
  border: 1px solid #d1d5db;
}

@media (max-width: 600px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .controls-container {
    width: 100%;
    justify-content: space-between;
  }
  
  .chat-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }
  
  .chat-avatar {
    align-self: center;
  }
  
  .recent-message {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>

