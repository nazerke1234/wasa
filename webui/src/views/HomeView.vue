<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">All conversations</h1>
      <div class="btn-toolbar mb-2 mb-md-0 button-group-container">
        
        <button type="button" class="btn-chat-primary" @click="startChat">
          Start Chat
        </button>

        <button type="button" class="btn-chat-secondary" @click="newGroup">
          New group
        </button>

        <button type="button" class="btn-chat-tertiary" @click="logOut">
          Log Out
        </button>
      </div>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg" />
    <div>
      <div v-if="conversations.length === 0">
        <p>You have no conversations yet. Click "Start Chat" to begin.</p>
      </div>
      <div v-else class="conversations-container">
        <div
          v-for="conv in conversations"
          :key="conv.id"
          class="conversation-block"
          @click="viewConversation(conv.id, conv.name)"
        >
          <div class="conversation-photo">
            <img
              v-if="conv.conversationPhoto.String"
              :src="'data:image/png;base64,' + conv.conversationPhoto.String"
              alt="Profile Picture"
              class="profile-picture"
            />
          </div>
          <div class="conversation-details">
            <h4>{{ conv.name }}</h4>
            <p v-if="conv.lastMessage" class="last-message">
              Last message by {{ conv.lastMessage.senderName }}:
              <img v-if="conv.lastMessage.attachment"
                   :src="'data:image/*;base64,' + conv.lastMessage.attachment"
                   class="attachment-thumbnail"
                   alt="Attachment">
              <span v-if="isForwarded(conv.lastMessage)" v-html="getFormattedMessage(conv.lastMessage)"></span>
              <span v-else>{{ getFormattedMessage(conv.lastMessage) }}</span>
              at {{ new Date(conv.lastMessage.timestamp).toLocaleString() }}
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
      username: "",
      errormsg: null,
      loading: false,
      conversations: [],
      pollIntervalId: null,
    };
  },
  methods: {
    async loadConversations() {
      this.errormsg = null;
      this.loading = true;
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          this.$router.push({ path: "/" });           //go to login
          return;
        }
        const response = await this.$axios.get("/conversations", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        this.conversations = response.data || [];
      } catch (error) {
        console.error("Error loading conversations:", error);
        this.errormsg = "Failed to load conversations. Please try again.";
      } finally {
        this.loading = false;
      }
    },
    viewConversation(conversationId, conversationName) {
      localStorage.setItem("conversationName", conversationName);
      this.$router.push({
        path: `/conversations/${conversationId}`,         //go to chatview
      });
    },
    truncateText(text, length = 50, clamp = '...') {
      if (!text || text.length <= length) {
        return text;
      }
      const lastSpaceIndex = text.substring(0, length).lastIndexOf(' '); // cut when the word ends
      if (lastSpaceIndex === -1) { //no index found
        return text.substring(0, length) + clamp;
      }
      return text.substring(0, lastSpaceIndex) + clamp;
    },
    isForwarded(message) {
      return message.content.includes("<strong>Forwarded from");
    },
    getFormattedMessage(message) {
      if (this.isForwarded(message)) {
        return message.content;    
      }
      return this.truncateText(message.content);
    },
    startChat() {
      this.$router.push({ path: "/search" });
    },
    logOut() {
      this.$router.push({ path: "/" });
    },
    newGroup() {
      this.$router.push({ path: "/new-group" });
    },
  },
  mounted() {                                      //implemented at the time the page loaded
    this.username = localStorage.getItem("name") || "Guest";
    this.loadConversations();
    this.pollIntervalId = setInterval(() => {
      this.loadConversations();
    }, 5000);                                      //refresh the page every 5 seconds
  },
  unmounted() {                                    //when i close the page(component)
    clearInterval(this.pollIntervalId);            //stop timer
  },
};
</script>. 

<style>
.pt-3 {
    padding-top: 1.5rem !important;
}
.pb-2 {
    padding-bottom: 0.5rem !important;
}
.mb-4 {
    margin-bottom: 1.5rem !important; 
}

.conversation-header {
    font-size: 2rem;
    font-weight: 600;
    color: #2c3e50;
}

.button-group-container {
    display: flex;
    gap: 10px; 
}


/* Start Chat */
.btn-chat-primary {
    background-color: #4a90e2;
    color: white;
    border: 1px solid #4a90e2;
    padding: 10px 18px;
    font-weight: 600;
    border-radius: 8px;
    transition: background-color 0.2s, box-shadow 0.2s;
}

.btn-chat-primary:hover {
    background-color: #3a74b6;
    box-shadow: 0 4px 10px rgba(74, 144, 226, 0.3);
}

/* New group */
.btn-chat-secondary {
    background-color: white;
    color: #4a90e2;
    border: 1px solid #4a90e2;
    padding: 10px 18px;
    font-weight: 600;
    border-radius: 8px;
    transition: background-color 0.2s, color 0.2s;
}

.btn-chat-secondary:hover {
    background-color: #e6f1fc;
    color: #3a74b6;
}

/* Log Out */
.btn-chat-tertiary {
    background-color: transparent;
    color: #7f8c8d;
    border: 1px solid #bdc3c7;
    padding: 10px 18px;
    font-weight: 500;
    border-radius: 8px;
    transition: background-color 0.2s, color 0.2s;
}

.btn-chat-tertiary:hover {
    background-color: #ecf0f1;
    color: #34495e;
    border-color: #34495e;
}

.conversations-container {
  display: flex;
  flex-direction: column;
}

.conversation-block {
  background-color: white; 
  padding: 15px;
  margin-bottom: 12px;
  cursor: pointer;
  border-radius: 10px;
  display: flex;
  align-items: center; 
  gap: 15px;
  border: 1px solid #ecf0f1; 
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05); 
  transition: all 0.2s ease; 
}

.conversation-photo {
  flex-shrink: 0; 
  width: 50px; 
  height: 50px; 
}

.profile-picture {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 50%;
  border: 2px solid #4a90e2;
}

.conversation-details h4 {
  margin-top: 0;
  margin-bottom: 2px;
  font-size: 1.1rem;
  font-weight: 600;
}

.last-message {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 0.9rem;
  color: #7f8c8d;
}

.attachment-thumbnail {
  width: 18px;
  height: 18px;
  object-fit: cover;
  border-radius: 2px;
  flex-shrink: 0;
}

.empty-state {
    padding: 40px;
    text-align: center;
    color: #95a5a6;
    font-style: italic;
}

@media (max-width: 600px) {
  .conversation-block p {
    -webkit-line-clamp: 3;
    line-clamp: 3;
  }
}
</style>



