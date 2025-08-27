<template>
  <div class="home-wrapper">
    <header class="header">
      <h2>{{ username }}, your chats</h2>
      <div class="actions">
        <button @click="refreshList" class="btn small secondary">Refresh</button>
        <button @click="logoutUser" class="btn small secondary">Log Out</button>
        <button @click="navigateToGroupCreation" class="btn small primary">New Group</button>
      </div>
    </header>

    <ErrorMsg v-if="errorMessage" :msg="errorMessage" />

    <section class="chat-section">
      <p v-if="chatThreads.length === 0" class="empty-text">No chat threads found.</p>
      <div v-else class="chat-list">
        <div
          v-for="chat in chatThreads"
          :key="chat.id"
          class="chat-card"
          @click="goToChat(chat.id, chat.name)"
        >
          <div class="chat-avatar">
            <img
              v-if="chat.conversationPhoto.String"
              :src="'data:image/png;base64,' + chat.conversationPhoto.String"
              alt="User"
              class="avatar"
            />
          </div>
          <div class="chat-meta">
            <h4 class="chat-title">{{ chat.name }}</h4>
            <p v-if="chat.lastMessage" class="chat-snippet">
              <strong>{{ chat.lastMessage.senderName }}:</strong>
              <img
                v-if="chat.lastMessage.attachment"
                :src="'data:image/*;base64,' + chat.lastMessage.attachment"
                class="chat-thumb"
                alt="Attachment"
              />
              <span v-if="isForwarded(chat.lastMessage)" v-html="formatMessage(chat.lastMessage)"></span>
              <span v-else>{{ formatMessage(chat.lastMessage) }}</span>
              <span class="timestamp">{{ formatTime(chat.lastMessage.timestamp) }}</span>
            </p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "HomeView",
  components: { ErrorMsg },
  data() {
    localStorage.removeItem("recipientId");
    return {
      username: "",
      chatThreads: [],
      errorMessage: null,
      intervalId: null,
    };
  },
  methods: {
    async fetchChats() {
      this.errorMessage = null;
      try {
        const token = localStorage.getItem("token");
        if (!token) return this.$router.push("/");

        const { data } = await this.$axios.get("/chats", {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.chatThreads = data || [];
      } catch (err) {
        this.errorMessage = "Unable to load chats.";
      }
    },
    goToChat(id, name) {
      localStorage.setItem("conversationName", name);
      this.$router.push(`/conversations/${id}`);
    },
    formatMessage(message) {
      return this.isForwarded(message) ? message.content : this.shorten(message.content);
    },
    shorten(text, limit = 60, suffix = "...") {
      if (!text || text.length <= limit) return text;
      const cut = text.substring(0, limit).lastIndexOf(" ");
      return text.substring(0, cut > 0 ? cut : limit) + suffix;
    },
    isForwarded(msg) {
      return msg.content.includes("Forwarded from");
    },
    formatTime(ts) {
      return new Date(ts).toLocaleString();
    },
    refreshList() {
      this.fetchChats();
    },
    logoutUser() {
      this.$router.push("/");
    },
    navigateToGroupCreation() {
      this.$router.push("/new-group");
    },
  },
  mounted() {
    this.username = localStorage.getItem("name") || "Guest";
    this.fetchChats();
    this.intervalId = setInterval(this.fetchChats, 1000);
  },
  unmounted() {
    clearInterval(this.intervalId);
  },
};
</script>

<style scoped>
.home-wrapper {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.btn.small {
  font-size: 14px;
}
.btn.secondary {
  background-color: #e0e0e0;
}
.btn.primary {
  background-color: #007bff;
  color: white;
}

.chat-section {
  margin-top: 20px;
}

.empty-text {
  color: #888;
  font-style: italic;
}

.chat-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chat-card {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 12px;
  display: flex;
  align-items: flex-start;
  gap: 16px;
  cursor: pointer;
  transition: background 0.2s;
}
.chat-card:hover {
  background-color: #efefef;
}

.chat-avatar {
  width: 60px;
  height: 60px;
}
.avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  object-fit: cover;
}

.chat-meta {
  flex: 1;
}

.chat-title {
  margin: 0;
  font-size: 16px;
  font-weight: bold;
}

.chat-snippet {
  margin-top: 4px;
  font-size: 14px;
  color: #555;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
}

.chat-thumb {
  width: 18px;
  height: 18px;
  object-fit: cover;
  border-radius: 3px;
}

.timestamp {
  font-size: 12px;
  color: #999;
}
</style>

