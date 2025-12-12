<template>
  <div class="chat-container">
    <div class="chat-header">
      <div class="chat-photo" v-if="conversationPhoto">
        <img :src="'data:image/jpeg;base64,' + conversationPhoto" alt="Chat Thumbnail" />
      </div>
      <h3>{{ convName }}</h3>
    </div>
    <div class="chat-messages" ref="chatMessages">
      <p v-if="messages.length === 0">No messages yet...</p>
      <div
        v-for="message in messages"
        :key="message.id"
        class="message"
        :class="message.senderId === userToken ? 'self' : 'other'"
        :style="message.senderId !== userToken && conversationType === 'group' ? { paddingLeft: '45px' } : {}"
      >
        <div v-if="conversationType === 'group' && message.senderId !== userToken" class="sender-thumbnail">
          <img :src="'data:image/jpeg;base64,' + message.senderPhoto" alt="Sender Photo" />
        </div>
        <div class="message-content">
          <div v-if="message.replyTo" class="reply-preview">
            <small>Replying to {{ message.replySenderName || 'Unknown' }}: {{ message.replyContent }}</small>
            <img
              v-if="message.replyAttachment"
              :src="'data:image/jpeg;base64,' + message.replyAttachment"
              alt="Reply Attachment"
              class="reply-attachment"
            />
          </div>
          <p v-if="message.content.startsWith('<strong>Forwarded')" v-html="message.content"></p>
          <p v-else>
            <strong>
              {{ message.senderId === userToken ? 'You' : (message.senderName || 'Unknown Sender') }}:
            </strong>
            {{ message.content }}
          </p>
          <div v-if="message.attachment" class="attachment-container">
            <img :src="'data:image/jpeg;base64,' + message.attachment" alt="Attachment" class="attachment-image" />
          </div>
          <small>{{ formatTimestamp(message.timestamp) }}</small>
          <div v-if="message.reactionCount > 0" class="reaction-count">
            ❤️ × {{ message.reactionCount }}
            <div class="reactors-list">
              <ul>
                <li v-for="(reactor, idx) in message.reactingUserNames" :key="reactor">
                  {{ idx + 1 }}. {{ reactor }}
                </li>
              </ul>
            </div>
          </div>
          <div class="action-buttons">
            <button v-if="message.senderId !== userToken" class="action-button reply-button" @click.stop="setReply(message)">
              ↩
            </button>
            <button
              v-if="message.senderId !== userToken"
              class="action-button heart-button"
              :class="{ 'has-reacted': (message.reactingUserNames || []).includes(userName) }"
              :disabled="message.reactionLoading"
              @click.stop="toggleReaction(message)"
            >
              ❤️
            </button>
            <button class="action-button forward-button" @click.stop="showForwardOptions(message.id)">
              →
            </button>
            <button v-if="message.senderId === userToken" class="action-button delete-button" @click.stop="deleteMessage(message)">
              ✖
            </button>
          </div>
          <div v-if="messageOptions[message.id]?.showForwardMenu" class="forward-options" @click.stop>
            <select id="forward-select" class="forward-select" v-model="messageOptions[message.id].selectedConversationId">
              <option value="" disabled>Select conversation</option>
              <option v-for="conv in messageOptions[message.id].forwardConversations" :key="conv.id" :value="conv.id">
                {{ conv.name }}
              </option>
              <option value="new">New contact</option>
            </select>
            <div v-if="messageOptions[message.id].selectedConversationId === 'new'" class="contact-search">
              <input type="text" v-model="messageOptions[message.id].contactQuery" placeholder="Enter contact name" @input="searchContact(message.id)" />
              <ul v-if="messageOptions[message.id].contactResults.length > 0" class="contact-results">
                <li v-for="contact in messageOptions[message.id].contactResults" :key="contact.id" @click="selectContact(contact, message.id)" class="contact-result">
                  {{ contact.name }}
                </li>
              </ul>
            </div>
            <div class="forward-buttons-container">
              <button class="button-style" v-if="messageOptions[message.id].selectedConversationId !== 'new'" :disabled="!messageOptions[message.id].selectedConversationId" @click.stop="forwardMessage(messageOptions[message.id].selectedConversationId, message.id)">
                Send
              </button>
              <button class="button-style" v-if="messageOptions[message.id].selectedConversationId === 'new'" :disabled="!messageOptions[message.id].selectedContactId" @click.stop="forwardToContact(messageOptions[message.id].selectedContactId, message.id)">
                Send
              </button>
              <button class="button-style" @click.stop="closeForwardMenu(message.id)">Cancel</button>
            </div>
            <div v-if="messageOptions[message.id].forwardConversations.length === 0">
              No conversation found.
            </div>
          </div>
        </div>
        <div class="message-status" v-if="message.status && message.senderId !== userToken">
          {{ message.status }}
        </div>
      </div>
    </div>
    <div v-if="replyToMessage" class="reply-preview-box">
      <div class="reply-info">
        <strong>Replying to {{ replyToMessage.senderName || 'Unknown' }}:</strong>
        <span class="reply-text">{{ replyToMessage.content }}</span>
        <img v-if="replyToMessage.attachment" :src="'data:image/jpeg;base64,' + replyToMessage.attachment" alt="Reply Attachment" class="reply-attachment-preview" />
      </div>
      <button class="cancel-reply-button" @click="cancelReply">✖</button>
    </div>
    <div class="chat-input">
      <input type="file" ref="fileInput" style="display: none" accept="image/*, .gif" @change="handleFileSelect" />
      <button class="attach-button" @click="triggerFileInput">
        Attach Image or GIF
        <span v-if="selectedFile" class="file-icon">🖼️</span>
      </button>
      <input v-model="message" class="message-input" type="text" placeholder="Type a message..." @input="toggleSendButton" />
      <button v-if="message.trim() || selectedFile" class="send-button" @click="sendMessage">
        Send
      </button>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
export default {
  name: "ChatView",
  data() {
    return {
      message: "",
      messages: [],
      conversations: [],
      userToken: localStorage.getItem("token"),
      convName: localStorage.getItem("conversationName") || "Unknown User",
      conversationPhoto: null,
      conversationType: null,
      conversationId: this.$route.params.uuid,
      messageOptions: {},
      selectedFile: null,
      pollIntervalId: null,
      firstLoad: true,
      replyToMessage: null
    };
  },
  computed: {
    userName() {
      return localStorage.getItem("name");
    }
  },
  methods: {

    triggerFileInput() {
      this.$refs.fileInput.click();
    },
    handleFileSelect(event) {
      this.selectedFile = event.target.files[0];
    },
    async sendMessage() {
      const token = localStorage.getItem("token");
      if (!token) {
        this.$router.push({ path: "/" });
        return;
      }
      const formData = new FormData();
      formData.append("content", this.message);
      if (this.replyToMessage) {
        formData.append("replyTo", this.replyToMessage.id);
      }
      if (this.selectedFile) {
        formData.append("attachment", this.selectedFile);
      }
      await axios.post(`/conversations/${this.conversationId}/message`, formData, {
        headers: { Authorization: `Bearer ${token}` }
      });
      this.message = "";
      this.selectedFile = null;
      this.$refs.fileInput.value = "";
      this.replyToMessage = null;
      await this.fetchMessages();
      this.$nextTick(() => {
        this.forceScrollToBottom();
      });
    },
    async fetchMessages() {
      const token = localStorage.getItem("token");
      if (!token) {
        this.$router.push({ path: "/" });
        return;
      }
      const response = await axios.get(`/conversations/${this.conversationId}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      this.messages = (response.data.messages || []).map(msg => ({
        ...msg,
        reactingUserNames: msg.reactingUserNames || [],
        showReactedList: false
      }));
      if (response.data.name) {
        this.convName = response.data.name;
      }
      if (response.data.conversationPhoto && response.data.conversationPhoto.String) {
        this.conversationPhoto = response.data.conversationPhoto.String;
      } else {
        this.conversationPhoto = null;
      }
      this.conversationType = response.data.type || "direct";
      this.$nextTick(() => {
        if (this.firstLoad) {
          this.forceScrollToBottom();
          this.firstLoad = false;
        }
      });
    },
    forceScrollToBottom() {
      const chat = this.$refs.chatMessages;
      if (chat) {
        chat.scrollTop = chat.scrollHeight;
      }
    },
    async toggleReaction(message) {
      const token = localStorage.getItem("token");
      if (!token || message.senderId === this.userToken) return;
      const hasReacted = (message.reactingUserNames || []).includes(this.userName);
      try {
        if (hasReacted) {
          await axios.delete(`/conversations/${this.conversationId}/message/${message.id}/comment`, {
            headers: { Authorization: `Bearer ${token}` }
          });
        } else {
          await axios.post(`/conversations/${this.conversationId}/message/${message.id}/comment`, {},
            { headers: { Authorization: `Bearer ${token}` } }
          );
        }
      } catch (err) {
        console.error("Error toggling reaction", err);
      } finally {
        await this.fetchMessages();
      }
    },
    async deleteMessage(message) {
      const token = localStorage.getItem("token");
      if (!token) {
        this.$router.push({ path: "/" });
        return;
      }
      await axios.delete(`/conversations/${this.conversationId}/message/${message.id}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      this.messages = this.messages.filter(m => m.id !== message.id);
    },
    formatTimestamp(timestamp) {
      const date = new Date(timestamp);
      return date.toLocaleString();
    },
    showForwardOptions(messageId) {
      this.closeAllMenus();
      if (!this.messageOptions[messageId]) {
        this.messageOptions[messageId] = {
          showForwardMenu: true,
          forwardConversations: [],
          selectedConversationId: "",
          contactQuery: "",
          contactResults: [],
          selectedContactId: ""
        };
        this.fetchForwardConversations(messageId);
      } else {
        this.messageOptions[messageId].showForwardMenu = !this.messageOptions[messageId].showForwardMenu;
      }
    },
    closeForwardMenu(messageId) {
      if (this.messageOptions[messageId]) {
        this.messageOptions[messageId].showForwardMenu = false;
      }
    },
    closeAllMenus() {
      for (const id in this.messageOptions) {
        this.messageOptions[id].showForwardMenu = false;
      }
    },
    handleOutsideClick(event) {
      const messageContent = this.$el.querySelector('.message-content');
      if (messageContent && !messageContent.contains(event.target)) {
        this.closeAllMenus();
      }
    },
    async fetchForwardConversations(messageId) {
      const token = localStorage.getItem("token");
      const response = await axios.get('/conversations', {
        headers: { Authorization: `Bearer ${token}` }
      });
      const conversations = response.data.filter(conv => conv.id !== this.conversationId);
      this.messageOptions[messageId].forwardConversations = conversations;
    },
    
    async searchContact(messageId) {
      const query = this.messageOptions[messageId].contactQuery;
      if (!query.trim()) {
        this.messageOptions[messageId].contactResults = [];
        return;
      }
      const token = localStorage.getItem("token");
      const response = await axios.get('/search', {
        params: { username: query },
        headers: { Authorization: `Bearer ${token}` }
      });
      this.messageOptions[messageId].contactResults = response.data;
    },
    selectContact(contact, messageId) {
      this.messageOptions[messageId].selectedContactId = contact.id;
      this.messageOptions[messageId].contactQuery = contact.name;
      this.messageOptions[messageId].contactResults = [];
    },
    async forwardToContact(selectedContactId, messageId) {
      const token = localStorage.getItem("token");
      if (!token) {
        this.$router.push({ path: "/" });
        return;
      }
      const conversationResponse = await axios.post(
        `/conversations`,
        { senderId: token, recipientId: selectedContactId },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      const targetConversationId = conversationResponse.data.conversationId;
      const forwarderName = localStorage.getItem("name") || "Unknown";
      await axios.post(
        `/conversations/${this.conversationId}/message/${messageId}/forward`,
        { sourceMessageId: messageId, targetConversationId: targetConversationId, forwarderName: forwarderName },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      alert("Message forwarded successfully!");
      this.closeForwardMenu(messageId);
    },
    async forwardMessage(targetConversationId, messageId) {
      const message = this.messages.find(m => m.id === messageId);
      if (!message) return;
      const token = localStorage.getItem("token");
      const forwarderName = localStorage.getItem("name") || "Unknown";
      await axios.post(
        `/conversations/${this.conversationId}/message/${messageId}/forward`,
        { targetConversationId: targetConversationId, forwarderName: forwarderName },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      alert("Message forwarded successfully!");
      this.closeForwardMenu(messageId);
    },
    setReply(message) {
      this.replyToMessage = message;
    },
    cancelReply() {
      this.replyToMessage = null;
    }
  },
  mounted() {
    this.fetchMessages();
    this.pollIntervalId = setInterval(() => {
      this.fetchMessages();
    }, 5000);
    document.addEventListener("click", this.handleOutsideClick);
  },
  beforeUnmount() {
    document.removeEventListener("click", this.handleOutsideClick);
    clearInterval(this.pollIntervalId);
  }
};
</script>

<style scoped>
/* ------------------------------------
   COLOR SYSTEM (EASY TO TUNE LATER)
------------------------------------- */
:root {
  --brand-color: #4a90e2;
  --brand-light: #e6f1fc;
  --self-bubble: #dff5e1;
  --other-bubble: #f2f4f7;
  --border-light: #d9d9d9;
  --text-dark: #222;
  --text-light: #555;
}

/* ------------------------------------
   CHAT LAYOUT
------------------------------------- */
.chat-container {
  display: flex;
  flex-direction: column;
  height: 92vh;
  background: #ffffff;
  overflow: hidden;
  border-radius: 12px;
  border: 1px solid var(--border-light);
}

/* ------------------------------------
   HEADER
------------------------------------- */
/* CHAT HEADER — modern and readable */
.chat-header {
  display: flex;
  align-items: center;
  padding: 15px;
  background-color: white; 
  border-bottom: 1px solid #3a7bc1;
  color:  #232323ff; 
}

/* conversation title */
.chat-header h3 {
  font-size: 30px;
  margin: 0;
  color: #232323ff;
  font-weight: 600;
  letter-spacing: 0.5px;
}


.chat-photo {
  width: 42px;
  height: 42px;
  margin-right: 12px;
}

.chat-photo img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

/* ------------------------------------
   MESSAGES AREA
------------------------------------- */
.chat-messages {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow-y: auto;
  padding: 22px;
  background: #fafafa;
  border-top: 1px solid #ccc;
  border-bottom: 1px solid #ccc;
}

/* ------------------------------------
   MESSAGE BUBBLES
------------------------------------- */
.message {
  position: relative;
  max-width: 68%;
  padding: 12px 16px;
  margin-bottom: 14px;
  border-radius: 12px;
  background-color: #f3fffeff;
  font-size: 0.95rem;
  box-shadow: 0 1px 3px rgba(0,0,0,0.07);
}

/* User message */
.message.self {
  margin-left: auto;
  background-color: #bcf9deff;
}

/* Group thumbnails */
.sender-thumbnail {
  position: absolute;
  left: 10px;
  top: 10px;
  width: 32px;
  height: 32px;
}

.sender-thumbnail img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

/* ------------------------------------
   TEXT & TIMESTAMP
------------------------------------- */
.message p {
  margin: 0;
  color: #2f2f2fff;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.message small {
  margin-top: 6px;
  font-size: 0.75rem;
  color: #777;
}

.message-content {
  position: relative;
  min-height: 40px;
}
/* ------------------------------------
   ATTACHMENTS
------------------------------------- */
.attachment-container {
  margin-top: 8px;
  width: 260px;
  max-height: 260px;
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid #d9dde0ff;
}

.attachment-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* ------------------------------------
   ACTION BUTTONS (Reply / Heart / Delete)
------------------------------------- */
.action-buttons {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  right: -40px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  opacity: 0;
  transition: opacity 0.25s;
}

.message.self .action-buttons {
  right: auto;
  left: -40px;
}

.message:hover .action-buttons {
  opacity: 1;
}

/* ACTION BUTTONS - Always visible, modern style */
.action-button {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: #f1f3f5; /* light gray */
  border: 1px solid #c5cbd3;
  font-size: 12px;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #6c7a89; /* icon color */
  transition: all 0.2s ease;
}

.action-button:hover {
  background-color: var(--brand-light);
  border-color: var(--brand-color);
  color: var(--brand-color);
}
.reply-button {
  font-size: 10px;
  margin-right: 5px;
}

/* If it’s the liked heart */
.has-reacted {
  background-color: #ffe3e3 !important;
  border-color: #ff8b8b !important;
  color: #ff4d4d !important;
}


/* ------------------------------------
   FORWARD MENU
------------------------------------- */
.forward-options {
  position: absolute;
  right: 0;
  top: 32px;
  width: 260px;
  background: #fff;
  border-radius: 8px;
  border: 1px solid #ccc;
  padding: 12px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.09);
  z-index: 100;
}

.forward-select,
.contact-search input {
  width: 100%;
  padding: 8px;
  border-radius: 6px;
  border: 1px solid #ccc;
  font-size: 14px;
  margin-bottom: 8px;
}
.forward-buttons-container {
  display: flex;
  gap: 10px;
  justify-content: center;
  margin-top: 10px;
}

/* Results */
.contact-results {
  max-height: 120px;
  overflow-y: auto;
  border: 1px solid #ccc;
  border-radius: 6px;
}

.contact-result {
  padding: 8px;
  cursor: pointer;
  transition: background 0.2s;
  border-bottom: 1px solid #eee;
}

.contact-result:hover {
  background-color: #f0f0f0;
}

.button-style {
  flex: 1;
  padding: 10px;
  background: #4a95ebff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.button-style:hover {
  background: #3a74b6;
}
.file-icon {
  font-size: 18px;
  margin-left: 5px;
}
.message-status {
  position: absolute;
  bottom: 5px;
  right: 10px;
  font-size: 12px;
  color: #646464ff;
}
.reactors-list ul {
  margin: 0;
  padding: 0;
  list-style: none;
  font-size: 0.8em;
  color: #444;
}
.reactors-list li {
  margin: 2px 0;
}

/* ------------------------------------
   REPLY PREVIEW
------------------------------------- */
.reply-preview-box {
  background-color: #d9dde0ff;
  border-left: 4px solid #4a95ebff;
  padding: 10px 14px;
  margin: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.reply-info {
  font-size: 0.9em;
  color: #444;
}
.reply-attachment,
.reply-attachment-preview {
  width: 21px;
  height: 21px;
  object-fit: cover;
  margin-left: 10px;
  border-radius: 4px;
}
.cancel-reply-button {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: #2f2f2fff;
}

/* ------------------------------------
   INPUT AREA
------------------------------------- */
.chat-input {
  display: flex;
  align-items: center;
  padding: 14px;
  border-top: 1px solid #d9dde0ff;
  gap: 10px;
}

.attach-button {
  background: #4a95ebff;
  color: white;
  border: none;
  padding: 10px 14px;
  border-radius: 20px;
  cursor: pointer;
  font-size: 13px;
}

.attach-button:hover {
  background: #3a74b6;
}

.message-input {
  flex: 1;
  padding: 11px 14px;
  border: 1px solid #d9dde0ff;
  border-radius: 20px;
  font-size: 14px;
  outline: none;
}

.send-button {
  background: #4a95ebff;
  color: white;
  border: none;
  padding: 11px 24px;
  border-radius: 20px;
  cursor: pointer;
  transition: background 0.2s;
}

.send-button:hover {
  background: #3a74b6;
}
</style>


